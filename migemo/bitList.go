package migemo

import "fmt"

type BitList struct {
	Words []uint64
	Size  int
}

func NewBitList() *BitList {
	return &BitList{
		Words: make([]uint64, 0, 8),
		Size:  0,
	}
}

func (self *BitList) Add(value bool) {
	if len(self.Words) < (self.Size+1+63)/64 {
		self.Words = append(self.Words, 0)
	}
	self.Set(self.Size, value)
	self.Size++
}

func (self *BitList) Set(pos int, value bool) {
	if self.Size < pos {
		panic(fmt.Sprintf("index out of range [%d] with length %d", pos, self.Size))
	}
	if value {
		self.Words[pos/64] |= 1 << (pos % 64)
	} else {
		self.Words[pos/64] &= ^(uint64(1) << (pos % 64))
	}
}

func (self *BitList) Get(pos int) bool {
	if self.Size < pos {
		panic(fmt.Sprintf("index out of range [%d] with length %d", pos, self.Size))
	}
	return self.Words[pos/64]>>(pos%64) == 1
}
