package migemo_test

import (
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestTernaryRegexGenerator_1(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("bad")))
	generator.Add(utf16.Encode([]rune("dad")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "(bad|dad)"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_2(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("bad")))
	generator.Add(utf16.Encode([]rune("bat")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "ba[dt]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_3(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("a")))
	generator.Add(utf16.Encode([]rune("b")))
	generator.Add(utf16.Encode([]rune("a")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "[ab]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_4(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("a.b")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "a\\.b"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_5(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("abc")))
	generator.Add(utf16.Encode([]rune("abcd")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "abc"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_6(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(*regex_operator)
	generator.Add(utf16.Encode([]rune("abcd")))
	generator.Add(utf16.Encode([]rune("abc")))
	result := string(utf16.Decode(generator.Generate()))
	expect := "abc"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}
