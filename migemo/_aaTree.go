package migemo

type AATreeNode struct {
	value rune
	left  *AATreeNode
	right *AATreeNode
	level int
}

func (this *AATreeNode) successor() *AATreeNode {
	t := this.right
	for t.left != nil {
		t = t.left
	}
	return t
}

func (this *AATreeNode) predecessor() *AATreeNode {
	t := this.left
	for t.left != nil {
		t = t.left
	}
	for t.right != nil {
		t = t.right
	}
	return t
}

type AATree struct {
	root *AATreeNode
}

func NewAATree() AATree {
	return AATree{
		root: nil,
	}
}

func (this *AATree) IsEmpty() bool {
	return this.root == nil
}
func (this *AATree) Clear() {
	this.root = nil
}

func skew(t *AATreeNode) *AATreeNode {
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

func split(t *AATreeNode) *AATreeNode {
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

func insert(x rune, t *AATreeNode) *AATreeNode {
	if t == nil {
		return &AATreeNode{
			value: x,
			level: 1,
			left:  nil,
			right: nil,
		}
	} else if x < t.value {
		t.left = insert(x, t.left)
	} else if x > t.value {
		t.right = insert(x, t.right)
	}
	t = skew(t)
	t = split(t)
	return t
}

func (this *AATree) Add(key rune) {
	this.root = insert(key, this.root)
}

func (this *AATree) Remove(key rune) {
	this.root = delete(key, this.root)
}

func (this *AATree) Contains(key rune) bool {
	t := this.root
	for t != nil {
		if t.value == key {
			return true
		} else if t.value < key {
			t = t.right
		} else if t.value > key {
			t = t.left
		}
	}
	return false
}

func delete(x rune, t *AATreeNode) *AATreeNode {
	if t == nil {
		return t
	} else if x > t.value {
		t.right = delete(x, t.right)
	} else if x < t.value {
		t.left = delete(x, t.left)
	} else {
		if t.left == nil && t.right == nil {
			return nil
		} else if t.left == nil {
			l := t.successor()
			t.right = delete(l.value, t.right)
			t.value = l.value
		} else {
			l := t.predecessor()
			t.left = delete(l.value, t.left)
			t.value = l.value
		}
	}
	t = decrease_level(t)
	t = skew(t)
	t.right = skew(t.right)
	if t.right != nil {
		t.right.right = skew(t.right.right)
	}
	t = split(t)
	t.right = split(t.right)
	return t
}

func decrease_level(t *AATreeNode) *AATreeNode {
	should_be := 0
	if t.left == nil || t.right == nil {
		should_be = 1
	} else if t.left.level < t.right.level {
		should_be = t.left.level + 1
	} else {
		should_be = t.right.level + 1
	}
	if should_be < t.level {
		t.level = should_be
		if should_be < t.right.level {
			t.right.level = should_be
		}
	}
	return t
}
