package migemo

// DoubleArrayBuilder は、DoubleArrayを生成する構造体
type DoubleArrayBuilder struct {
	base    []int16
	check   []int16
	keys    []string
	indices []int16
}

// BuildDoubleArray は、文字列配列keysからDoubleArrayを生成する
func BuildDoubleArray(keys []string) *DoubleArray {
	indices := make([]int16, len(keys))
	builder := NewDoubleArrayBuilder(keys, indices)
	builder.traverse(0, 0, len(keys), 0)
	f := func(e uint8) int {
		return int(e)
	}
	return NewDoubleArray(builder.base, builder.check, f, 128)
}

// NewDoubleArrayBuilder は、DoubleArrayBuilderを初期化する
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

func (builder *DoubleArrayBuilder) traverse(index int16, start int, end int, offset int) {
	if len(builder.keys[start]) == offset {
		builder.indices[start] = index
		start++
		if start == end {
			return
		}
	}

	// enumerate children chars
	var childrenChars = make([]uint8, 0)
	lastChar := uint8(0)
	for i := start; i < end; i++ {
		currentChar := builder.keys[i][offset]
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
			builder.ensureDoubleArray(a)
			if builder.check[a] >= 0 {
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
	builder.base[index] = childrenOffset
	for i := 0; i < len(childrenChars); i++ {
		a := childrenOffset + int16(childrenChars[i])
		builder.check[a] = index
	}

	// visit children recursively
	lastChar = builder.keys[start][offset]
	startPos := start
	for i := start; i < end; i++ {
		currentChar := builder.keys[i][offset]
		if currentChar != lastChar {
			a := childrenOffset + int16(lastChar)
			builder.traverse(a, startPos, i, offset+1)
			startPos = i
			lastChar = currentChar
		}
	}
	builder.traverse(childrenOffset+int16(lastChar), startPos, end, offset+1)
}

func (builder *DoubleArrayBuilder) ensureDoubleArray(a int16) {
	if int16(len(builder.base)) <= a {
		count := int(a) - len(builder.base) + 1
		for i := 0; i < count; i++ {
			builder.base = append(builder.base, -1)
			builder.check = append(builder.check, -1)
		}
	}
}
