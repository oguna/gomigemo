package migemo_test

import (
	"io"
	"os"
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestCompactDictionary_1(t *testing.T) {
	f, err := os.Open("../testdata/todofuken-dict")
	if err != nil {
		t.Fatalf("failed to open dict: %v", err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read dict: %v", err)
	}
	dict := migemo.NewCompactDictionary(buf)
	list := []string{}
	fn := func(s []uint16) {
		list = append(list, string(utf16.Decode(s)))
	}
	dict.Search(utf16.Encode([]rune("とうきょうと")), fn)
	for _, w := range list {
		if w == "東京都" {
			return
		}
	}
	t.Error()
}
func TestCompactDictionary_2(t *testing.T) {
	f, err := os.Open("../testdata/todofuken-dict")
	if err != nil {
		t.Fatalf("failed to open dict: %v", err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read dict: %v", err)
	}
	dict := migemo.NewCompactDictionary(buf)
	list := []string{}
	fn := func(s []uint16) {
		list = append(list, string(utf16.Decode(s)))
	}
	dict.PredictiveSearch(utf16.Encode([]rune("か")), fn)
	for _, w := range list {
		if w == "神奈川県" {
			return
		}
	}
	t.Error()
}
