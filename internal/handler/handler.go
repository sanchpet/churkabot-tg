// Package handler turns incoming messages into the bot's replies.
package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/sanchpet/churkabot-tg/internal/rules"
)

const (
	// Telegram keeps unconfirmed updates for a day; after downtime the bot should not answer
	// yesterday's messages.
	staleAfter = 10 * time.Minute
	// Telegram allows about 20 messages a minute in a group; one sticker per chat per cooldown
	// keeps a run of trigger words from turning into a flood.
	cooldown = 15 * time.Second
	// Handlers run on a context that outlives shutdown, so a reply in flight is finished, not cut.
	sendTimeout = 15 * time.Second
)

// Sender is the part of the Bot API the handler uses.
type Sender interface {
	SendSticker(ctx context.Context, params *bot.SendStickerParams) (*models.Message, error)
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
	SendAnimation(ctx context.Context, params *bot.SendAnimationParams) (*models.Message, error)
}

// Tease answers a share of one user's group messages, whatever they say, with the Uimis animation.
// A zero UserID turns it off.
type Tease struct {
	UserID int64
	Chance float64
}

// Handler answers group messages with stickers and admins' private messages with diagnostics.
type Handler struct {
	api    Sender
	admins map[int64]bool
	tease  Tease
	log    *slog.Logger
	now    func() time.Time
	roll   func() float64

	mu        sync.Mutex
	lastReply map[int64]time.Time
}

// New returns a Handler that sends through api and treats the given user IDs as admins.
func New(api Sender, admins map[int64]bool, tease Tease, log *slog.Logger) *Handler {
	return &Handler{api: api, admins: admins, tease: tease, log: log, now: time.Now, roll: rand.Float64, lastReply: map[int64]time.Time{}}
}

// Handle processes one update.
func (h *Handler) Handle(ctx context.Context, u *models.Update) {
	m := u.Message
	// Bots may now see each other's messages in groups; answering them risks a loop. Automatic
	// forwards are channel posts copied into a linked group, not something a member said.
	if m == nil || m.From == nil || m.IsAutomaticForward || (m.From.IsBot && !anonymousAdmin(m)) {
		return
	}
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), sendTimeout)
	defer cancel()

	switch m.Chat.Type {
	case models.ChatTypePrivate:
		h.private(ctx, m)
	case models.ChatTypeGroup, models.ChatTypeSupergroup:
		h.group(ctx, m)
	}
}

func (h *Handler) group(ctx context.Context, m *models.Message) {
	if h.now().Sub(time.Unix(int64(m.Date), 0)) > staleAfter {
		return
	}
	if h.tease.UserID != 0 && m.From.ID == h.tease.UserID && h.roll() < h.tease.Chance {
		h.teaseReply(ctx, m)
		return
	}
	text, entities := content(m)
	match, ok := rules.Find(text, entities)
	if !ok {
		return
	}
	log := h.log.With("chat_id", m.Chat.ID, "user_id", m.From.ID, "rule", match.Rule, "word", match.Word)
	if !h.allow(m.Chat.ID) {
		log.Debug("reaction skipped: cooldown")
		return
	}
	_, err := h.api.SendSticker(ctx, &bot.SendStickerParams{
		ChatID:          m.Chat.ID,
		Sticker:         &models.InputFileString{Data: match.Sticker.FileID},
		ReplyParameters: &models.ReplyParameters{MessageID: m.ID},
	})
	if err != nil {
		log.Warn("reaction failed", "error", err)
		return
	}
	log.Info("reaction sent", "sticker", match.Sticker.Name)
}

func (h *Handler) teaseReply(ctx context.Context, m *models.Message) {
	log := h.log.With("chat_id", m.Chat.ID, "user_id", m.From.ID)
	if !h.allow(m.Chat.ID) {
		log.Debug("tease skipped: cooldown")
		return
	}
	_, err := h.api.SendAnimation(ctx, &bot.SendAnimationParams{
		ChatID:          m.Chat.ID,
		Animation:       &models.InputFileString{Data: rules.Uimis.FileID},
		ReplyParameters: &models.ReplyParameters{MessageID: m.ID},
	})
	if err != nil {
		log.Warn("tease failed", "error", err)
		return
	}
	log.Info("tease sent")
}

// allow reports whether the chat is out of its cooldown, and starts a new one if it is.
func (h *Handler) allow(chatID int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := h.now()
	if now.Sub(h.lastReply[chatID]) < cooldown {
		return false
	}
	h.lastReply[chatID] = now
	return true
}

func (h *Handler) private(ctx context.Context, m *models.Message) {
	if !h.admins[m.From.ID] {
		if command(m.Text) == "/start" {
			h.reply(ctx, m, "Я работаю в группах: добавь меня в чат.")
		}
		h.log.Info("private message from a non-admin ignored", "user_id", m.From.ID)
		return
	}
	switch {
	case m.Sticker != nil:
		s := m.Sticker
		h.reply(ctx, m, fmt.Sprintf("file_id: %s\nfile_unique_id: %s\nset: %s\nemoji: %s", s.FileID, s.FileUniqueID, s.SetName, s.Emoji))
	case m.Animation != nil:
		a := m.Animation
		h.reply(ctx, m, fmt.Sprintf("file_id: %s\nfile_unique_id: %s\nsize: %dx%d, %ds", a.FileID, a.FileUniqueID, a.Width, a.Height, a.Duration))
	case len(m.Photo) > 0:
		p := m.Photo[0]
		for _, size := range m.Photo[1:] {
			if size.Width*size.Height > p.Width*p.Height {
				p = size
			}
		}
		h.reply(ctx, m, fmt.Sprintf("file_id: %s\nsize: %dx%d", p.FileID, p.Width, p.Height))
	case command(m.Text) == "/start":
		h.reply(ctx, m, fmt.Sprintf("Привет, %s! Я чуркабот.", m.From.FirstName))
	case command(m.Text) == "/help":
		h.reply(ctx, m, "Пришли стикер, гифку или фото — отвечу их file_id.\nПришли текст — скажу, какое правило на него сработает в группе.")
	default:
		text, entities := content(m)
		if text == "" {
			return
		}
		match, ok := rules.Find(text, entities)
		switch {
		case !ok:
			h.reply(ctx, m, "no match")
		case match.Word == "":
			h.reply(ctx, m, fmt.Sprintf("%s → %s", match.Rule, match.Sticker.Name))
		default:
			h.reply(ctx, m, fmt.Sprintf("%s ← «%s» → %s", match.Rule, match.Word, match.Sticker.Name))
		}
	}
}

func (h *Handler) reply(ctx context.Context, m *models.Message, text string) {
	_, err := h.api.SendMessage(ctx, &bot.SendMessageParams{ChatID: m.Chat.ID, Text: text})
	if err != nil {
		h.log.Warn("private reply failed", "user_id", m.From.ID, "error", err)
	}
}

// anonymousAdmin reports whether a group admin posted as the group: Telegram then sets From to the
// GroupAnonymousBot, a bot account, and SenderChat to the group itself.
func anonymousAdmin(m *models.Message) bool {
	return m.SenderChat != nil && m.SenderChat.ID == m.Chat.ID
}

// content returns the text a rule reads: the message text, or the caption of a photo or video.
func content(m *models.Message) (string, []models.MessageEntity) {
	if m.Text != "" {
		return m.Text, m.Entities
	}
	return m.Caption, m.CaptionEntities
}

// command returns the command a message starts with, without the @botname suffix.
func command(text string) string {
	fields := strings.Fields(text)
	if len(fields) == 0 || !strings.HasPrefix(fields[0], "/") {
		return ""
	}
	name, _, _ := strings.Cut(fields[0], "@")
	return name
}
