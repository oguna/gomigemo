package migemo_test

import (
	"os"
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestLoadFromText(t *testing.T) {
	f, err := os.Open("../testdata/todofuken.txt")
	if err != nil {
		panic(err)
	}
	dict := migemo.BuildDictionaryFromMigemoDictFile(f)
	matches := make([]string, 0, 2)
	dict.PredictiveSearch(utf16.Encode([]rune("おお")), func(word []uint16) {
		matches = append(matches, string(utf16.Decode(word)))
	})
	if matches[0] != "大分県" || matches[1] != "大阪府" {
		t.Error("")
	}
}

func TestCompactDictionaryBuilder_ExtractTail(t *testing.T) {
	words := []string{"a", "aaa", "b", "cc"}
	tails := migemo.ExtractTail(words)
	expectedTails := []uint32{0, 1, 0, 1}
	if len(tails) != len(expectedTails) {
		t.Error()
	}
	for i := 0; i < len(tails); i++ {
		if tails[i] != expectedTails[i] {
			t.Fatalf("#%d expected:%d actual:%d", i, expectedTails[i], tails[i])
		}
	}
}
