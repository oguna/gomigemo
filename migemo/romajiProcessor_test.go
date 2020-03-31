package migemo_test

import (
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestRomajiProcessor_1(t *testing.T) {
	processor := migemo.NewRomajiProcessor()
	hiragana := string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("ro-maji")))))
	if hiragana != "ろーまじ" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("atti")))))
	if hiragana != "あっち" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("att")))))
	if hiragana != "あっt" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("www")))))
	if hiragana != "wっw" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("kk")))))
	if hiragana != "っk" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("n")))))
	if hiragana != "ん" {
		t.Error()
	}
	hiragana = string(utf16.Decode(processor.RomajiToHiragana(utf16.Encode([]rune("kensaku")))))
	if hiragana != "けんさく" {
		t.Error()
	}
}

func TestRomajiProcessor_romajiToHiraganaPredictively_1(t *testing.T) {
	processor := migemo.NewRomajiProcessor()
	r := processor.RomajiToHiraganaPredictively(utf16.Encode([]rune("kiku")))
	prefix := string(utf16.Decode(r.Prefix))
	if prefix != "きく" {
		t.Error()
	}
	if len(r.Suffixes) != 1 {
		t.Error()
	}
	if len(r.Suffixes[0]) != 0 {
		t.Error()
	}
}

func TestRomajiProcessor_romajiToHiraganaPredictively_2(t *testing.T) {
	processor := migemo.NewRomajiProcessor()
	r := processor.RomajiToHiraganaPredictively(utf16.Encode([]rune("ky")))
	prefix := string(utf16.Decode(r.Prefix))
	if prefix != "" {
		t.Error()
	}
	if len(r.Suffixes) != 5 {
		t.Error()
	}
	suffixes := make(map[string]int)
	for _, s := range r.Suffixes {
		s := string(utf16.Decode(s))
		suffixes[s] = 0
	}
	_, kya := suffixes["きゃ"]
	if !kya {
		t.Error()
	}
	_, kixi := suffixes["きぃ"]
	if !kixi {
		t.Error()
	}
	_, kixe := suffixes["きぇ"]
	if !kixe {
		t.Error()
	}
	_, kixyu := suffixes["きゅ"]
	if !kixyu {
		t.Error()
	}
	_, kixyo := suffixes["きょ"]
	if !kixyo {
		t.Error()
	}
}
