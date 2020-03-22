package migemo

import "regexp"
import "strings"
import "unicode/utf16"

func QueryAWord(word string, dict *CompactDictionary) string {
    var utf16word = utf16.Encode([]rune(word))
    var generator = NewRegexGenerator()
    generator.Add(utf16word)
    var lower = utf16.Encode([]rune(strings.ToLower(word)))
    if dict != nil {
        dict.PredictiveSearch(lower, func(word []uint16) {
            generator.Add(word)
        })
    }
    var zen = ConvertHan2Zen(word)
    generator.Add(utf16.Encode([]rune(zen)))
    var han = ConvertZen2Han(word)
    generator.Add(utf16.Encode([]rune(han)))

    var romajiProcessor = NewRomajiProcessor()
    var hiraganaResult = romajiProcessor.RomajiToHiraganaPredictively(lower)
    for _, a := range hiraganaResult.Suffixes {
        var hira = append(hiraganaResult.Prefix, a...)
        generator.Add(hira)
        if dict != nil {
            dict.PredictiveSearch(hira, func(word []uint16) {
                generator.Add(word)
            })
        }
        var kata = ConvertHira2Kata(string([]rune(utf16.Decode(hira))))
        generator.Add(utf16.Encode([]rune(kata)))
        generator.Add(utf16.Encode([]rune(ConvertZen2Han(kata))))
    }
    return string([]rune(utf16.Decode(generator.Generate())))
}

func Query(word string, dict *CompactDictionary) string {
    if len(word) == 0 {
        return ""
    }
    words := parseQuery(word)
    results := make([]string, len(words))
    for i, w := range words {
        results[i] = QueryAWord(w, dict)
    }
    return strings.Join(results, "")
}

func parseQuery(query string) []string {
    // TODO: regexpの処理は遅いため、別の実装に置き換えるべき
    var re = regexp.MustCompile("[^A-Z\\s]+|[A-Z]{2,}|([A-Z][^A-Z\\s]+)|([A-Z]\\s*$)");
    return re.FindAllString(query, -1)
}