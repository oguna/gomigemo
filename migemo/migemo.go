package migemo

import (
	"strings"
	"unicode/utf16"
)

// QueryAWord は、migemoクエリを処理する
func QueryAWord(word string, dict *CompactDictionary, operator *RegexOperator) string {
	var utf32word = []rune(word)
	var generator = NewTernaryRegexGenerator(operator)
	generator.Add(utf32word)
	var lower = strings.ToLower(word)
	if dict != nil {
		var utf16lower = utf16.Encode([]rune(lower))
		dict.PredictiveSearch(utf16lower, func(word []uint16) {
			generator.Add(utf16.Decode(word))
		})
	}
	var zen = ConvertHan2Zen(word)
	generator.Add([]rune(zen))
	var han = ConvertZen2Han(word)
	generator.Add([]rune(han))

	var romajiProcessor = NewRomajiProcessor2()
	var hiraganaResult = romajiProcessor.RomajiToHiraganaPredictively(lower)
	for _, a := range hiraganaResult.Suffixes {
		var hira = hiraganaResult.Prefix + a
		var utf32hira = []rune(hira)
		var utf16hira = utf16.Encode([]rune(hira))
		generator.Add(utf32hira)
		if dict != nil {
			dict.PredictiveSearch(utf16hira, func(word []uint16) {
				generator.Add(utf16.Decode(word))
			})
		}
		var kata = ConvertHira2Kata(string([]rune(utf16.Decode(utf16hira))))
		generator.Add([]rune(kata))
		generator.Add([]rune(ConvertZen2Han(kata)))
	}
	return string([]rune(generator.Generate()))
}

// Query は、migemoクエリを処理する
func Query(word string, dict *CompactDictionary, operator *RegexOperator) string {
	if len(word) == 0 {
		return ""
	}
	words := parseQuery(word)
	results := make([]string, len(words))
	for i, w := range words {
		results[i] = QueryAWord(w, dict, operator)
	}
	return strings.Join(results, "")
}

func parseQuery(query string) []string {
	parser := NewMigemoParser(query)
	words := make([]string, 0, 8)
	for true {
		w := parser.Next()
		if w == "" {
			break
		}
		words = append(words, w)
	}
	return words
}
