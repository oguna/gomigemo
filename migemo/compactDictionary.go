package migemo

import "encoding/binary"

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

func (this *CompactDictionary) Search(key []uint16, f func([]uint16)) {
	var keyIndex = this.keyTrie.Get(key)
	if keyIndex != -1 {
		var valueStartPos = this.mappingBitVector.Select(uint32(keyIndex), false)
		var valueEndPos = this.mappingBitVector.NextClearBit(valueStartPos + 1)
		var size = uint(valueEndPos - valueStartPos - 1)
		if size > 0 {
			var offset = this.mappingBitVector.Rank(valueStartPos, false)
			for i := uint(0); i < size; i++ {
				f(this.valueTrie.GetKey(this.mapping[valueStartPos-uint(offset)+i]))
			}
		}
	}
}

func (this *CompactDictionary) PredictiveSearch(key []uint16, f func([]uint16)) {
	var keyIndex = this.keyTrie.Get(key)
	if keyIndex > 1 {
		this.keyTrie.Iterate(keyIndex, func(i int) {
			var valueStartPos uint = this.mappingBitVector.Select(uint32(i), false)
			var valueEndPos uint = this.mappingBitVector.NextClearBit(valueStartPos + 1)
			var size = valueEndPos - valueStartPos - 1
			var offset = this.mappingBitVector.Rank(valueStartPos, false)
			for j := uint(0); j < size; j++ {
				f(this.valueTrie.GetKey(this.mapping[valueStartPos-offset+j]))
			}
		})
	}
}
