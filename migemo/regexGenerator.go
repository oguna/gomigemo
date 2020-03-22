package migemo

import "unicode/utf16"

type RegexNode struct {
	code uint16
	child *RegexNode
	next *RegexNode
}

type RegexGenerator struct {
	or []uint16;
	beginGroup []uint16;
	endGroup []uint16;
	beginClass []uint16;
	endClass []uint16;
	newline []uint16;
	root *RegexNode;
	escapeCharacters map[uint16]struct{};
}

func NewRegexNode(code uint16) *RegexNode {
	return &RegexNode {
		code: code,
		child: nil,
		next: nil,
	}
}

func NewRegexGenerator() *RegexGenerator {
	const ESCAPE = "\\.[]{}()*+-?^$|"
	var escapeCharacters = make(map[uint16]struct{})
	for _, r := range ESCAPE {
		var c = uint16(r)
		escapeCharacters[c] = struct{}{}
	}
	return &RegexGenerator {
		or: utf16.Encode([]rune("|")),
		beginGroup: utf16.Encode([]rune("(")),
		endGroup: utf16.Encode([]rune(")")),
		beginClass: utf16.Encode([]rune("[")),
		endClass: utf16.Encode([]rune("]")),
		newline: utf16.Encode([]rune("")),
		escapeCharacters: escapeCharacters,
	}
}

func _add(node *RegexNode, word []uint16, offset int) *RegexNode {
	if node == nil {
		if (offset >= len(word)) {
			return nil;
		}
		node = NewRegexNode(word[offset]);
		if (offset < len(word) - 1) {
			node.child = _add(nil, word, offset + 1);
		}
		return node;
	}
	var thisNode = node;
	var code = word[offset];
	if code < node.code {
		var newNode = NewRegexNode(code);
		newNode.next = node;
		node = newNode;
		if (offset < len(word)) {
			node.child = _add(nil, word, offset + 1);
		}
		thisNode = node;
	} else {
		for node.next != nil && node.next.code <= code {
			node = node.next;
		}
		if node.code == code {
			if node.child == nil {
				return thisNode
			}
		} else {
			var newNode = NewRegexNode(code);
			newNode.next = node.next;
			node.next = newNode;
			node = newNode;
		}
		if (len(word) == offset + 1) {
			node.child = nil;
			return thisNode;
		}
		node.child = _add(node.child, word, offset + 1);
	}
	return thisNode;
}

func (this *RegexGenerator) Add(word []uint16) {
	if (len(word) == 0) {
		return;
	}
	this.root = _add(this.root, word, 0);
}

func (this *RegexGenerator) generateStub(node *RegexNode) []uint16 {
	var brother = 1;
	var haschild = 0;
	var buf []uint16;
	for tmp := node; tmp != nil; tmp = tmp.next {
		if (tmp.next != nil) {
			brother++;
		}
		if (tmp.child != nil) {
			haschild++;
		}
	}
	var nochild = brother - haschild;

	if (brother > 1 && haschild > 0) {
		buf = append(buf, this.beginGroup...)
	}

	if (nochild > 0) {
		if (nochild > 1) {
			buf = append(buf, this.beginClass...)
		}
		for tmp := node; tmp != nil; tmp = tmp.next {
			if (tmp.child != nil) {
				continue;
			}
			var _, ok = this.escapeCharacters[tmp.code]
			if (ok) {
				buf = append(buf, 92);
			}
			buf = append(buf, tmp.code);
		}
		if (nochild > 1) {
			buf = append(buf, this.endClass...);
		}
	}

	if (haschild > 0) {
		if (nochild > 0) {
			buf = append(buf, this.or...);
		}
		var tmp *RegexNode = nil;
		for tmp = node; tmp.child == nil; tmp = tmp.next {
		}
		for true {
			var _, ok = this.escapeCharacters[tmp.code]
			if (ok) {
				buf = append(buf, 92);
			}
			buf = append(buf, tmp.code);
			if this.newline != nil {
				buf = append(buf, this.newline...);
			}
			buf = append(buf, this.generateStub(tmp.child)...);
			for tmp = tmp.next; tmp != nil && tmp.child == nil; tmp = tmp.next {
			}
			if (tmp == nil) {
				break;
			}
			if (haschild > 1) {
				buf = append(buf, this.or...);
			}
		}
	}
	if (brother > 1 && haschild > 0) {
		buf = append(buf, this.endGroup...);
	}
	return buf;
}

func (this *RegexGenerator) Generate() []uint16 {
	if (this.root == nil) {
		return []uint16 {}
	} else {
		return this.generateStub(this.root);
	}
}