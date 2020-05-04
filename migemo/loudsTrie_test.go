package migemo_test

import (
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestLoudsTrie_PredictiveSearchDepthFirst(t *testing.T) {
	// bady,bad,bank,box,dad,dance
	bitVectorWords := []uint64{0b01000100010010110101110101101101}
	bitVectorSize := uint32(32)
	bitVector := migemo.NewBitVector(bitVectorWords, bitVectorSize)
	edgeString := "  bdaoabdnxdnykce"
	edges := make([]uint16, 17)
	for i := 0; i < 17; i++ {
		edges[i] = uint16(edgeString[i])
	}
	trie := migemo.NewLoudsTrie(bitVector, edges)
	nodes := make([]int, 0)
	keys := make([]string, 0)
	trie.PredictiveSearchDepthFirst(4, func(node int, key []uint16) {
		nodes = append(nodes, node)
		keys = append(keys, string(utf16.Decode(key)))
	})
	expectedNodes := []int{4, 7, 13, 8, 9, 14}
	expectedKeys := []string{"", "b", "by", "d", "n", "nk"}
	if len(expectedNodes) != len(nodes) {
		t.Errorf("invalid size. expected:%d actual:%d", len(expectedNodes), len(nodes))
	}
	for i := 0; i < len(expectedNodes); i++ {
		if expectedNodes[i] != nodes[i] {
			t.Errorf("expected:%d actual:%d", expectedNodes[i], nodes[i])
		}
		if expectedKeys[i] != keys[i] {
			t.Errorf("expected:%s actual:%s", expectedKeys[i], keys[i])
		}
	}
}

func TestLoudsTrie_PredictiveSearchBreadthFirst(t *testing.T) {
	// bady,bad,bank,box,dad,dance
	bitVectorWords := []uint64{0b01000100010010110101110101101101}
	bitVectorSize := uint32(32)
	bitVector := migemo.NewBitVector(bitVectorWords, bitVectorSize)
	edgeString := "  bdaoabdnxdnykce"
	edges := make([]uint16, 17)
	for i := 0; i < 17; i++ {
		edges[i] = uint16(edgeString[i])
	}
	trie := migemo.NewLoudsTrie(bitVector, edges)
	list := make([]int, 0)
	trie.PredictiveSearchBreadthFirst(4, func(node int) {
		list = append(list, node)
	})
	expectedList := []int{4, 7, 8, 9, 13, 14}
	if len(expectedList) != len(list) {
		t.Errorf("invalid size. expected:%d actual:%d", len(expectedList), len(list))
	}
	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != list[i] {
			t.Errorf("expected:%d actual:%d", expectedList[i], list[i])
		}
	}
}
