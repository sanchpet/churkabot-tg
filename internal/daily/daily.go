// Package daily sends the morning sticker. It runs as a one-shot job, so the scheduler, not the
// bot, remembers whether today's run happened.
package daily

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/sanchpet/churkabot-tg/internal/rules"
)

// ErrPermanent marks a failure that a retry will not fix: the job should fail, not run again.
var ErrPermanent = errors.New("permanent failure")

// Sender is the part of the Bot API the job uses.
type Sender interface {
	SendSticker(ctx context.Context, params *bot.SendStickerParams) (*models.Message, error)
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}

// Send posts the daily sticker to chatID. Once the sticker is out it returns nil whatever happens
// next: a non-zero exit after a successful send would make the job retry and post it twice.
// Failures an admin has to act on are also sent to the admins in private.
func Send(ctx context.Context, api Sender, chatID int64, admins []int64, log *slog.Logger) error {
	err := sendSticker(ctx, api, chatID)

	var migrated *bot.MigrateError
	if errors.As(err, &migrated) {
		newID := int64(migrated.MigrateToChatID)
		log.Error("chat became a supergroup; set CHURKABOT_DAILY_CHAT_ID to the new ID", "chat_id", chatID, "new_chat_id", newID)
		notify(ctx, api, admins, log, fmt.Sprintf("Чат %d стал супергруппой %d: обнови CHURKABOT_DAILY_CHAT_ID.", chatID, newID))
		chatID, err = newID, sendSticker(ctx, api, newID)
	}

	switch {
	case err == nil:
		log.Info("daily sticker sent", "chat_id", chatID, "sticker", rules.Joker.Name)
		return nil
	case errors.Is(err, bot.ErrorUnauthorized):
		return fmt.Errorf("%w: %w", ErrPermanent, err)
	case errors.Is(err, bot.ErrorForbidden), errors.Is(err, bot.ErrorBadRequest):
		notify(ctx, api, admins, log, fmt.Sprintf("Ежедневный стикер в чат %d не ушёл: %v", chatID, err))
		return fmt.Errorf("%w: %w", ErrPermanent, err)
	default:
		return err
	}
}

// sendSticker sends once, and once more after the wait Telegram asks for on a 429.
func sendSticker(ctx context.Context, api Sender, chatID int64) error {
	params := &bot.SendStickerParams{ChatID: chatID, Sticker: &models.InputFileString{Data: rules.Joker.FileID}}
	_, err := api.SendSticker(ctx, params)
	var limited *bot.TooManyRequestsError
	if !errors.As(err, &limited) {
		return err
	}
	select {
	case <-time.After(time.Duration(limited.RetryAfter) * time.Second):
	case <-ctx.Done():
		return err
	}
	_, err = api.SendSticker(ctx, params)
	return err
}

func notify(ctx context.Context, api Sender, admins []int64, log *slog.Logger, text string) {
	for _, id := range admins {
		if _, err := api.SendMessage(ctx, &bot.SendMessageParams{ChatID: id, Text: text}); err != nil {
			log.Warn("admin notification failed", "user_id", id, "error", err)
		}
	}
}
