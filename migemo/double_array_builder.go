package migemo

type DoubleArrayBuilder struct {
	base    []int16
	check   []int16
	keys    []string
	indices []int16
}

func BuildDoubleArray(keys []string) *DoubleArray {
	indices := make([]int16, len(keys))
	builder := NewDoubleArrayBuilder(keys, indices)
	builder.traverse(0, 0, len(keys), 0)
	f := func(e uint8) int {
		return int(e)
	}
	return NewDoubleArray(builder.base, builder.check, f, 128)
}

func NewDoubleArrayBuilder(keys []string, indices []int16) *DoubleArrayBuilder {
	base := make([]int16, 16)
	for i := 0; i < len(base); i++ {
		base[i] = -1
	}
	check := make([]int16, 16)
	for i := 0; i < len(check); i++ {
		check[i] = -1
	}
	return &DoubleArrayBuilder{
		base,
		check,
		keys,
		indices,
	}
}

func (this *DoubleArrayBuilder) traverse(index int16, start int, end int, offset int) {
	if len(this.keys[start]) == offset {
		this.indices[start] = index
		start++
		if start == end {
			return
		}
	}

	// enumerate children chars
	var childrenChars = make([]uint8, 0)
	lastChar := uint8(0)
	for i := start; i < end; i++ {
		currentChar := this.keys[i][offset]
		if currentChar != lastChar {
			lastChar = currentChar
			childrenChars = append(childrenChars, currentChar)
		}
	}

	// find children offset
	childrenOffset := int16(0)
	for true {
		conflict := false
		for i := 0; i < len(childrenChars); i++ {
			a := childrenOffset + int16(childrenChars[i])
			this.ensureDoubleArray(a)
			if this.check[a] >= 0 {
				conflict = true
				break
			}
		}
		if conflict {
			childrenOffset++
		} else {
			break
		}
	}

	// mark to base and check
	this.base[index] = childrenOffset
	for i := 0; i < len(childrenChars); i++ {
		a := childrenOffset + int16(childrenChars[i])
		this.check[a] = index
	}

	// visit children recursively
	lastChar = this.keys[start][offset]
	startPos := start
	for i := start; i < end; i++ {
		currentChar := this.keys[i][offset]
		if currentChar != lastChar {
			a := childrenOffset + int16(lastChar)
			this.traverse(a, startPos, i, offset+1)
			startPos = i
			lastChar = currentChar
		}
	}
	this.traverse(childrenOffset+int16(lastChar), startPos, end, offset+1)
}

func (this *DoubleArrayBuilder) ensureDoubleArray(a int16) {
	if int16(len(this.base)) <= a {
		count := int(a) - len(this.base) + 1
		for i := 0; i < count; i++ {
			this.base = append(this.base, -1)
			this.check = append(this.check, -1)
		}
	}
}
