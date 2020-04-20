package migemo

import (
	"math"
	"math/bits"
)

type BitVector struct {
	words      []uint64
	sizeInBits uint32
	lb         []uint32
	sb         []uint16
}

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

func (this *BitVector) Rank(pos uint, b bool) uint {
	/*if pos < 0 && this.sizeInBits <= pos {
		return nil;
	}*/
	var count1 uint = uint(this.sb[pos/64]) + uint(this.lb[pos/512])
	var word = this.words[pos/64]
	var shiftSize = 64 - (pos & 63)
	var mask uint64 = 0
	if shiftSize < 64 {
		mask = uint64(0xFFFFFFFFFFFFFFFF) >> shiftSize
	}
	count1 += uint(bits.OnesCount64(word & mask))
	if b {
		return count1
	} else {
		return pos - count1
	}
}

func (this *BitVector) Select(count uint32, b bool) uint {
	var lbIndex uint32 = this.lowerBoundBinarySearchLB(count, b) - 1
	var countInLb uint32
	var countInSb uint32
	if b {
		countInLb = count - this.lb[lbIndex]
	} else {
		countInLb = count - uint32(512*lbIndex-this.lb[lbIndex])
	}
	var sbIndex = this.lowerBoundBinarySearchSB(uint16(countInLb), lbIndex*8, lbIndex*8+8, b) - 1
	if b {
		countInSb = countInLb - uint32(this.sb[sbIndex])
	} else {
		countInSb = countInLb - (64*(sbIndex%8) - uint32(this.sb[sbIndex]))
	}
	var word = this.words[sbIndex]
	if !b {
		word = ^word
	}
	return uint(sbIndex*64) + selectInWord(word, uint(countInSb)) - 1
}

func selectInWord(word uint64, count uint) uint {
	var lower_bit_count = uint(bits.OnesCount32(uint32(word)))
	var i uint = 0
	if lower_bit_count < count {
		word = word >> 32
		count = count - lower_bit_count
		i = 32
	}
	var lower16bit_count = uint(bits.OnesCount16(uint16(word)))
	if lower16bit_count < count {
		word = word >> 16
		count = count - lower16bit_count
		i = i + 16
	}
	var lower8bit_count = uint(bits.OnesCount8(uint8(word)))
	if lower8bit_count < count {
		word = word >> 8
		count = count - lower8bit_count
		i = i + 8
	}
	var lower4bit_count = uint(bits.OnesCount8(uint8(word) & 0b1111))
	if lower4bit_count < count {
		word = word >> 4
		count = count - lower4bit_count
		i = i + 4
	}
	for count > 0 {
		count = count - uint(word&1)
		word = word >> 1
		i = i + 1
	}
	return i
}

func (this *BitVector) lowerBoundBinarySearchLB(key uint32, b bool) uint32 {
	var high = len(this.lb)
	var low = -1
	if b {
		for high-low > 1 {
			var mid = int(high+low) >> 1
			if this.lb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	} else {
		for high-low > 1 {
			var mid = int(high+low) >> 1
			if uint32(mid<<9)-this.lb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	}
	return uint32(high)
}

func (this *BitVector) lowerBoundBinarySearchSB(key uint16, fromIndex uint32, toIndex uint32, b bool) uint32 {
	var high = toIndex
	var low = fromIndex - 1
	if b {
		for high-low > 1 {
			mid := (high + low) >> 1
			if this.sb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	} else {
		for high-low > 1 {
			mid := (high + low) >> 1
			if uint16(mid&7)<<6-this.sb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		}
	}
	return high
}

func (this *BitVector) NextClearBit(fromIndex uint) uint {
	var u = int(fromIndex / 64)
	if u >= len(this.words) {
		return fromIndex
	}
	var word uint64 = ^this.words[u] & uint64(^uint64(0)<<(fromIndex&63))
	for true {
		if word != 0 {
			return uint(u*64) + uint(bits.TrailingZeros64(word))
		}
		u = u + 1
		if u == len(this.words) {
			return uint(len(this.words)) * 64
		}
		word = ^this.words[u]
	}
	return math.MaxUint32 // Unreachable here
}

func (this *BitVector) Size() int {
	return int(this.sizeInBits)
}

func (this *BitVector) Get(pos uint32) bool {
	/*
		if (pos < 0 && this.sizeInBits <= pos) {
			return nil;
		}
	*/
	return ((this.words[pos>>6] >> (pos & 63)) & 1) == 1
}
