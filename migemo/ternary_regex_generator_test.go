package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestTernaryRegexGenerator_1(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("bad"))
	generator.Add([]rune("dad"))
	result := string(generator.Generate())
	expect := "(bad|dad)"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_2(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("bad"))
	generator.Add([]rune("bat"))
	result := string(generator.Generate())
	expect := "ba[dt]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_3(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("a"))
	generator.Add([]rune("b"))
	generator.Add([]rune("a"))
	result := string(generator.Generate())
	expect := "[ab]"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_4(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("a.b"))
	result := string(generator.Generate())
	expect := "a\\.b"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_5(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("abc"))
	generator.Add([]rune("abcd"))
	result := string(generator.Generate())
	expect := "abc"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}

func TestTernaryRegexGenerator_6(t *testing.T) {
	regex_operator := migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	generator := migemo.NewTernaryRegexGenerator(regex_operator)
	generator.Add([]rune("abcd"))
	generator.Add([]rune("abc"))
	result := string(generator.Generate())
	expect := "abc"
	if result != expect {
		t.Error("result: ", result, "\nexpected: ", expect, "\n")
	}
}
