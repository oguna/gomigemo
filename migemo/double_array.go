package migemo

// DoubleArray は、高速に検索可能なトライ木
type DoubleArray struct {
	base          []int16
	check         []int16
	charConverter func(uint8) int
	charSize      int
}

// NewDoubleArray は、DoubleArrayを初期化する
func NewDoubleArray(base []int16, check []int16, code func(uint8) int, charSize int) *DoubleArray {
	return &DoubleArray{
		base,
		check,
		code,
		charSize,
	}
}

func (doubleArray *DoubleArray) traverse(n int16, k int) int16 {
	m := doubleArray.base[n] + int16(k)
	if doubleArray.check[m] == n {
		return m
	}
	return -1
}

// Lookup は、指定した文字列のノード番号を返す
func (doubleArray *DoubleArray) Lookup(str string) int16 {
	if len(str) == 0 {
		return 0
	}
	n := int16(0)
	for i := 0; i < len(str); i++ {
		c := doubleArray.charConverter(str[i])
		if c < 1 {
			panic("")
		}
		n = doubleArray.traverse(n, c)
		if n == -1 {
			return -1
		}
	}
	return n
}

// CommonPrefixSearch は、指定した文字列のノードまでにたどる全てのノード番号を関数fに返す
func (doubleArray *DoubleArray) CommonPrefixSearch(key string, f func(node int16)) {
	index := int16(0)
	offset := 0
	for index != -1 {
		lastIndex := index
		if offset == len(key) {
			index = -1
		} else {
			c := doubleArray.charConverter(key[offset])
			index = doubleArray.traverse(index, c)
			offset++
		}
		f(lastIndex)
	}
}

// PredictiveSearch は、接頭辞keyが含まれている全てのノードのノード番号を関数fに返す
func (doubleArray *DoubleArray) PredictiveSearch(key string, f func(node int16)) {
	n := doubleArray.Lookup(key)
	if n == -1 {
		return
	}
	doubleArray.visitRecursive(n, f)
}

func (doubleArray *DoubleArray) visitRecursive(n int16, f func(node int16)) {
	f(n)
	for i := 0; i < doubleArray.charSize; i++ {
		m := int(doubleArray.base[n]) + i + 1
		if m >= len(doubleArray.check) {
			return
		}
		if doubleArray.check[m] == n {
			doubleArray.visitRecursive(int16(m), f)
		}
	}
}
