package migemo

type TernaryRegexNode struct {
	value uint16
	child *TernaryRegexNode
	left  *TernaryRegexNode
	right *TernaryRegexNode
	level int
}

func (this *TernaryRegexNode) successor() *TernaryRegexNode {
	t := this.right
	for t.left != nil {
		t = t.left
	}
	return t
}

func (this *TernaryRegexNode) predecessor() *TernaryRegexNode {
	t := this.left
	for t.left != nil {
		t = t.left
	}
	for t.right != nil {
		t = t.right
	}
	return t
}

type TernaryRegexGenerator struct {
	root             *TernaryRegexNode
	operator         RegexOperator
	escapeCharacters [2]uint64
}

func initializeEscapeCharacters() [2]uint64 {
	const ESCAPE = "\\.[]{}()*+-?^$|"
	bits := [2]uint64{}
	for _, r := range ESCAPE {
		var c = uint16(r)
		bits[c/64] |= 1 << (c % 64)
	}
	return bits
}

func (this *TernaryRegexGenerator) isEscapeCharacter(c uint16) bool {
	if c < 128 {
		return (this.escapeCharacters[c/64]>>(c%64))&1 == 1
	}
	return false
}

func NewTerminaryRegexGenerator(operator RegexOperator) *TernaryRegexGenerator {
	return &TernaryRegexGenerator{
		root:             nil,
		operator:         operator,
		escapeCharacters: initializeEscapeCharacters(),
	}
}

func skew(t *TernaryRegexNode) *TernaryRegexNode {
	if t == nil {
		return nil
	} else if t.left == nil {
		return t
	} else if t.left.level == t.level {
		l := t.left
		t.left = l.right
		l.right = t
		return l
	} else {
		return t
	}
}

func split(t *TernaryRegexNode) *TernaryRegexNode {
	if t == nil {
		return nil
	} else if t.right == nil || t.right.right == nil {
		return t
	} else if t.level == t.right.right.level {
		r := t.right
		t.right = r.left
		r.left = t
		r.level = r.level + 1
		return r
	} else {
		return t
	}
}

func insert(x uint16, t *TernaryRegexNode) (*TernaryRegexNode, *TernaryRegexNode, bool) {
	var r *TernaryRegexNode
	var inserted bool = false
	if t == nil {
		r = &TernaryRegexNode{
			value: x,
			level: 1,
			left:  nil,
			right: nil,
		}
		return r, r, true
	} else if x < t.value {
		t.left, r, inserted = insert(x, t.left)
	} else if x > t.value {
		t.right, r, inserted = insert(x, t.right)
	} else {
		return t, t, false
	}
	t = skew(t)
	t = split(t)
	return t, r, inserted
}

func (this *TernaryRegexGenerator) Add(word []uint16) {
	if len(word) == 0 {
		return
	}
	this.root = add(this.root, word, 0)
}

func add(node *TernaryRegexNode, word []uint16, offset int) *TernaryRegexNode {
	if offset < len(word) {
		node, target, inserted := insert(word[offset], node)
		if inserted || target.child != nil {
			target.child = add(target.child, word, offset+1)
		}
		return node
	} else {
		return nil
	}
}

func traverseSiblings(node *TernaryRegexNode, f func(node *TernaryRegexNode)) {
	if node != nil {
		traverseSiblings(node.left, f)
		f(node)
		traverseSiblings(node.right, f)
	}
}

func (this *TernaryRegexGenerator) generate(node *TernaryRegexNode, buf *[]uint16) {
	var brother = 0
	var haschild = 0
	traverseSiblings(node, func(node *TernaryRegexNode) {
		brother++
		if node.child != nil {
			haschild++
		}
	})
	var nochild = brother - haschild

	if brother > 1 && haschild > 0 {
		*buf = append(*buf, this.operator.beginGroup...)
	}

	if nochild > 0 {
		if nochild > 1 {
			*buf = append(*buf, this.operator.beginClass...)
		}
		traverseSiblings(node, func(node *TernaryRegexNode) {
			if node.child != nil {
				return
			}
			if this.isEscapeCharacter(node.value) {
				*buf = append(*buf, 92)
			}
			*buf = append(*buf, node.value)
		})
		if nochild > 1 {
			*buf = append(*buf, this.operator.endClass...)
		}
	}

	if haschild > 0 {
		if nochild > 0 {
			*buf = append(*buf, this.operator.or...)
		}
		traverseSiblings(node, func(node *TernaryRegexNode) {
			if node.child != nil {
				if this.isEscapeCharacter(node.value) {
					*buf = append(*buf, 92)
				}
				*buf = append(*buf, node.value)
				if this.operator.newline != nil { // TODO: always true
					*buf = append(*buf, this.operator.newline...)
				}
				this.generate(node.child, buf)
				if haschild > 1 {
					*buf = append(*buf, this.operator.or...)
				}
			}
		})
		if haschild > 1 {
			*buf = (*buf)[:len(*buf)-1]
		}
	}
	if brother > 1 && haschild > 0 {
		*buf = append(*buf, this.operator.endGroup...)
	}
}

func (this *TernaryRegexGenerator) Generate() []uint16 {
	if this.root == nil {
		return []uint16{}
	} else {
		buffer := []uint16{}
		this.generate(this.root, &buffer)
		return buffer
	}
}
