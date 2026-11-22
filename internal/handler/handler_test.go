package handler

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/sanchpet/churkabot-tg/internal/rules"
)

type fake struct {
	stickers []*bot.SendStickerParams
	messages []*bot.SendMessageParams
}

func (f *fake) SendSticker(_ context.Context, p *bot.SendStickerParams) (*models.Message, error) {
	f.stickers = append(f.stickers, p)
	return &models.Message{}, nil
}

func (f *fake) SendMessage(_ context.Context, p *bot.SendMessageParams) (*models.Message, error) {
	f.messages = append(f.messages, p)
	return &models.Message{}, nil
}

const admin = 1001

var start = time.Date(2026, 11, 22, 9, 0, 0, 0, time.UTC)

func setup() (*Handler, *fake, *time.Time) {
	f := &fake{}
	now := start
	h := New(f, map[int64]bool{admin: true}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	h.now = func() time.Time { return now }
	return h, f, &now
}

func msg(chat models.ChatType, from int64, text string) *models.Update {
	return &models.Update{Message: &models.Message{
		ID:   7,
		Date: int(start.Unix()),
		Chat: models.Chat{ID: -42, Type: chat},
		From: &models.User{ID: from, FirstName: "Саша"},
		Text: text,
	}}
}

func TestGroupReactsWithReply(t *testing.T) {
	h, f, _ := setup()
	h.Handle(context.Background(), msg(models.ChatTypeSupergroup, 5, "точно?"))
	if len(f.stickers) != 1 {
		t.Fatalf("stickers = %d", len(f.stickers))
	}
	p := f.stickers[0]
	if p.ChatID != int64(-42) || p.ReplyParameters.MessageID != 7 || p.Sticker.(*models.InputFileString).Data != rules.Proof.FileID {
		t.Fatalf("sticker params = %+v", p)
	}
}

func TestGroupReadsCaptions(t *testing.T) {
	h, f, _ := setup()
	u := msg(models.ChatTypeGroup, 5, "")
	u.Message.Caption = "вот пруфы"
	h.Handle(context.Background(), u)
	if len(f.stickers) != 1 {
		t.Fatal("a caption must be read like text")
	}
}

func TestGroupSkipsCaptionEntities(t *testing.T) {
	h, f, _ := setup()
	u := msg(models.ChatTypeGroup, 5, "")
	u.Message.Caption = "#пруфы"
	u.Message.CaptionEntities = []models.MessageEntity{{Type: models.MessageEntityTypeHashtag, Offset: 0, Length: 6}}
	h.Handle(context.Background(), u)
	if len(f.stickers) != 0 {
		t.Fatal("a hashtag in a caption must not trigger")
	}
}

func TestGroupAnswersAnonymousAdmin(t *testing.T) {
	h, f, _ := setup()
	u := msg(models.ChatTypeSupergroup, 1087968824, "точно")
	u.Message.From.IsBot = true
	u.Message.SenderChat = &models.Chat{ID: -42, Type: models.ChatTypeSupergroup}
	h.Handle(context.Background(), u)
	if len(f.stickers) != 1 {
		t.Fatal("an admin posting as the group is a person, not a bot")
	}
}

func TestGroupIgnores(t *testing.T) {
	cases := map[string]func(*models.Update){
		"no trigger":        func(u *models.Update) { u.Message.Text = "привет" },
		"bot sender":        func(u *models.Update) { u.Message.From.IsBot = true },
		"no sender":         func(u *models.Update) { u.Message.From = nil },
		"automatic forward": func(u *models.Update) { u.Message.IsAutomaticForward = true },
		"stale":             func(u *models.Update) { u.Message.Date -= int((11 * time.Minute).Seconds()) },
		"channel":           func(u *models.Update) { u.Message.Chat.Type = models.ChatTypeChannel },
	}
	for name, mutate := range cases {
		h, f, _ := setup()
		u := msg(models.ChatTypeGroup, 5, "точно")
		mutate(u)
		h.Handle(context.Background(), u)
		if len(f.stickers)+len(f.messages) != 0 {
			t.Errorf("%s: the bot answered", name)
		}
	}
}

func TestGroupCooldownIsPerChat(t *testing.T) {
	h, f, now := setup()
	h.Handle(context.Background(), msg(models.ChatTypeGroup, 5, "точно"))
	h.Handle(context.Background(), msg(models.ChatTypeGroup, 5, "точно"))
	other := msg(models.ChatTypeGroup, 5, "точно")
	other.Message.Chat.ID = -43
	h.Handle(context.Background(), other)
	if len(f.stickers) != 2 {
		t.Fatalf("within the cooldown: stickers = %d, want 2 (one per chat)", len(f.stickers))
	}
	*now = now.Add(cooldown)
	later := msg(models.ChatTypeGroup, 5, "точно")
	later.Message.Date = int(now.Unix())
	h.Handle(context.Background(), later)
	if len(f.stickers) != 3 {
		t.Fatal("after the cooldown the chat must get a sticker again")
	}
}

func TestPrivateNonAdmin(t *testing.T) {
	h, f, _ := setup()
	h.Handle(context.Background(), msg(models.ChatTypePrivate, 5, "точно"))
	u := msg(models.ChatTypePrivate, 5, "")
	u.Message.Sticker = &models.Sticker{FileID: "x"}
	h.Handle(context.Background(), u)
	if len(f.messages)+len(f.stickers) != 0 {
		t.Fatal("a non-admin gets nothing but the /start line")
	}
	h.Handle(context.Background(), msg(models.ChatTypePrivate, 5, "/start"))
	if len(f.messages) != 1 {
		t.Fatal("a non-admin's /start gets one line")
	}
}

func TestPrivateAdmin(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*models.Message)
		want   string
	}{
		{"start", func(m *models.Message) { m.Text = "/start@churkabot" }, "Привет, Саша"},
		{"diagnose", func(m *models.Message) { m.Text = "наверное, это факт" }, "fact ← «факт» → proof"},
		{"no match", func(m *models.Message) { m.Text = "наверное" }, "no match"},
		{"shouting", func(m *models.Message) { m.Text = "ЗАЧЕМ ТАК" }, "shouting → silence"},
		{"sticker", func(m *models.Message) { m.Sticker = &models.Sticker{FileID: "stk", SetName: "set"} }, "file_id: stk"},
		{"largest photo", func(m *models.Message) {
			m.Photo = []models.PhotoSize{{FileID: "small", Width: 90, Height: 90}, {FileID: "big", Width: 1280, Height: 960}, {FileID: "mid", Width: 320, Height: 240}}
		}, "file_id: big\nsize: 1280x960"},
	}
	for _, c := range cases {
		h, f, _ := setup()
		u := msg(models.ChatTypePrivate, admin, "")
		c.mutate(u.Message)
		h.Handle(context.Background(), u)
		if len(f.stickers) != 0 || len(f.messages) != 1 || !strings.Contains(f.messages[0].Text, c.want) {
			t.Errorf("%s: messages = %+v, want one containing %q", c.name, f.messages, c.want)
		}
	}
}

// The handler's context is cancelled on shutdown; the reply in flight must still go out.
func TestReplySurvivesCancellation(t *testing.T) {
	h, f, _ := setup()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	h.api = &ctxCheck{fake: f, t: t}
	h.Handle(ctx, msg(models.ChatTypeGroup, 5, "точно"))
	if len(f.stickers) != 1 {
		t.Fatal("no sticker sent")
	}
}

type ctxCheck struct {
	*fake
	t *testing.T
}

func (c *ctxCheck) SendSticker(ctx context.Context, p *bot.SendStickerParams) (*models.Message, error) {
	if ctx.Err() != nil {
		c.t.Fatal("send got a cancelled context")
	}
	return c.fake.SendSticker(ctx, p)
}
