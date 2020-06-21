package migemo

import "fmt"

// BitList は、ビット配列を効率的に格納する伸長可能な構造体
type BitList struct {
	Words []uint64
	Size  int
}

// NewBitList は、BitListを長さ=0で初期化する
func NewBitList() *BitList {
	return &BitList{
		Words: make([]uint64, 0, 8),
		Size:  0,
	}
}

// NewBitListWithSize は、BitListを長さ=sizeで初期化する
func NewBitListWithSize(size int) *BitList {
	return &BitList{
		Words: make([]uint64, (size+63)/64),
		Size:  size,
	}
}

// Add は、ビットをリストの末尾に追加する
func (bitList *BitList) Add(value bool) {
	if len(bitList.Words) < (bitList.Size+1+63)/64 {
		bitList.Words = append(bitList.Words, 0)
	}
	bitList.Set(bitList.Size, value)
	bitList.Size++
}

// Set は、指定したposのビット値をvalueに設定する
func (bitList *BitList) Set(pos int, value bool) {
	if bitList.Size < pos {
		panic(fmt.Sprintf("index out of range [%d] with length %d", pos, bitList.Size))
	}
	if value {
		bitList.Words[pos/64] |= uint64(1) << (pos % 64)
	} else {
		bitList.Words[pos/64] &= ^(uint64(1) << (pos % 64))
	}
}

// Get は、指定したposのビット値を取得する
func (bitList *BitList) Get(pos int) bool {
	if bitList.Size < pos {
		panic(fmt.Sprintf("index out of range [%d] with length %d", pos, bitList.Size))
	}
	return (bitList.Words[pos/64]>>(pos%64))&1 == 1
}
