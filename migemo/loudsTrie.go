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

// PredictiveSearchDepthFirst は、指定したノードから葉の方向に全てのノードを深さ優先で巡る
func (this *LoudsTrie) PredictiveSearchDepthFirst(index int, f func(int, []uint16)) {
	key := make([]uint16, 0, 8)
	f(index, key)
	childPos := this.bitVector.Select(uint32(index), false) + 1
	if this.bitVector.Get(uint32(childPos)) {
		child := int(this.bitVector.Rank(childPos, true)) + 1
		for this.bitVector.Get(uint32(childPos)) {
			key = append(key, this.edges[child])
			this.predictiveSearchDepthFirstInternal(child, &key, f)
			key = (key)[:len(key)-1]
			child++
			childPos++
		}
	}
}

func (this *LoudsTrie) predictiveSearchDepthFirstInternal(index int, key *[]uint16, f func(int, []uint16)) {
	f(index, *key)
	childPos := this.bitVector.Select(uint32(index), false) + 1
	if this.bitVector.Get(uint32(childPos)) {
		child := int(this.bitVector.Rank(childPos, true)) + 1
		for this.bitVector.Get(uint32(childPos)) {
			*key = append(*key, this.edges[child])
			this.predictiveSearchDepthFirstInternal(child, key, f)
			*key = (*key)[:len(*key)-1]
			child++
			childPos++
		}
	}
}

// PredictiveSearchBreadthFirst は、指定したノードから葉の方向に全てのノードを幅優先で巡る．
func (this *LoudsTrie) PredictiveSearchBreadthFirst(node int, f func(int)) {
	lower := uint(node)
	upper := uint(node + 1)
	for upper-lower > 0 {
		for i := lower; i < upper; i++ {
			f(int(i))
		}
		lower = this.bitVector.Rank(this.bitVector.Select(uint32(lower), false)+1, true) + 1
		upper = this.bitVector.Rank(this.bitVector.Select(uint32(upper), false)+1, true) + 1
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
