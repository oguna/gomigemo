package migemo

import (
	"math"
	"math/bits"
)

// BitVector は、RankやSelectの計算が高速なビット配列
type BitVector struct {
	words      []uint64
	sizeInBits uint32
	lb         []uint32
	sb         []uint16
}

// NewBitVector は、BitVectorを初期化する
func NewBitVector(words []uint64, sizeInBits uint32) *BitVector {
	if (sizeInBits+63)/64 != uint32(len(words)) {
		return nil
	}
	lb := make([]uint32, (sizeInBits+511)/512)
	sb := make([]uint16, len(lb)*8)
	var sum = 0
	var sumInLb = 0
	for i := 0; i < len(sb); i++ {
		var bc = 0
		if i < len(words) {
			bc = bits.OnesCount64(words[i])
		}
		sb[i] = uint16(sumInLb)
		sumInLb += bc
		if (i & 7) == 7 {
			lb[i>>3] = uint32(sum)
			sum += sumInLb
			sumInLb = 0
		}
	}
	return &BitVector{
		words,
		sizeInBits,
		lb,
		sb,
	}
}

// Rank は、pos位置のb値が何個あるかを返す
func (bitVector *BitVector) Rank(pos uint, b bool) uint {
	/*if pos < 0 && this.sizeInBits <= pos {
		return nil;
	}*/
	var count1 uint = uint(bitVector.sb[pos/64]) + uint(bitVector.lb[pos/512])
	var word = bitVector.words[pos/64]
	var mask = uint64(0xFFFFFFFFFFFFFFFF) >> (64 - pos&63)
	count1 += uint(bits.OnesCount64(word & mask))
	if b {
		return count1
	}
	return pos - count1
}

// Select は、b値のcount番目の位置を返す
func (bitVector *BitVector) Select(count uint32, b bool) uint {
	var lbIndex uint32 = bitVector.lowerBoundBinarySearchLB(count, b) - 1
	var countInLb uint32
	var countInSb uint32
	if b {
		countInLb = count - bitVector.lb[lbIndex]
	} else {
		countInLb = count - uint32(512*lbIndex-bitVector.lb[lbIndex])
	}
	var sbIndex = bitVector.lowerBoundBinarySearchSB(uint16(countInLb), lbIndex*8, lbIndex*8+8, b) - 1
	if b {
		countInSb = countInLb - uint32(bitVector.sb[sbIndex])
	} else {
		countInSb = countInLb - (64*(sbIndex%8) - uint32(bitVector.sb[sbIndex]))
	}
	var word = bitVector.words[sbIndex]
	if !b {
		word = ^word
	}
	return uint(sbIndex*64) + selectInWord(word, uint(countInSb)) - 1
}

func selectInWord(word uint64, count uint) uint {
	var lowerBitCount = uint(bits.OnesCount32(uint32(word)))
	var i uint = 0
	if lowerBitCount < count {
		word = word >> 32
		count = count - lowerBitCount
		i = 32
	}
	var lower16bitCount = uint(bits.OnesCount16(uint16(word)))
	if lower16bitCount < count {
		word = word >> 16
		count = count - lower16bitCount
		i = i + 16
	}
	var lower8bitCount = uint(bits.OnesCount8(uint8(word)))
	if lower8bitCount < count {
		word = word >> 8
		count = count - lower8bitCount
		i = i + 8
	}
	var lower4bitCount = uint(bits.OnesCount8(uint8(word) & 0b1111))
	if lower4bitCount < count {
		word = word >> 4
		count = count - lower4bitCount
		i = i + 4
	}
	for count > 0 {
		count = count - uint(word&1)
		word = word >> 1
		i = i + 1
	}
	return i
}

func (bitVector *BitVector) lowerBoundBinarySearchLB(key uint32, b bool) uint32 {
	var high = len(bitVector.lb)
	var low = -1
	if b {
		for high-low > 1 {
			var mid = int(high+low) >> 1
			if bitVector.lb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	} else {
		for high-low > 1 {
			var mid = int(high+low) >> 1
			if uint32(mid<<9)-bitVector.lb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	}
	return uint32(high)
}

func (bitVector *BitVector) lowerBoundBinarySearchSB(key uint16, fromIndex uint32, toIndex uint32, b bool) uint32 {
	var high = toIndex
	var low = fromIndex - 1
	if b {
		for high-low > 1 {
			mid := (high + low) >> 1
			if bitVector.sb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	} else {
		for high-low > 1 {
			mid := (high + low) >> 1
			if uint16(mid&7)<<6-bitVector.sb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	}
	return high
}

// NextClearBit は、fromIndexから次の0ビットの位置を取得する
func (bitVector *BitVector) NextClearBit(fromIndex uint) uint {
	var u = int(fromIndex / 64)
	if u >= len(bitVector.words) {
		return fromIndex
	}
	var word uint64 = ^bitVector.words[u] & uint64(^uint64(0)<<(fromIndex&63))
	for true {
		if word != 0 {
			return uint(u*64) + uint(bits.TrailingZeros64(word))
		}
		u = u + 1
		if u == len(bitVector.words) {
			return uint(len(bitVector.words)) * 64
		}
		word = ^bitVector.words[u]
	}
	return math.MaxUint32 // Unreachable here
}

// Size は、ビット配列の長さを返す
func (bitVector *BitVector) Size() int {
	return int(bitVector.sizeInBits)
}

// Get は、pos位置のビット値を返す
func (bitVector *BitVector) Get(pos uint32) bool {
	/*
		if (pos < 0 && this.sizeInBits <= pos) {
			return nil;
		}
	*/
	return ((bitVector.words[pos>>6] >> (pos & 63)) & 1) == 1
}
