package migemo

import "math/bits"

type BitVector struct {
	words      []uint32
	sizeInBits uint32
	lb         []uint32
	sb         []uint16
}

func NewBitVector(words []uint32, sizeInBits uint32) *BitVector {
	if (sizeInBits+63)>>5 != uint32(len(words)) {
		return nil
	}
	lb := make([]uint32, (sizeInBits+511)>>9)
	sb := make([]uint16, len(lb)*8)
	var sum = 0
	var sumInLb = 0
	for i := 0; i < len(sb); i++ {
		var bc = 0
		if i < len(words)>>1 {
			bc = bits.OnesCount32(words[i*2]) + bits.OnesCount32(words[i*2+1])
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

func (this *BitVector) Rank(pos uint32, b bool) uint32 {
	/*if pos < 0 && this.sizeInBits <= pos {
		return nil;
	}*/
	var count1 uint32 = uint32(this.sb[pos>>6]) + this.lb[pos>>9]
	var posInDWord uint32 = pos & uint32(63)
	if posInDWord >= 32 {
		count1 += uint32(bits.OnesCount32(this.words[(pos>>5)&uint32(0xFFFFFFFE)]))
	}
	posInWord := pos & 31
	var mask uint32 = 0x7FFFFFFF >> (31 - posInWord)
	count1 += uint32(bits.OnesCount32(this.words[pos>>5] & mask))
	if b {
		return count1
	} else {
		return pos - count1
	}
}

func (this *BitVector) Select(count uint32, b bool) uint32 {
	var lbIndex uint32 = this.lowerBoundBinarySearchLB(count, b) - 1
	var countInLb uint32
	var countInSb uint32
	if b {
		countInLb = count - this.lb[lbIndex]
	} else {
		countInLb = count - uint32(512*lbIndex-this.lb[lbIndex])
	}
	var sbIndex = this.lowerBoundBinarySearchSB(countInLb, lbIndex*8, lbIndex*8+8, b) - 1
	if b {
		countInSb = countInLb - uint32(this.sb[sbIndex])
	} else {
		countInSb = countInLb - (64*(sbIndex%8) - uint32(this.sb[sbIndex]))
	}
	var wordL = this.words[sbIndex*2]
	var wordU = this.words[sbIndex*2+1]
	if !b {
		wordL = ^wordL
		wordU = ^wordU
	}
	var lowerBitCount = uint32(bits.OnesCount32(wordL))
	var i uint32 = 0
	if countInSb > lowerBitCount {
		wordL = wordU
		countInSb -= lowerBitCount
		i = 32
	}
	for countInSb > 0 {
		countInSb = countInSb - (wordL & 1)
		wordL = wordL >> 1
		i++
	}
	return sbIndex*64 + (i - 1)
}

func (this *BitVector) lowerBoundBinarySearchLB(key uint32, b bool) uint32 {
	var high = len(this.lb)
	var low = -1
	for high-low > 1 {
		var mid = int(high+low) >> 1
		if b {
			if this.lb[mid] < key {
				low = mid
			} else {
				high = mid
			}
		} else {
			if 512*mid-int(this.lb[mid]) < int(key) {
				low = mid
			} else {
				high = mid
			}
		}
	}
	return uint32(high)
}

func (this *BitVector) lowerBoundBinarySearchSB(key uint32, fromIndex uint32, toIndex uint32, b bool) uint32 {
	var high = toIndex
	var low = fromIndex - 1
	for high-low > 1 {
		var mid = (high + low) >> 1
		if b && uint32(this.sb[mid]) < key {
			low = mid
		} else if !b && (uint32(64*(mid&7)-uint32(this.sb[mid])) < uint32(key)) {
			low = mid
		} else {
			high = mid
		}
	}
	return high
}

func (this *BitVector) NextClearBit(fromIndex uint32) uint32 {
	for fromIndex < this.sizeInBits && this.Get(fromIndex) {
		fromIndex = fromIndex + 1
	}
	if fromIndex > this.sizeInBits {
		return this.sizeInBits + 1
	}
	return fromIndex
	/*
		var u = int(fromIndex >> 5);
		var word uint32 = ^this.words[u] & (0xffffffff  << fromIndex); // TODO
		for true {
			if word != 0 {
				return (u * 32 + bits.TrailingZeros32(word));
			}
			u = u + 1
			if u == len(this.words) {
				return -1;
			}
			word = ^this.words[u];
		}
		return -1;
	*/
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
	return ((this.words[pos>>5] >> (pos & 31)) & 1) == 1
}
