package config

import "testing"

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestParse(t *testing.T) {
	c, err := parse(env(map[string]string{
		"TELEGRAM_BOT_TOKEN":      " 123:abc ",
		"CHURKABOT_ADMIN_IDS":     "1001, 1002,",
		"CHURKABOT_DAILY_CHAT_ID": "-1009999999999",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if c.Token != "123:abc" || !c.Admins[1001] || !c.Admins[1002] || len(c.Admins) != 2 || c.DailyChatID != -1009999999999 {
		t.Fatalf("parse = %+v", c)
	}
}

func TestParseRejects(t *testing.T) {
	for name, m := range map[string]map[string]string{
		"no token":      {},
		"blank token":   {"TELEGRAM_BOT_TOKEN": "  "},
		"username":      {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_ADMIN_IDS": "sanchpet"},
		"chat username": {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_DAILY_CHAT_ID": "@memes"},
	} {
		if _, err := parse(env(m)); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}
