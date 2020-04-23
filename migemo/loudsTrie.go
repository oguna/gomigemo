package migemo

type LoudsTrie struct {
	bitVector *BitVector
	edges     []uint16
}

func NewLoudsTrie(bitVector *BitVector, edges []uint16) *LoudsTrie {
	return &LoudsTrie{
		bitVector,
		edges,
	}
}

func (this *LoudsTrie) ReverseLookup(index uint32, key *[]uint16) int {
	offset := len(*key)
	for index > 1 {
		*key = append(*key, this.edges[index])
		index = this.Parent(index)
	}
	for i, j := offset, len(*key)-1; i < j; i, j = i+1, j-1 {
		(*key)[i], (*key)[j] = (*key)[j], (*key)[i]
	}
	return len(*key) - offset
}

func (this *LoudsTrie) Parent(x uint32) uint32 {
	return uint32(this.bitVector.Rank(this.bitVector.Select(x, true), false))
}

func (this *LoudsTrie) FirstChild(x uint32) int {
	y := this.bitVector.Select(x, false) + 1
	if this.bitVector.Get(uint32(y)) {
		return int(this.bitVector.Rank(y, true)) + 1
	} else {
		return -1
	}
}

func (this *LoudsTrie) Traverse(index uint32, c uint16) int {
	firstChild := this.FirstChild(index)
	if firstChild == -1 {
		return -1
	}
	var childStartBit = this.bitVector.Select(uint32(firstChild), true)
	var childEndBit = this.bitVector.NextClearBit(childStartBit)
	var childSize = childEndBit - childStartBit
	var result = binarySearchUint16(this.edges, uint32(firstChild), uint32(firstChild)+uint32(childSize), c)
	if result >= 0 {
		return result
	} else {
		return -1
	}
}

func (this *LoudsTrie) Lookup(key []uint16) int {
	var nodeIndex int = 1
	for _, c := range key {
		nodeIndex = this.Traverse(uint32(nodeIndex), c)
		if nodeIndex == -1 {
			break
		}
	}
	if nodeIndex >= 0 {
		return nodeIndex
	} else {
		return -1
	}
}

func (this *LoudsTrie) PredictiveSearch(index int, f func(int)) {
	f(index)
	var child = this.FirstChild(uint32(index))
	if child == -1 {
		return
	}
	var childPos = this.bitVector.Select(uint32(child), true)
	for this.bitVector.Get(uint32(childPos)) {
		this.PredictiveSearch(child, f)
		child++
		childPos++
	}
}

func (this *LoudsTrie) Size() int {
	return len(this.edges) - 2
}

func binarySearchUint16(a []uint16, fromIndex uint32, toIndex uint32, key uint16) int {
	var low = fromIndex
	var high = toIndex - 1
	for low <= high {
		var mid = (low + high) >> 1
		var midVal = a[mid]
		if midVal < key {
			low = mid + 1
		} else if midVal > key {
			high = mid - 1
		} else {
			return int(mid)
		}
	}
	return -int(low + 1)
}
