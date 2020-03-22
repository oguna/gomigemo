package migemo

import (
	"testing"
	"unicode/utf16"
	"github.com/oguna/gomigemo/migemo"
)

func TestRegexGenerator_1(t *testing.T) {
	generator := migemo.NewRegexGenerator()
	generator.Add(utf16.Encode([]rune("bad")))
	generator.Add(utf16.Encode([]rune("dad")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "(bad|dad)"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestRegexGenerator_2(t *testing.T) {
	generator := migemo.NewRegexGenerator()
	generator.Add(utf16.Encode([]rune("bad")))
	generator.Add(utf16.Encode([]rune("bat")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "ba[dt]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestRegexGenerator_3(t *testing.T) {
	generator := migemo.NewRegexGenerator()
	generator.Add(utf16.Encode([]rune("a")))
	generator.Add(utf16.Encode([]rune("b")))
	generator.Add(utf16.Encode([]rune("a")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "[ab]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestRegexGenerator_4(t *testing.T) {
	generator := migemo.NewRegexGenerator()
	generator.Add(utf16.Encode([]rune("a.b")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "a\\.b"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}