package migemo_test

import (
	"testing"

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
	list := make([]int, 0)
	trie.PredictiveSearchDepthFirst(4, func(node int) {
		list = append(list, node)
	})
	expectedList := []int{4, 7, 13, 8, 9, 14}
	if len(expectedList) != len(list) {
		t.Errorf("invalid size. expected:%d actual:%d", len(expectedList), len(list))
	}
	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != list[i] {
			t.Errorf("expected:%d actual:%d", expectedList[i], list[i])
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
