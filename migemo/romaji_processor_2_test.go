package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestRomajiProcessor_2(t *testing.T) {
	processor := migemo.NewRomajiProcessor2()
	testcases := map[string]string{
		"ro-maji": "ろーまじ",
		"atti":    "あっち",
		"att":     "あっt",
		"www":     "wっw",
		"kk":      "っk",
		"n":       "ん",
		"kensaku": "けんさく",
	}
	for k, v := range testcases {
		hiragana := processor.RomajiToHiragana(k)
		if hiragana != v {
			t.Fail()
		}
	}
}

func TestRomajiProcessor2_romajiToHiraganaPredictively_1(t *testing.T) {
	processor := migemo.NewRomajiProcessor2()
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

func TestRomajiProcessor2_romajiToHiraganaPredictively_2(t *testing.T) {
	processor := migemo.NewRomajiProcessor2()
	r := processor.RomajiToHiraganaPredictively("saky")
	if r.Prefix != "さ" {
		t.Error()
	}
	if len(r.Suffixes) != 5 {
		t.Error()
	}
	suffixes := make(map[string]int)
	for _, s := range r.Suffixes {
		suffixes[s] = 0
	}
	containsAll := []string{"きゃ", "きぃ", "きぇ", "きゅ", "きょ"}
	for _, e := range containsAll {
		_, a := suffixes[e]
		if !a {
			t.Error()
		}
	}
}
