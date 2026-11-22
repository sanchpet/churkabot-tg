package rules

import (
	"strings"
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestFindFires(t *testing.T) {
	cases := map[string]string{
		"ну ты и чурка":            "churka",
		"Чур меня":                 "churka",
		"чурочек мой":              "churka",
		"Тише, пожалуйста":         "quiet",
		"тишина в зале":            "quiet",
		"молчи уже":                "shut-up",
		"Заткнись!":                "shut-up",
		"угомонитесь":              "shut-up",
		"не шуми":                  "noise",
		"а доказательства где?":    "proof",
		"подтверждено":             "proof",
		"докажи":                   "proof",
		"это факт":                 "fact",
		"Факты такие":              "fact",
		"ты уверен?":               "sure",
		"пруфы будут?":             "pruf",
		"не верю":                  "believe",
		"Верно":                    "believe",
		"точно?":                   "exactly",
		"точняк":                   "exactly",
		"ну ТИШЕ":                  "quiet",
		"Тишком пробрался":         "quiet",
		"факт, а не чурка":         "churka", // table order, not word order
		"ЗАТКНИТЕСЬ ВСЕ":           "shut-up",
		"ХВАТИТ УЖЕ":               "shouting",
		"ПОЧЕМУ":                   "shouting",
		"ну КОГДА УЖЕ":             "shouting",
		"Я СКАЗАЛ НЕТ":             "shouting",
		"молчать, тише!":           "quiet",
		"подтверждена гипотеза":    "proof",
		"я уверена":                "sure",
		"пруфани":                  "pruf",
		"пруфай давай":             "pruf",
		"ты заткнёшься?":           "shut-up",
		"я молчу":                  "shut-up",
		"тишь да гладь":            "quiet",
		"всё верно, факт":          "fact",
		"точная копия":             "exactly",
		"KUBERNETES IS DOWN AGAIN": "shouting",
	}
	for text, want := range cases {
		m, ok := Find(text, nil)
		if !ok || m.Rule != want {
			t.Errorf("Find(%q) = %q, %v; want %q", text, m.Rule, ok, want)
		}
	}
}

// Each of these fired a sticker in the Java bot, which matched roots as substrings.
func TestFindIgnoresWordsThatOnlyContainARoot(t *testing.T) {
	for _, text := range []string{
		"наверное", "наверно", "проверю", "достаточно", "восточно", "неточность", "точность",
		"этих", "стихи", "стихия", "ответишь", "платишь", "по умолчанию", "молчаливый",
		"артефакт", "фактически", "фактура", "фактор", "чурчхела", "чураться", "шумахер",
		"верность", "уверенность", "тихоокеанский", "тихий",
		"подними HTTPS", "AWS EKS", "JSON API", "США И РФ", "это ОЧЕНЬ важно", "ок",
		"", "😀😀😀",
	} {
		if m, ok := Find(text, nil); ok {
			t.Errorf("Find(%q) fired %q on %q", text, m.Rule, m.Word)
		}
	}
}

func TestFindSkipsEntities(t *testing.T) {
	// "тише" sits inside code, "факт" inside a link, "чурка" inside a mention.
	text := "смотри `тише` тут https://x.example/факт @чурка"
	entities := []models.MessageEntity{
		{Type: models.MessageEntityTypeCode, Offset: 7, Length: 6},
		{Type: models.MessageEntityTypeURL, Offset: 18, Length: 22},
		{Type: models.MessageEntityTypeMention, Offset: 41, Length: 6},
	}
	if _, ok := Find(text, nil); !ok {
		t.Fatal("without entities the text must match, or this test proves nothing")
	}
	if got := Strip(text, entities); got != "смотри"+strings.Repeat(" ", 8)+"тут"+strings.Repeat(" ", 30) {
		t.Fatalf("Strip = %q", got)
	}
	if m, ok := Find(text, entities); ok {
		t.Fatalf("fired %q on %q in %q", m.Rule, m.Word, Strip(text, entities))
	}
	// Bold is the sender's prose and still counts.
	bold := []models.MessageEntity{{Type: models.MessageEntityTypeBold, Offset: 0, Length: 5}}
	if _, ok := Find("точно", bold); !ok {
		t.Fatal("bold text must still match")
	}
}

func TestStripCountsUTF16(t *testing.T) {
	// The emoji takes two UTF-16 units, so the code entity starts at 3, not 2.
	text := "😀 `факт` верно"
	got := Strip(text, []models.MessageEntity{{Type: models.MessageEntityTypeCode, Offset: 3, Length: 6}})
	if want := "😀        верно"; got != want {
		t.Fatalf("Strip = %q, want %q", got, want)
	}
}

func TestStripToleratesBadOffsets(t *testing.T) {
	got := Strip("точно", []models.MessageEntity{{Type: models.MessageEntityTypeCode, Offset: 3, Length: 100}})
	if got != "точ  " {
		t.Fatalf("Strip = %q", got)
	}
}
