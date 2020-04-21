package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestAATree(t *testing.T) {
	tree := migemo.NewAATree()
	tree.Add('a')
	if !tree.Contains('a') {
		t.Error()
	}
	tree.Add('b')
	if !tree.Contains('b') {
		t.Error()
	}
	tree.Add('c')
	if !tree.Contains('c') {
		t.Error()
	}
	tree.Add('d')
	if !tree.Contains('d') {
		t.Error()
	}
	if tree.Contains('e') {
		t.Error()
	}
	tree.Remove('a')
	if tree.Contains('e') {
		t.Error()
	}
	tree.Add('a')
	if !tree.Contains('a') {
		t.Error()
	}
}
