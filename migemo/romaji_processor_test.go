package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestRomajiProcessor_1(t *testing.T) {
	processor := migemo.NewRomajiProcessor()
	hiragana := processor.RomajiToHiragana("ro-maji")
	if hiragana != "ろーまじ" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("atti")
	if hiragana != "あっち" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("att")
	if hiragana != "あっt" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("www")
	if hiragana != "wっw" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("kk")
	if hiragana != "っk" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("n")
	if hiragana != "ん" {
		t.Error()
	}
	hiragana = processor.RomajiToHiragana("kensaku")
	if hiragana != "けんさく" {
		t.Error()
	}
}

func TestRomajiProcessor_romajiToHiraganaPredictively_1(t *testing.T) {
	processor := migemo.NewRomajiProcessor()
	r := processor.RomajiToHiraganaPredictively("kiku")
	if r.Prefix != "きく" {
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
	r := processor.RomajiToHiraganaPredictively("ky")
	if r.Prefix != "" {
		t.Error()
	}
	if len(r.Suffixes) != 5 {
		t.Error()
	}
	suffixes := make(map[string]int)
	for _, s := range r.Suffixes {
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
