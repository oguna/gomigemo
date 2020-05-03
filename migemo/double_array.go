package migemo

// DoubleArray は、高速に検索可能なトライ木
type DoubleArray struct {
	base          []int16
	check         []int16
	charConverter func(uint8) int
	charSize      int
}

func NewDoubleArray(base []int16, check []int16, code func(uint8) int, charSize int) *DoubleArray {
	return &DoubleArray{
		base,
		check,
		code,
		charSize,
	}
}

func (this *DoubleArray) traverse(n int16, k int) int16 {
	m := this.base[n] + int16(k)
	if this.check[m] == n {
		return m
	} else {
		return -1
	}
}

func (this *DoubleArray) Lookup(str string) int16 {
	if len(str) == 0 {
		return 0
	}
	n := int16(0)
	for i := 0; i < len(str); i++ {
		c := this.charConverter(str[i])
		if c < 1 {
			panic("")
		}
		n = this.traverse(n, c)
		if n == -1 {
			return -1
		}
	}
	return n
}

func (this *DoubleArray) CommonPrefixSearch(key string, f func(node int16)) {
	index := int16(0)
	offset := 0
	for index != -1 {
		lastIndex := index
		if offset == len(key) {
			index = -1
		} else {
			c := this.charConverter(key[offset])
			index = this.traverse(index, c)
			offset++
		}
		f(lastIndex)
	}
}

func (this *DoubleArray) PredictiveSearch(key string, f func(node int16)) {
	n := this.Lookup(key)
	if n == -1 {
		return
	} else {
		this.visitRecursive(n, f)
	}
}

func (this *DoubleArray) visitRecursive(n int16, f func(node int16)) {
	f(n)
	for i := 0; i < this.charSize; i++ {
		m := int(this.base[n]) + i + 1
		if m >= len(this.check) {
			return
		}
		if this.check[m] == n {
			this.visitRecursive(int16(m), f)
		}
	}
}
