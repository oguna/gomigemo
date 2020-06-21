package migemo

import "unicode/utf16"

// RegexNode は、RegexGeneratorに格納するノード
type RegexNode struct {
	code  uint16
	child *RegexNode
	next  *RegexNode
}

// RegexGenerator は、二重連鎖木による正規表現を生成する構造体
type RegexGenerator struct {
	operator         RegexOperator
	root             *RegexNode
	escapeCharacters [2]uint64
}

// RegexOperator は、正規表現の記号を格納する
type RegexOperator struct {
	or         []uint16
	beginGroup []uint16
	endGroup   []uint16
	beginClass []uint16
	endClass   []uint16
	newline    []uint16
}

// NewRegexOperator は、RegexOperatorを初期化する
func NewRegexOperator(or string, beginGroup string, endGroup string, beginClass string, endClass string, newline string) *RegexOperator {
	return &RegexOperator{
		or:         utf16.Encode([]rune(or)),
		beginGroup: utf16.Encode([]rune(beginGroup)),
		endGroup:   utf16.Encode([]rune(endGroup)),
		beginClass: utf16.Encode([]rune(beginClass)),
		endClass:   utf16.Encode([]rune(endClass)),
		newline:    utf16.Encode([]rune(newline)),
	}
}

// NewRegexNode は、RegexNodeを初期化する
func NewRegexNode(code uint16) *RegexNode {
	return &RegexNode{
		code:  code,
		child: nil,
		next:  nil,
	}
}

// NewRegexGenerator は、RegexGeneratorを初期化する
func NewRegexGenerator(operator RegexOperator) *RegexGenerator {
	const ESCAPE = "\\.[]{}()*+-?^$|"
	bits := [2]uint64{}
	for _, r := range ESCAPE {
		var c = uint16(r)
		bits[c/64] |= 1 << (c % 64)
	}
	return &RegexGenerator{
		root:             nil,
		operator:         operator,
		escapeCharacters: bits,
	}
}

func (generator *RegexGenerator) isEscapeCharacter(c uint16) bool {
	if c < 128 {
		return (generator.escapeCharacters[c/64]>>(c%64))&1 == 1
	}
	return false
}

func _add(node *RegexNode, word []uint16, offset int) *RegexNode {
	if node == nil {
		if offset >= len(word) {
			return nil
		}
		node = NewRegexNode(word[offset])
		if offset < len(word)-1 {
			node.child = _add(nil, word, offset+1)
		}
		return node
	}
	var thisNode = node
	var code = word[offset]
	if code < node.code {
		var newNode = NewRegexNode(code)
		newNode.next = node
		node = newNode
		if offset < len(word) {
			node.child = _add(nil, word, offset+1)
		}
		thisNode = node
	} else {
		for node.next != nil && node.next.code <= code {
			node = node.next
		}
		if node.code == code {
			if node.child == nil {
				return thisNode
			}
		} else {
			var newNode = NewRegexNode(code)
			newNode.next = node.next
			node.next = newNode
			node = newNode
		}
		if len(word) == offset+1 {
			node.child = nil
			return thisNode
		}
		node.child = _add(node.child, word, offset+1)
	}
	return thisNode
}

// Add は、RegexGeneratorに単語を追加する
func (generator *RegexGenerator) Add(word []uint16) {
	if len(word) == 0 {
		return
	}
	generator.root = _add(generator.root, word, 0)
}

func (generator *RegexGenerator) generateStub(node *RegexNode) []uint16 {
	var brother = 1
	var haschild = 0
	var buf []uint16
	for iter := node; iter != nil; iter = iter.next {
		if iter.next != nil {
			brother++
		}
		if iter.child != nil {
			haschild++
		}
	}
	var nochild = brother - haschild

	if brother > 1 && haschild > 0 {
		buf = append(buf, generator.operator.beginGroup...)
	}

	if nochild > 0 {
		if nochild > 1 {
			buf = append(buf, generator.operator.beginClass...)
		}
		for iter := node; iter != nil; iter = iter.next {
			if iter.child == nil {
				if generator.isEscapeCharacter(iter.code) {
					buf = append(buf, 92)
				}
				buf = append(buf, iter.code)
			}
		}
		if nochild > 1 {
			buf = append(buf, generator.operator.endClass...)
		}
	}

	if haschild > 0 {
		if nochild > 0 {
			buf = append(buf, generator.operator.or...)
		}
		for iter := node; iter != nil; iter = iter.next {
			if iter.child != nil {
				if generator.isEscapeCharacter(iter.code) {
					buf = append(buf, 92)
				}
				buf = append(buf, iter.code)
				if generator.operator.newline != nil { // TODO: always true
					buf = append(buf, generator.operator.newline...)
				}
				buf = append(buf, generator.generateStub(iter.child)...)
				if haschild > 1 && iter.next != nil {
					buf = append(buf, generator.operator.or...)
				}
			}
		}
	}
	if brother > 1 && haschild > 0 {
		buf = append(buf, generator.operator.endGroup...)
	}
	return buf
}

// Generate は、RegexGenertorに追加した単語から、正規表現を生成する
func (generator *RegexGenerator) Generate() []uint16 {
	if generator.root == nil {
		return []uint16{}
	}
	return generator.generateStub(generator.root)
}
