package migemo_test

import (
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestLoudsTrieBuilder_1(t *testing.T) {
	words := [][]uint16{}
	words = append(words, utf16.Encode([]rune("baby")))
	words = append(words, utf16.Encode([]rune("bad")))
	words = append(words, utf16.Encode([]rune("bank")))
	words = append(words, utf16.Encode([]rune("box")))
	words = append(words, utf16.Encode([]rune("dad")))
	words = append(words, utf16.Encode([]rune("dance")))
	trie := migemo.BuildLoudsTrie(words, nil)
	if trie.Get(utf16.Encode([]rune("box"))) != 10 {
		t.Error()
	}
}

func TestLoudsTrieBuilder_2(t *testing.T) {
	words := [][]uint16{}
	words = append(words, utf16.Encode([]rune("a")))
	words = append(words, utf16.Encode([]rune("aa")))
	words = append(words, utf16.Encode([]rune("ab")))
	words = append(words, utf16.Encode([]rune("bb")))
	trie := migemo.BuildLoudsTrie(words, nil)
	if trie.Get(utf16.Encode([]rune(""))) != 1 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("a"))) != 2 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("b"))) != 3 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("aa"))) != 4 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("ab"))) != 5 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("bb"))) != 6 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("bbb"))) != -1 {
		t.Error()
	}
	if trie.Get(utf16.Encode([]rune("c"))) != -1 {
		t.Error()
	}
}
