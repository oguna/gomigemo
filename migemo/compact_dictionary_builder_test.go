package migemo_test

import (
	"bytes"
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
	if matches[1] != "大分県" || matches[0] != "大阪府" {
		t.Errorf("expected:[大阪府 大分県] actual:%v", matches)
	}
}

func TestLoadFromText2(t *testing.T) {
	text := "けんさく\t検索"
	buff := bytes.NewBufferString(text)
	dict := migemo.BuildDictionaryFromMigemoDictFile(buff)
	matches := make([]string, 0, 2)
	dict.Search(utf16.Encode([]rune("けんさく")), func(word []uint16) {
		matches = append(matches, string(utf16.Decode(word)))
	})
	if len(matches) == 0 || matches[0] != "検索" {
		t.Errorf("expected:[検索] actual:%v", matches)
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
