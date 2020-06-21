package migemo_test

import (
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestLoudsTrieBuilder_1(t *testing.T) {
	words := []string{"baby", "bad", "bank", "box", "dad", "dance"}
	keys := make([][]uint16, len(words))
	for i := 0; i < len(words); i++ {
		keys[i] = utf16.Encode([]rune(words[i]))
	}
	trie, nodes, _ := migemo.BuildLoudsTrie(keys)
	if trie.Lookup(utf16.Encode([]rune("box"))) != 10 {
		t.Error()
	}
	var expectedNodes = []uint32{13, 8, 14, 10, 11, 16}
	for i := 0; i < len(nodes); i++ {
		if nodes[i] != expectedNodes[i] {
			t.Error(words[i])
		}
	}
	for i := 0; i < len(nodes); i++ {
		if uint32(trie.Lookup(keys[i])) != expectedNodes[i] {
			t.Error(words[i])
		}
	}
}

func TestLoudsTrieBuilder_2(t *testing.T) {
	words := []string{"a", "aa", "ab", "bb"}
	keys := make([][]uint16, len(words))
	for i := 0; i < len(words); i++ {
		keys[i] = utf16.Encode([]rune(words[i]))
	}
	trie, nodes, err := migemo.BuildLoudsTrie(keys)
	if err != nil {
		t.Error()
	}
	var expectedNodes = []uint32{2, 4, 5, 6}
	for i := 0; i < len(nodes); i++ {
		if expectedNodes[i] != nodes[i] {
			t.Error(words[i])
		}
	}
	for i := 0; i < len(nodes); i++ {
		if uint32(trie.Lookup(keys[i])) != expectedNodes[i] {
			t.Error(words[i])
		}
	}
	if trie.Lookup(utf16.Encode([]rune(""))) != 1 {
		t.Error()
	}
	if trie.Lookup(utf16.Encode([]rune("bbb"))) != -1 {
		t.Error()
	}
	if trie.Lookup(utf16.Encode([]rune("c"))) != -1 {
		t.Error()
	}
}

func TestLoudsTrieBuilder_NotSorted(t *testing.T) {
	words := []string{"aa", "a"}
	keys := make([][]uint16, len(words))
	for i := 0; i < len(words); i++ {
		keys[i] = utf16.Encode([]rune(words[i]))
	}
	trie, nodes, err := migemo.BuildLoudsTrie(keys)
	if trie != nil || nodes != nil || err == nil {
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
	if trie.Lookup(utf16.Encode([]rune("box"))) != 10 {
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
