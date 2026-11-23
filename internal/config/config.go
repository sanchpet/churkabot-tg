// Package config reads the bot's settings from the environment.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds every setting. IDs are numeric: a username can be changed or released and taken by
// someone else, a user ID cannot.
type Config struct {
	Token       string
	Admins      map[int64]bool
	DailyChatID int64
	// TeaseUserID gets the Uimis animation on a TeaseChance share of their group messages; 0 is off.
	TeaseUserID int64
	TeaseChance float64
}

const defaultTeaseChance = 0.03

// Load reads TELEGRAM_BOT_TOKEN, CHURKABOT_ADMIN_IDS (comma-separated user IDs),
// CHURKABOT_DAILY_CHAT_ID, CHURKABOT_TEASE_USER_ID and CHURKABOT_TEASE_CHANCE (0 to 1, default
// 0.03). Only the token is required here; a command that needs more checks it.
func Load() (Config, error) {
	return parse(os.Getenv)
}

func parse(getenv func(string) string) (Config, error) {
	c := Config{Token: strings.TrimSpace(getenv("TELEGRAM_BOT_TOKEN")), Admins: map[int64]bool{}}
	if c.Token == "" {
		return c, errors.New("TELEGRAM_BOT_TOKEN is not set")
	}
	for _, field := range strings.Split(getenv("CHURKABOT_ADMIN_IDS"), ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		id, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return c, fmt.Errorf("CHURKABOT_ADMIN_IDS: %q is not a user ID", field)
		}
		c.Admins[id] = true
	}
	if v := strings.TrimSpace(getenv("CHURKABOT_DAILY_CHAT_ID")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return c, fmt.Errorf("CHURKABOT_DAILY_CHAT_ID: %q is not a chat ID", v)
		}
		c.DailyChatID = id
	}
	if v := strings.TrimSpace(getenv("CHURKABOT_TEASE_USER_ID")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return c, fmt.Errorf("CHURKABOT_TEASE_USER_ID: %q is not a user ID", v)
		}
		c.TeaseUserID = id
	}
	c.TeaseChance = defaultTeaseChance
	if v := strings.TrimSpace(getenv("CHURKABOT_TEASE_CHANCE")); v != "" {
		p, err := strconv.ParseFloat(v, 64)
		if err != nil || p < 0 || p > 1 {
			return c, fmt.Errorf("CHURKABOT_TEASE_CHANCE: %q is not a probability between 0 and 1", v)
		}
		c.TeaseChance = p
	}
	return c, nil
}
