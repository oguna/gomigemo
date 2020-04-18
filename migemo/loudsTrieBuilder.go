package migemo

import (
	"errors"
)

type Level struct {
	louds  []bool
	outs   []bool
	labels []uint16
}

func (self *Level) Size() int {
	return len(self.louds) + len(self.outs) + len(self.labels)
}

type LoudsTrieBuilder struct {
	levels  []Level
	lastKey []uint16
}

func NewLoudsTrieBuilder() *LoudsTrieBuilder {
	level0 := Level{
		louds:  []bool{true, false},
		outs:   []bool{false},
		labels: []uint16{' ', ' '},
	}
	level1 := Level{
		louds: []bool{false},
	}
	levels := []Level{level0, level1}
	return &LoudsTrieBuilder{
		levels:  levels,
		lastKey: []uint16{},
	}
}

func (self *LoudsTrieBuilder) Add(key []uint16) error {
	if CompareUtf16String(key, self.lastKey) <= 0 {
		return errors.New("key must be larger than last added key")
	}
	if len(key) == 0 {
		self.levels[0].outs[0] = true
		return nil
	}
	if len(key)+1 >= len(self.levels) {
		self.levels = append(self.levels, make([]Level, len(key)+2-len(self.levels))...)
	}
	i := 0
	for ; i < len(key); i++ {
		var level = &self.levels[i+1]
		if (i == len(self.lastKey)) || key[i] != level.labels[len(level.labels)-1] {
			level.louds[len(self.levels[i+1].louds)-1] = true
			level.louds = append(level.louds, false)
			level.outs = append(level.outs, false)
			level.labels = append(level.labels, key[i])
			break
		}
	}
	for i++; i < len(key); i++ {
		var level = &self.levels[i+1]
		level.louds = append(level.louds, true, false)
		level.outs = append(level.outs, false)
		level.labels = append(level.labels, key[i])
	}
	self.levels[len(key)+1].louds = append(self.levels[len(key)+1].louds, true)
	self.levels[len(key)].outs[len(self.levels[len(key)].outs)-1] = true
	self.lastKey = make([]uint16, len(key))
	copy(self.lastKey, key)
	return nil
}

func (self *LoudsTrieBuilder) Build() *LoudsTrie {
	louds := []bool{}
	outs := []bool{}
	labels := []uint16{}
	for _, level := range self.levels {
		louds = append(louds, level.louds...)
		outs = append(outs, level.outs...)
		labels = append(labels, level.labels...)
	}
	louds = louds[:len(louds)-1]
	words := make([]uint64, (len(louds)+63)/64)
	for i := 0; i < len(louds); i++ {
		if louds[i] {
			words[i/64] |= 1 << (i % 64)
		}
	}
	var bitVector = NewBitVector(words, uint32(len(louds)))
	return NewLoudsTrie(bitVector, labels)
}

func CompareUtf16String(a []uint16, b []uint16) int {
	var min = len(a)
	if min > len(b) {
		min = len(b)
	}
	for i := 0; i < min; i++ {
		if a[i] > b[i] {
			return 1
		} else if a[i] < b[i] {
			return -1
		}
	}
	if len(a) == len(b) {
		return 0
	} else if len(a) > len(b) {
		return 1
	} else {
		return -1
	}
}

func BuildLoudsTrie(keys [][]uint16) (*LoudsTrie, []uint32, error) {
	for i := 0; i < len(keys)-1; i++ {
		if CompareUtf16String(keys[i], keys[i+1]) >= 0 {
			return nil, nil, errors.New("keys need be ordered")
		}
	}
	var nodes = make([]uint32, len(keys))
	for i := 0; i < len(nodes); i++ {
		nodes[i] = 1
	}
	var cursor = 0
	var currentNode uint32 = 1
	var edges = []uint16{0x30, 0x30} // TODO: '0'で穴埋めを'\0'にするか、なくす
	var louds = NewBitList()
	louds.Add(true)
	for true {
		var lastChar uint16 = 0
		var lastParent uint32 = 0
		var restKeys uint32 = 0
		for i := 0; i < len(keys); i++ {
			if len(keys[i]) < cursor {
				continue
			}
			if len(keys[i]) == cursor {
				louds.Add(false)
				lastParent = nodes[i]
				continue
			}
			var currentChar = keys[i][cursor]
			var currentParent = nodes[i]
			if lastParent != currentParent {
				louds.Add(false)
				louds.Add(true)
				edges = append(edges, currentChar)
				currentNode = currentNode + 1
			} else if lastChar != currentChar {
				louds.Add(true)
				edges = append(edges, currentChar)
				currentNode = currentNode + 1
			}
			nodes[i] = currentNode
			lastChar = currentChar
			lastParent = currentParent
			restKeys++
		}
		if restKeys == 0 {
			break
		}
		cursor++
	}
	var bitVector = NewBitVector(louds.Words, uint32(louds.Size))
	return NewLoudsTrie(bitVector, edges), nodes, nil
}
