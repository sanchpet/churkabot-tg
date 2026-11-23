package daily

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// script answers each sendSticker call with the next error in errs, nil once they run out.
type script struct {
	errs     []error
	stickers []int64
	messages []int64
}

func (s *script) SendSticker(_ context.Context, p *bot.SendStickerParams) (*models.Message, error) {
	s.stickers = append(s.stickers, p.ChatID.(int64))
	if len(s.errs) == 0 {
		return &models.Message{}, nil
	}
	err := s.errs[0]
	s.errs = s.errs[1:]
	return nil, err
}

func (s *script) SendMessage(_ context.Context, p *bot.SendMessageParams) (*models.Message, error) {
	s.messages = append(s.messages, p.ChatID.(int64))
	// A failed notification must not turn a sent sticker into a failed job.
	return nil, errors.New("blocked by user")
}

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestSend(t *testing.T) {
	cases := []struct {
		name      string
		errs      []error
		stickers  []int64
		notified  int
		permanent bool
		fails     bool
	}{
		{name: "ok", stickers: []int64{-1}},
		{name: "migrated", errs: []error{&bot.MigrateError{MigrateToChatID: -1002}}, stickers: []int64{-1, -1002}, notified: 1},
		{name: "rate limited", errs: []error{&bot.TooManyRequestsError{RetryAfter: 0}}, stickers: []int64{-1, -1}},
		{name: "kicked", errs: []error{fmt.Errorf("%w, bot was kicked", bot.ErrorForbidden)}, stickers: []int64{-1}, notified: 1, permanent: true, fails: true},
		{name: "revoked token", errs: []error{fmt.Errorf("%w, Unauthorized", bot.ErrorUnauthorized)}, stickers: []int64{-1}, permanent: true, fails: true},
		{name: "network", errs: []error{errors.New("connection reset")}, stickers: []int64{-1}, fails: true},
	}
	for _, c := range cases {
		s := &script{errs: c.errs}
		err := Send(context.Background(), s, -1, []int64{1001}, discard)
		if (err != nil) != c.fails || errors.Is(err, ErrPermanent) != c.permanent {
			t.Errorf("%s: err = %v", c.name, err)
		}
		if fmt.Sprint(s.stickers) != fmt.Sprint(c.stickers) || len(s.messages) != c.notified {
			t.Errorf("%s: stickers to %v, %d notifications", c.name, s.stickers, len(s.messages))
		}
	}
}
