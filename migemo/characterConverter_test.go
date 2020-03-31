package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestCharacterConverter_ConvertHira2Kata(t *testing.T) {
	expected := "ア"
	actual := migemo.ConvertHira2Kata("あ")
	if actual != expected {
		t.Error("result: ", actual, "\nexpected: ", expected, "\n")
	}
}

func TestCharacterConverter_ConvertHan2Zen(t *testing.T) {
	expected := "ア"
	actual := migemo.ConvertHan2Zen("ｱ")
	if actual != expected {
		t.Error("result: ", actual, "\nexpected: ", expected, "\n")
	}
}

func TestCharacterConverter_ConvertZen2Han(t *testing.T) {
	expected := "ｱ"
	actual := migemo.ConvertZen2Han("ア")
	if actual != expected {
		t.Error("result: ", actual, "\nexpected: ", expected, "\n")
	}
}
