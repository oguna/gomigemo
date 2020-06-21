package migemo_test

import (
	"testing"

	"github.com/oguna/gomigemo/migemo"
)

func TestMigemoParser_1(t *testing.T) {
	query := "toukyouOosaka nagoyaFUKUOKAhokkaido "
	words := []string{"toukyou", "Oosaka", "nagoya", "FUKUOKA", "hokkaido", ""}
	iter := migemo.NewMigemoParser(query)
	for _, w := range words {
		if w != iter.Next() {
			t.Error()
		}
	}
}

func TestMigemoParser_2(t *testing.T) {
	query := "a"
	words := []string{"a", ""}
	iter := migemo.NewMigemoParser(query)
	for _, w := range words {
		if w != iter.Next() {
			t.Error()
		}
	}
}

func TestMigemoParser_3(t *testing.T) {
	query := "A"
	words := []string{"A", ""}
	iter := migemo.NewMigemoParser(query)
	for _, w := range words {
		if w != iter.Next() {
			t.Error()
		}
	}
}

func TestMigemoParser_4(t *testing.T) {
	query := "あ"
	words := []string{"あ", ""}
	iter := migemo.NewMigemoParser(query)
	for _, w := range words {
		if w != iter.Next() {
			t.Error()
		}
	}
}
