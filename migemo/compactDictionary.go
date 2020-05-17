package migemo

import (
	"bufio"
	"encoding/binary"
	"os"
)

// CompactDictionary は、読み毎に複数の単語を格納した辞書
type CompactDictionary struct {
	keyTrie          *LoudsTrie
	valueTrie        *LoudsTrie
	mappingBitVector *BitVector
	mapping          []uint32
	// hasMappingBitList は、あるノードがマッピングを持つかを格納する
	hasMappingBitList *BitList
}

// NewCompactDictionary は、バイト配列からCompactDictionaryを読み込む
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
	hasMappingBitList := createHasMappingBitList(mappingBitVector)
	return &CompactDictionary{
		keyTrie,
		valueTrie,
		mappingBitVector,
		mapping,
		hasMappingBitList,
	}
}

func createHasMappingBitList(mappingBitVector *BitVector) *BitList {
	numOfNodes := mappingBitVector.Rank(uint(mappingBitVector.Size())+1, false)
	bitList := NewBitListWithSize(int(numOfNodes))
	bitPosition := uint(0)
	for node := 1; node < int(numOfNodes); node++ {
		var hasMapping = mappingBitVector.Get(uint32(bitPosition) + 1)
		bitList.Set(node, hasMapping)
		bitPosition = mappingBitVector.NextClearBit(bitPosition + 1)
	}
	return bitList
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

// Search は、キーに一致する単語をコールバック関数fに返す
func (compactDictionary *CompactDictionary) Search(key []uint16, f func([]uint16)) {
	var keyIndex = compactDictionary.keyTrie.Lookup(key)
	if keyIndex != -1 {
		var valueStartPos = compactDictionary.mappingBitVector.Select(uint32(keyIndex), false)
		var valueEndPos = compactDictionary.mappingBitVector.NextClearBit(valueStartPos + 1)
		var size = uint(valueEndPos - valueStartPos - 1)
		if size > 0 {
			var offset = compactDictionary.mappingBitVector.Rank(valueStartPos, false)
			word := make([]uint16, 0, 16)
			for i := uint(0); i < size; i++ {
				compactDictionary.valueTrie.ReverseLookup(compactDictionary.mapping[valueStartPos-offset+i], &word)
				f(word)
				word = word[:0]
			}
		}
	}
}

// PredictiveSearch は、接頭辞がkeyに一致する全ての単語をコールバック関数fに返す
func (compactDictionary *CompactDictionary) PredictiveSearch(key []uint16, f func([]uint16)) {
	var keyIndex = compactDictionary.keyTrie.Lookup(key)
	word := make([]uint16, 0, 16)
	if keyIndex > 1 {
		compactDictionary.keyTrie.PredictiveSearchBreadthFirst(keyIndex, func(i int) {
			if compactDictionary.hasMappingBitList.Get(i) {
				var valueStartPos uint = compactDictionary.mappingBitVector.Select(uint32(i), false)
				var valueEndPos uint = compactDictionary.mappingBitVector.NextClearBit(valueStartPos + 1)
				var size = valueEndPos - valueStartPos - 1
				var offset = compactDictionary.mappingBitVector.Rank(valueStartPos, false)
				for j := uint(0); j < size; j++ {
					compactDictionary.valueTrie.ReverseLookup(compactDictionary.mapping[valueStartPos-offset+j], &word)
					f(word)
					word = word[:0]
				}
			}
		})
	}
}

// Save は、ファイルに辞書の内容を書き込む
func (compactDictionary *CompactDictionary) Save(fp *os.File) {
	writer := bufio.NewWriter(fp)
	buffer := make([]byte, 8)
	// output key trie
	keyTriEdgesLength := uint32((len(compactDictionary.keyTrie.edges)))
	binary.BigEndian.PutUint32(buffer, keyTriEdgesLength)
	writer.Write(buffer[0:4])
	for i := 0; i < len(compactDictionary.keyTrie.edges); i++ {
		buffer[0] = encode(compactDictionary.keyTrie.edges[i])
		writer.Write(buffer[:1])
	}
	binary.BigEndian.PutUint32(buffer, uint32(compactDictionary.keyTrie.bitVector.Size()))
	writer.Write(buffer[0:4])
	keyTrieBitVectorWords := compactDictionary.keyTrie.bitVector.words
	for i := 0; i < len(keyTrieBitVectorWords); i++ {
		binary.BigEndian.PutUint64(buffer, keyTrieBitVectorWords[i])
		writer.Write(buffer)
	}
	// output value trie
	valueTriEdgesLength := uint32((len(compactDictionary.valueTrie.edges)))
	binary.BigEndian.PutUint32(buffer, valueTriEdgesLength)
	writer.Write(buffer[:4])
	for i := 0; i < len(compactDictionary.valueTrie.edges); i++ {
		binary.BigEndian.PutUint16(buffer, compactDictionary.valueTrie.edges[i])
		writer.Write(buffer[:2])
	}
	binary.BigEndian.PutUint32(buffer, uint32(compactDictionary.valueTrie.bitVector.Size()))
	writer.Write(buffer[0:4])
	valueTrieBitVectorWords := compactDictionary.valueTrie.bitVector.words
	for i := 0; i < len(valueTrieBitVectorWords); i++ {
		binary.BigEndian.PutUint64(buffer, valueTrieBitVectorWords[i])
		writer.Write(buffer)
	}

	// output mapping trie
	binary.BigEndian.PutUint32(buffer, uint32(compactDictionary.mappingBitVector.sizeInBits))
	writer.Write(buffer[:4])
	for _, w := range compactDictionary.mappingBitVector.words {
		binary.BigEndian.PutUint64(buffer, w)
		writer.Write(buffer)
	}
	binary.BigEndian.PutUint32(buffer, uint32(len(compactDictionary.mapping)))
	writer.Write(buffer[:4])
	for _, x := range compactDictionary.mapping {
		binary.BigEndian.PutUint32(buffer, x)
		writer.Write(buffer[:4])
	}
	writer.Flush()
}
