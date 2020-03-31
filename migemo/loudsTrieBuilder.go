package migemo

import "unicode/utf16"

func BuildLoudsTrie(keys [][]uint16, generatedIndexes []uint32) *LoudsTrie {
	var memo []uint32
	if generatedIndexes == nil {
		memo = make([]uint32, len(keys))
	} else if len(generatedIndexes) == len(keys) {
		memo = generatedIndexes
	} else {
		return nil
	}

	for i := 0; i < len(keys); i++ {
		if i > 0 && string(utf16.Decode(keys[i-1])) > string(utf16.Decode(keys[i])) {
			return nil
		}
	}
	for i := 0; i < len(memo); i++ {
		memo[i] = 1
	}
	var offset = 0
	var currentNode uint32 = 1
	var edges = []uint16{0x30, 0x30} // TODO: '0'で穴埋めを'\0'にするか、なくす
	var childSizes = make([]uint32, 128)
	for true {
		var lastChar uint16 = 0
		var lastParent uint32 = 0
		var restKeys uint32 = 0
		for i := 0; i < len(keys); i++ {
			if memo[i] < 0 {
				continue
			}
			if len(keys[i]) <= offset {
				memo[i] = -memo[i]
				continue
			}
			var currentChar = keys[i][offset]
			var currentParent = memo[i]
			if lastChar != currentChar || lastParent != currentParent {
				if uint32(len(childSizes)) <= memo[i] {
					var tmp = make([]uint32, len(childSizes)*2)
					copy(tmp, childSizes)
					childSizes = tmp
				}
				childSizes[memo[i]]++
				currentNode = currentNode + 1
				edges = append(edges, currentChar)
				lastChar = currentChar
				lastParent = currentParent
			}
			memo[i] = currentNode
			restKeys++
		}
		if restKeys == 0 {
			break
		}
		offset++
	}

	for i := 0; i < len(memo); i++ {
		memo[i] = -memo[i]
	}

	var numOfChildren uint32 = 0
	for i := uint32(1); i <= currentNode; i++ {
		numOfChildren = numOfChildren + childSizes[i]
	}
	var numOfNodes = currentNode
	var bitVectorWords = make([]uint32, ((numOfChildren+uint32(numOfNodes)+63+1)>>6)*2)
	var bitVectorIndex uint32 = 1
	bitVectorWords[0] = 1
	for i := uint32(1); i <= currentNode; i++ {
		bitVectorIndex++
		var childSize = childSizes[i]
		for j := uint32(0); j < childSize; j++ {
			bitVectorWords[bitVectorIndex>>5] |= 1 << (bitVectorIndex & 31)
			bitVectorIndex++
		}
	}

	var bitVector = NewBitVector(bitVectorWords, bitVectorIndex)
	return NewLoudsTrie(bitVector, edges)
}
