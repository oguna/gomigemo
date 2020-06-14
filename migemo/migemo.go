package migemo

import (
	"regexp"
	"strings"
	"unicode/utf16"
)

// QueryAWord は、migemoクエリを処理する
func QueryAWord(word string, dict *CompactDictionary, operator *RegexOperator) string {
	var utf16word = utf16.Encode([]rune(word))
	var generator = NewTernaryRegexGenerator(*operator)
	generator.Add(utf16word)
	var lower = strings.ToLower(word)
	if dict != nil {
		var utf16lower = utf16.Encode([]rune(lower))
		dict.PredictiveSearch(utf16lower, func(word []uint16) {
			generator.Add(word)
		})
	}
	var zen = ConvertHan2Zen(word)
	generator.Add(utf16.Encode([]rune(zen)))
	var han = ConvertZen2Han(word)
	generator.Add(utf16.Encode([]rune(han)))

	var romajiProcessor = NewRomajiProcessor2()
	var hiraganaResult = romajiProcessor.RomajiToHiraganaPredictively(lower)
	for _, a := range hiraganaResult.Suffixes {
		var hira = hiraganaResult.Prefix + a
		var utf16hira = utf16.Encode([]rune(hira))
		generator.Add(utf16hira)
		if dict != nil {
			dict.PredictiveSearch(utf16hira, func(word []uint16) {
				generator.Add(word)
			})
		}
		var kata = ConvertHira2Kata(string([]rune(utf16.Decode(utf16hira))))
		generator.Add(utf16.Encode([]rune(kata)))
		generator.Add(utf16.Encode([]rune(ConvertZen2Han(kata))))
	}
	return string([]rune(utf16.Decode(generator.Generate())))
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
	// TODO: regexpの処理は遅いため、別の実装に置き換えるべき
	var re = regexp.MustCompile("[^A-Z\\s]+|[A-Z]{2,}|([A-Z][^A-Z\\s]+)|([A-Z]\\s*$)")
	return re.FindAllString(query, -1)
}
