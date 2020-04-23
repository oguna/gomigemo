package migemo

import (
	"encoding/binary"
	"os"
)

type CompactDictionary struct {
	keyTrie          *LoudsTrie
	valueTrie        *LoudsTrie
	mappingBitVector *BitVector
	mapping          []uint32
}

func (this *CompactDictionary) GetMappingBitVector() *BitVector {
	return this.mappingBitVector
}

func NewCompactDictionary(buffer []uint8) *CompactDictionary {
	var offset = 0
	keyTrie, offset := readTrie(buffer, offset, true)
	valueTrie, offset := readTrie(buffer, offset, false)
	var mappingBitVectorSize = binary.BigEndian.Uint32(buffer[offset:])
	offset = offset + 4
	var mappingBitVectorWords = make([]uint64, (mappingBitVectorSize+63)/64)
	for i := 0; i < len(mappingBitVectorWords); i++ {
		mappingBitVectorWords[i] = binary.BigEndian.Uint64(buffer[offset:])
		offset = offset + 8
	}
	mappingBitVector := NewBitVector(mappingBitVectorWords, mappingBitVectorSize)
	var mappingSize = binary.BigEndian.Uint32(buffer[offset:])
	offset = offset + 4
	mapping := make([]uint32, mappingSize)
	for i := uint32(0); i < mappingSize; i++ {
		mapping[i] = binary.BigEndian.Uint32(buffer[offset:])
		offset += 4
	}
	if offset != len(buffer) {
		return nil
	}
	return &CompactDictionary{
		keyTrie,
		valueTrie,
		mappingBitVector,
		mapping,
	}
}

func readTrie(buffer []uint8, offset int, compactHiragana bool) (*LoudsTrie, int) {
	var keyTrieEdgeSize = binary.BigEndian.Uint32(buffer[offset:])
	offset = offset + 4
	var keyTrieEdges = make([]uint16, keyTrieEdgeSize)
	for i := uint32(0); i < keyTrieEdgeSize; i++ {
		var c uint16
		if compactHiragana {
			c = decode(buffer[offset])
			offset = offset + 1
		} else {
			c = binary.BigEndian.Uint16(buffer[offset:])
			offset = offset + 2
		}
		keyTrieEdges[i] = c
	}
	var keyTrieBitVectorSize = binary.BigEndian.Uint32(buffer[offset:])
	offset = offset + 4
	var keyTrieBitVectorWords = make([]uint64, (keyTrieBitVectorSize+63)/64)
	for i := 0; i < len(keyTrieBitVectorWords); i++ {
		keyTrieBitVectorWords[i] = binary.BigEndian.Uint64(buffer[offset:])
		offset = offset + 8
	}
	return NewLoudsTrie(NewBitVector(keyTrieBitVectorWords, keyTrieBitVectorSize), keyTrieEdges), offset
}

func decode(c uint8) uint16 {
	if 0x20 <= c && c <= 0x7e {
		return uint16(c)
	}
	if 0xa1 <= c && c <= 0xf6 {
		return uint16(c) + 0x3040 - 0xa0
	}
	return 0
}

func encode(c uint16) uint8 {
	if 0x20 <= c && c <= 0x7e {
		return uint8(c)
	}
	if 0x3041 <= c && c <= 0x3096 {
		return uint8(c - 0x3040 + 0xa0)
	}
	if 0x30fc == c {
		return uint8(c - 0x3040 + 0xa0)
	}
	return 0
}

func (this *CompactDictionary) Search(key []uint16, f func([]uint16)) {
	var keyIndex = this.keyTrie.Lookup(key)
	if keyIndex != -1 {
		var valueStartPos = this.mappingBitVector.Select(uint32(keyIndex), false)
		var valueEndPos = this.mappingBitVector.NextClearBit(valueStartPos + 1)
		var size = uint(valueEndPos - valueStartPos - 1)
		if size > 0 {
			var offset = this.mappingBitVector.Rank(valueStartPos, false)
			word := make([]uint16, 0, 16)
			for i := uint(0); i < size; i++ {
				this.valueTrie.ReverseLookup(this.mapping[valueStartPos-offset+i], &word)
				f(word)
				word = word[:0]
			}
		}
	}
}

func (this *CompactDictionary) PredictiveSearch(key []uint16, f func([]uint16)) {
	var keyIndex = this.keyTrie.Lookup(key)
	if keyIndex > 1 {
		this.keyTrie.PredictiveSearch(keyIndex, func(i int) {
			var valueStartPos uint = this.mappingBitVector.Select(uint32(i), false)
			var valueEndPos uint = this.mappingBitVector.NextClearBit(valueStartPos + 1)
			var size = valueEndPos - valueStartPos - 1
			var offset = this.mappingBitVector.Rank(valueStartPos, false)
			word := make([]uint16, 0, 16)
			for j := uint(0); j < size; j++ {
				this.valueTrie.ReverseLookup(this.mapping[valueStartPos-offset+j], &word)
				f(word)
				word = word[:0]
			}
		})
	}
}

func (this *CompactDictionary) Save(fp *os.File) {
	buffer := make([]byte, 8)
	// output key trie
	keyTriEdgesLength := uint32((len(this.keyTrie.edges)))
	binary.BigEndian.PutUint32(buffer, keyTriEdgesLength)
	fp.Write(buffer[0:4])
	for i := 0; i < len(this.keyTrie.edges); i++ {
		buffer[0] = encode(this.keyTrie.edges[i])
		fp.Write(buffer[:1])
	}
	binary.BigEndian.PutUint32(buffer, uint32(this.keyTrie.bitVector.Size()))
	fp.Write(buffer[0:4])
	keyTrieBitVectorWords := this.keyTrie.bitVector.words
	for i := 0; i < len(keyTrieBitVectorWords); i++ {
		binary.BigEndian.PutUint64(buffer, keyTrieBitVectorWords[i])
		fp.Write(buffer)
	}
	// output value trie
	valueTriEdgesLength := uint32((len(this.valueTrie.edges)))
	binary.BigEndian.PutUint32(buffer, valueTriEdgesLength)
	fp.Write(buffer[:4])
	for i := 0; i < len(this.valueTrie.edges); i++ {
		binary.BigEndian.PutUint16(buffer, this.valueTrie.edges[i])
		fp.Write(buffer[:2])
	}
	binary.BigEndian.PutUint32(buffer, uint32(this.valueTrie.bitVector.Size()))
	fp.Write(buffer[0:4])
	valueTrieBitVectorWords := this.valueTrie.bitVector.words
	for i := 0; i < len(valueTrieBitVectorWords); i++ {
		binary.BigEndian.PutUint64(buffer, valueTrieBitVectorWords[i])
		fp.Write(buffer)
	}

	// output mapping trie
	binary.BigEndian.PutUint32(buffer, uint32(this.mappingBitVector.sizeInBits))
	fp.Write(buffer[:4])
	for _, w := range this.mappingBitVector.words {
		binary.BigEndian.PutUint64(buffer, w)
		fp.Write(buffer)
	}
	binary.BigEndian.PutUint32(buffer, uint32(len(this.mapping)))
	fp.Write(buffer[:4])
	for _, x := range this.mapping {
		binary.BigEndian.PutUint32(buffer, x)
		fp.Write(buffer[:4])
	}
}
