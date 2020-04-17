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
	trie, _, _ := migemo.BuildLoudsTrie(words)
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
	trie, _, _ := migemo.BuildLoudsTrie(words)
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

func TestLoudsTrieBuilderNew(t *testing.T) {
	words := [][]uint16{}
	words = append(words, utf16.Encode([]rune("baby")))
	words = append(words, utf16.Encode([]rune("bad")))
	words = append(words, utf16.Encode([]rune("bank")))
	words = append(words, utf16.Encode([]rune("box")))
	words = append(words, utf16.Encode([]rune("dad")))
	words = append(words, utf16.Encode([]rune("dance")))
	builder := migemo.NewLoudsTrieBuilder()
	for _, w := range words {
		builder.Add(w)
	}
	var trie = builder.Build()
	if trie.Get(utf16.Encode([]rune("box"))) != 10 {
		t.Error()
	}
}

func TestLoudsTrieBuilderNew_InvalidArgument(t *testing.T) {
	builder := migemo.NewLoudsTrieBuilder()
	var err1 = builder.Add(utf16.Encode([]rune("bad")))
	if err1 != nil {
		t.Error()
	}
	var err2 = builder.Add(utf16.Encode([]rune("bad")))
	if err2 == nil {
		t.Error()
	}
	var err3 = builder.Add(utf16.Encode([]rune("baby")))
	if err3 == nil {
		t.Error()
	}
}
