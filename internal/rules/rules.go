// Package rules decides which sticker, if any, answers a message.
package rules

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf16"

	"github.com/go-telegram/bot/models"
)

// Sticker is a sticker the bot can send, by its Telegram file_id.
type Sticker struct {
	Name   string
	FileID string
}

// The file_ids belong to this bot: another bot cannot send them.
var (
	Joker   = Sticker{"joker", "CAACAgIAAxkBAAMxaNYNncvw1laDYnDFYW2N-PBvprEAAlYTAALzWwFIbKJ09BOwAAHkNgQ"}
	Silence = Sticker{"silence", "CAACAgIAAxkBAAOLaNYh6pGe8n_qvV40Hj_SP9vqQJYAAqEVAAINfRBIc9_-Nn_UjbI2BA"}
	Proof   = Sticker{"proof", "CAACAgIAAxkBAAIBl2jiG4Qy5rmnf9zTzIILZbanfSDOAAKXFQAC1MIRSFaAhf9LJ34gNgQ"}
)

// Stickers lists every sticker, for the startup check of their file_ids.
var Stickers = []Sticker{Joker, Silence, Proof}

type rule struct {
	name    string
	sticker Sticker
	word    *regexp.Regexp
}

// Each pattern matches one whole word, lower-cased with ё folded to е. Forms are listed rather than
// left open: an open suffix turns "факт" into "фактор" and "верно" into "наверное".
var table = []rule{
	{"churka", Joker, word(`чур|чурк(а|и|е|у|ой|ам|ами|ах)|чурок|чурочк(а|и|е|у|ой|ам|ами|ах)|чурочек`)},
	{"quiet", Silence, word(`тихо|тише|тишь|тишин(а|ы|е|у|ой)|тишком`)},
	{"shut-up", Silence, word(`молч(ать|у|и|ите|ит|им|ишь|ат|а|ал|ала|али)|заткн(ись|итесь|ешься|етесь|улся|улась|улись)|угомон(ись|итесь|иться)`)},
	{"noise", Silence, word(`шум|шума|шуметь|шумит|шумишь|шуми|шумите`)},
	{"proof", Proof, word(`доказательств(о|а|у|ом|ами|ах)|подтвержден(о|а|ы)?|докаж(и|ите|у|ет|ешь|ем|ут)`)},
	{"fact", Proof, word(`факт(ы|а|у|ом|ов|ам|ами|ах)?`)},
	{"sure", Proof, word(`уверен(а|ы)?`)},
	{"pruf", Proof, word(`пруф(ы|а|у|ом|ов|ам|ами|ах)?|пруфай|пруфани`)},
	{"believe", Proof, word(`верю|веришь|верим|верите|верят|верно`)},
	{"exactly", Proof, word(`точно|точный|точная|точные|точняк`)},
}

func word(alternatives string) *regexp.Regexp {
	return regexp.MustCompile(`^(?:` + alternatives + `)$`)
}

// Match is the verdict on one message: the sticker to send, and the rule and word that chose it.
type Match struct {
	Sticker Sticker
	Rule    string
	Word    string
}

// Find returns the first rule, in table order, that matches a word of text, then the shouting
// check. Words inside links, mentions, commands and code do not count.
func Find(text string, entities []models.MessageEntity) (Match, bool) {
	words := Words(Strip(text, entities))
	for _, r := range table {
		for _, w := range words {
			if r.word.MatchString(normalize(w)) {
				return Match{r.sticker, r.name, w}, true
			}
		}
	}
	if Shouting(words) {
		return Match{Silence, "shouting", ""}, true
	}
	return Match{}, false
}

// Shouting reports whether the upper-case words in a message add up to at least six letters.
// Short Latin words in capitals (API, JSON, HTTPS) are taken for acronyms and single letters are
// ignored, so a technical sentence does not read as a shout.
func Shouting(words []string) bool {
	letters := 0
	for _, w := range words {
		n := len([]rune(w))
		if n < 2 || !upper(w) || (latin(w) && n <= 5) {
			continue
		}
		letters += n
	}
	return letters >= 6
}

// Words splits text into runs of letters.
func Words(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool { return !unicode.IsLetter(r) })
}

// Strip blanks the parts of text covered by entities whose content is not the sender's prose.
// Entity offsets count UTF-16 code units.
func Strip(text string, entities []models.MessageEntity) string {
	units := utf16.Encode([]rune(text))
	for _, e := range entities {
		if !skipped[e.Type] {
			continue
		}
		for i := max(e.Offset, 0); i < e.Offset+e.Length && i < len(units); i++ {
			units[i] = ' '
		}
	}
	return string(utf16.Decode(units))
}

var skipped = map[models.MessageEntityType]bool{
	models.MessageEntityTypeURL:         true,
	models.MessageEntityTypeTextLink:    true,
	models.MessageEntityTypeEmail:       true,
	models.MessageEntityTypeMention:     true,
	models.MessageEntityTypeTextMention: true,
	models.MessageEntityTypeHashtag:     true,
	models.MessageEntityTypeCashtag:     true,
	models.MessageEntityTypeBotCommand:  true,
	models.MessageEntityTypeCode:        true,
	models.MessageEntityTypePre:         true,
}

func normalize(w string) string {
	return strings.ReplaceAll(strings.ToLower(w), "ё", "е")
}

func upper(w string) bool {
	for _, r := range w {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func latin(w string) bool {
	for _, r := range w {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
