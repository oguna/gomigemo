package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestDoubleArrayBuilder_test1(t *testing.T) {
	keys := []string{"ab", "ac", "b", "da"}
	for i := range keys {
		chars := []rune(keys[i])
		for j := range chars {
			chars[j] -= 'a' - 1
		}
		keys[i] = string(chars)
	}
	indices := []int16{5, 6, 2, 3}
	trie := migemo.BuildDoubleArray(keys)
	for i := range keys {
		if trie.Lookup(keys[i]) != indices[i] {
			t.Error()
		}
	}
	if trie.Lookup("") != 0 {
		t.Error()
	}
}

func TestDoubleArrayBuilder_test2(t *testing.T) {
	keys := []string{"ab", "ac", "b", "da"}
	trie := migemo.BuildDoubleArray(keys)
	for i := range keys {
		if trie.Lookup(keys[i]) <= 0 {
			t.Error()
		}
	}
	if trie.Lookup("c") != -1 {
		t.Error()
	}
}
