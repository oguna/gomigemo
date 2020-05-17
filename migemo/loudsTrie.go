package migemo

// LoudsTrie は、LOUDS(level order unary degree sequence)を実装したもの
type LoudsTrie struct {
	bitVector *BitVector
	edges     []uint16
}

// NewLoudsTrie は、LoudsTrieを初期化する
func NewLoudsTrie(bitVector *BitVector, edges []uint16) *LoudsTrie {
	return &LoudsTrie{
		bitVector,
		edges,
	}
}

// ReverseLookup は、ノード番号indexからkeyを復元する
func (trie *LoudsTrie) ReverseLookup(index uint32, key *[]uint16) int {
	offset := len(*key)
	for index > 1 {
		*key = append(*key, trie.edges[index])
		index = trie.Parent(index)
	}
	for i, j := offset, len(*key)-1; i < j; i, j = i+1, j-1 {
		(*key)[i], (*key)[j] = (*key)[j], (*key)[i]
	}
	return len(*key) - offset
}

// Parent は、ノード番号indexの親を返す
func (trie *LoudsTrie) Parent(x uint32) uint32 {
	return uint32(trie.bitVector.Rank(trie.bitVector.Select(x, true), false))
}

// FirstChild は、ノード番号xのはじめの子供のノード番号を返す。子供がなければ-1
func (trie *LoudsTrie) FirstChild(x uint32) int {
	y := trie.bitVector.Select(x, false) + 1
	if trie.bitVector.Get(uint32(y)) {
		return int(trie.bitVector.Rank(y, true)) + 1
	}
	return -1
}

// Traverse は、ノード番号indexの子ノードのうち、ラベルcを持つノード番号を返す。見つからなければ-1
func (trie *LoudsTrie) Traverse(index uint32, c uint16) int {
	firstChild := trie.FirstChild(index)
	if firstChild == -1 {
		return -1
	}
	var childStartBit = trie.bitVector.Select(uint32(firstChild), true)
	var childEndBit = trie.bitVector.NextClearBit(childStartBit)
	var childSize = childEndBit - childStartBit
	var result = binarySearchUint16(trie.edges, uint32(firstChild), uint32(firstChild)+uint32(childSize), c)
	if result >= 0 {
		return result
	}
	return -1
}

// Lookup は、検索対象keyのノード番号を返す。見つからければ-1
func (trie *LoudsTrie) Lookup(key []uint16) int {
	var nodeIndex int = 1
	for _, c := range key {
		nodeIndex = trie.Traverse(uint32(nodeIndex), c)
		if nodeIndex == -1 {
			break
		}
	}
	if nodeIndex >= 0 {
		return nodeIndex
	}
	return -1
}

// PredictiveSearchDepthFirst は、指定したノードから葉の方向に全てのノードを深さ優先で巡る
func (trie *LoudsTrie) PredictiveSearchDepthFirst(index int, f func(int, []uint16)) {
	key := make([]uint16, 0, 8)
	f(index, key)
	childPos := trie.bitVector.Select(uint32(index), false) + 1
	if trie.bitVector.Get(uint32(childPos)) {
		child := int(trie.bitVector.Rank(childPos, true)) + 1
		for trie.bitVector.Get(uint32(childPos)) {
			key = append(key, trie.edges[child])
			trie.predictiveSearchDepthFirstInternal(child, &key, f)
			key = (key)[:len(key)-1]
			child++
			childPos++
		}
	}
}

func (trie *LoudsTrie) predictiveSearchDepthFirstInternal(index int, key *[]uint16, f func(int, []uint16)) {
	f(index, *key)
	childPos := trie.bitVector.Select(uint32(index), false) + 1
	if trie.bitVector.Get(uint32(childPos)) {
		child := int(trie.bitVector.Rank(childPos, true)) + 1
		for trie.bitVector.Get(uint32(childPos)) {
			*key = append(*key, trie.edges[child])
			trie.predictiveSearchDepthFirstInternal(child, key, f)
			*key = (*key)[:len(*key)-1]
			child++
			childPos++
		}
	}
}

// PredictiveSearchBreadthFirst は、指定したノードから葉の方向に全てのノードを幅優先で巡る．
func (trie *LoudsTrie) PredictiveSearchBreadthFirst(node int, f func(int)) {
	lower := uint(node)
	upper := uint(node + 1)
	for upper-lower > 0 {
		for i := lower; i < upper; i++ {
			f(int(i))
		}
		lower = trie.bitVector.Rank(trie.bitVector.Select(uint32(lower), false)+1, true) + 1
		upper = trie.bitVector.Rank(trie.bitVector.Select(uint32(upper), false)+1, true) + 1
	}
}

// Size は、ノードの個数を返す
func (trie *LoudsTrie) Size() int {
	return len(trie.edges) - 2
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
