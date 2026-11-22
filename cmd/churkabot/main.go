// Command churkabot is a Telegram bot that answers trigger words in group chats with stickers.
//
//	churkabot poll   long-polls Telegram and answers messages (the default)
//	churkabot daily  sends the daily sticker once and exits
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/sanchpet/churkabot-tg/internal/config"
	"github.com/sanchpet/churkabot-tg/internal/daily"
	"github.com/sanchpet/churkabot-tg/internal/handler"
	"github.com/sanchpet/churkabot-tg/internal/rules"
)

// Exit codes: 1 is worth retrying, 2 is not (bad configuration, revoked token, bot kicked).
const (
	exitRetry     = 1
	exitPermanent = 2
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cmd := "poll"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	cfg, err := config.Load()
	if err != nil {
		log.Error("configuration", "error", err)
		os.Exit(exitPermanent)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	switch cmd {
	case "poll":
		err = poll(ctx, cfg, log)
	case "daily":
		err = sendDaily(ctx, cfg, log)
	default:
		err = fmt.Errorf("%w: unknown command %q, want poll or daily", daily.ErrPermanent, cmd)
	}
	if err != nil {
		log.Error(cmd+" failed", "error", err)
		stop()
		if errors.Is(err, daily.ErrPermanent) || errors.Is(err, bot.ErrorUnauthorized) {
			os.Exit(exitPermanent)
		}
		os.Exit(exitRetry)
	}
}

func poll(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	var h *handler.Handler
	b, err := bot.New(cfg.Token,
		bot.WithSkipGetMe(),
		// One update at a time, confirmed only after its handler returns: at-least-once delivery,
		// and on shutdown Start waits for the reply in flight.
		bot.WithNotAsyncHandlers(),
		bot.WithWorkers(1),
		bot.WithAllowedUpdates(bot.AllowedUpdates{models.AllowedUpdateMessage}),
		bot.WithErrorsHandler(func(err error) { log.Warn("polling", "error", err) }),
		bot.WithDefaultHandler(func(ctx context.Context, _ *bot.Bot, u *models.Update) { h.Handle(ctx, u) }),
	)
	if err != nil {
		return err
	}
	h = handler.New(b, cfg.Admins, log)

	me, err := b.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("getMe: %w", err)
	}
	if !me.CanReadAllGroupMessages {
		log.Warn("privacy mode is on: in groups the bot sees only commands and replies to it; turn it off in BotFather")
	}
	if len(cfg.Admins) == 0 {
		log.Warn("CHURKABOT_ADMIN_IDS is empty: nobody can use the private-chat helpers")
	}
	// Polling fails with 409 while a webhook is set.
	if _, err := b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{}); err != nil {
		log.Warn("deleteWebhook", "error", err)
	}
	// A sticker whose file_id no longer works is found now rather than on its first trigger.
	for _, s := range rules.Stickers {
		if _, err := b.GetFile(ctx, &bot.GetFileParams{FileID: s.FileID}); err != nil {
			log.Error("sticker file_id is not usable", "sticker", s.Name, "error", err)
		}
	}

	log.Info("polling", "bot", me.Username, "admins", len(cfg.Admins))
	b.Start(ctx)
	log.Info("stopped")
	return nil
}

func sendDaily(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	if cfg.DailyChatID == 0 {
		return fmt.Errorf("%w: CHURKABOT_DAILY_CHAT_ID is not set", daily.ErrPermanent)
	}
	b, err := bot.New(cfg.Token, bot.WithSkipGetMe())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	admins := make([]int64, 0, len(cfg.Admins))
	for id := range cfg.Admins {
		admins = append(admins, id)
	}
	return daily.Send(ctx, b, cfg.DailyChatID, admins, log)
}
