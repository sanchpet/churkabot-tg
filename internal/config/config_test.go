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
	if c.TeaseUserID != 0 || c.TeaseChance != 0.03 {
		t.Fatalf("tease defaults = %d, %v", c.TeaseUserID, c.TeaseChance)
	}
}

func TestParseTease(t *testing.T) {
	c, err := parse(env(map[string]string{
		"TELEGRAM_BOT_TOKEN":      "t",
		"CHURKABOT_TEASE_USER_ID": "1003",
		"CHURKABOT_TEASE_CHANCE":  "0.5",
	}))
	if err != nil || c.TeaseUserID != 1003 || c.TeaseChance != 0.5 {
		t.Fatalf("parse = %+v, %v", c, err)
	}
}

func TestParseRejects(t *testing.T) {
	for name, m := range map[string]map[string]string{
		"no token":       {},
		"blank token":    {"TELEGRAM_BOT_TOKEN": "  "},
		"username":       {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_ADMIN_IDS": "sanchpet"},
		"chat username":  {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_DAILY_CHAT_ID": "@memes"},
		"tease username": {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_TEASE_USER_ID": "@someone"},
		"percent chance": {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_TEASE_CHANCE": "3"},
		"word chance":    {"TELEGRAM_BOT_TOKEN": "t", "CHURKABOT_TEASE_CHANCE": "rarely"},
	} {
		if _, err := parse(env(m)); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}
