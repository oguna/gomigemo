package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestDoubleArray_Lookup(t *testing.T) {
	base := []int16{0, 3, -1, -1, 2, -1, -1, -1, -1}
	check := []int16{-1, 0, 0, 4, 0, 1, 1, -1, -1}
	charConverter := func(c uint8) int {
		if 'a' <= c && c <= 'z' {
			return int(c-'a') + 1
		} else {
			return -1
		}
	}
	charSize := 26
	trie := migemo.NewDoubleArray(base, check, charConverter, charSize)
	if trie.Lookup("") != 0 {
		t.Error()
	}
	if trie.Lookup("a") != 1 {
		t.Error()
	}
	if trie.Lookup("ab") != 5 {
		t.Error()
	}
	if trie.Lookup("ac") != 6 {
		t.Error()
	}
	if trie.Lookup("b") != 2 {
		t.Error()
	}
	if trie.Lookup("d") != 4 {
		t.Error()
	}
	if trie.Lookup("da") != 3 {
		t.Error()
	}
	if trie.Lookup("c") != -1 {
		t.Error()
	}
}

func TestDoubleArray_PredictiveSearch(t *testing.T) {
	base := []int16{0, 3, -1, -1, 2, -1, -1, -1, -1}
	check := []int16{-1, 0, 0, 4, 0, 1, 1, -1, -1}
	charConverter := func(c uint8) int {
		if 'a' <= c && c <= 'z' {
			return int(c-'a') + 1
		} else {
			return -1
		}
	}
	charSize := 26
	trie := migemo.NewDoubleArray(base, check, charConverter, charSize)
	list := make([]int16, 0)
	trie.PredictiveSearch("ab", func(node int16) { list = append(list, node) })
	if len(list) == 1 && list[0] != 5 {
		t.Error()
	}
	list = list[:0]
	trie.PredictiveSearch("a", func(node int16) { list = append(list, node) })
	if len(list) == 3 && list[0] != 1 && list[1] != 5 && list[2] == 6 {
		t.Error()
	}
	list = list[:0]
}

func TestDoubleArray_CommonPrefixSearch(t *testing.T) {
	base := []int16{0, 3, -1, -1, 2, -1, -1, -1, -1}
	check := []int16{-1, 0, 0, 4, 0, 1, 1, -1, -1}
	charConverter := func(c uint8) int {
		if 'a' <= c && c <= 'z' {
			return int(c-'a') + 1
		} else {
			return -1
		}
	}
	charSize := 26
	trie := migemo.NewDoubleArray(base, check, charConverter, charSize)
	list := make([]int16, 0)
	trie.CommonPrefixSearch("ab", func(node int16) { list = append(list, node) })
	if len(list) == 3 && list[0] != 0 && list[1] == 1 && list[2] == 5 {
		t.Error()
	}
	list = list[:0]
	trie.CommonPrefixSearch("ac", func(node int16) { list = append(list, node) })
	if len(list) == 3 && list[0] != 0 && list[1] == 1 && list[2] == 6 {
		t.Error()
	}
	list = list[:0]
	trie.CommonPrefixSearch("b", func(node int16) { list = append(list, node) })
	if len(list) == 2 && list[0] != 0 && list[1] == 2 {
		t.Error()
	}
	list = list[:0]
	trie.CommonPrefixSearch("", func(node int16) { list = append(list, node) })
	if len(list) == 1 && list[0] != 0 {
		t.Error()
	}
	list = list[:0]
	trie.CommonPrefixSearch("c", func(node int16) { list = append(list, node) })
	if len(list) == 1 && list[0] != 0 {
		t.Error()
	}
	list = list[:0]
}
