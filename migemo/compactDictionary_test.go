package migemo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"unicode/utf16"

	"github.com/oguna/gomigemo/migemo"
)

func TestCompactDictionary_1(t *testing.T) {
	f, err := os.Open("../migemo-compact-dict")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	dict := migemo.NewCompactDictionary(buf)
	list := []string{}
	fn := func(s []uint16) {
		list = append(list, string(utf16.Decode(s)))
	}
	dict.Search(utf16.Encode([]rune("きかい")), fn)
	for _, w := range list {
		if w == "機械" {
			return
		}
	}
	t.Error()
}

func TestCompactDictionary_2(t *testing.T) {
	f, err := os.Open("../migemo-compact-dict")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	dict := migemo.NewCompactDictionary(buf)
	list := []string{}
	fn := func(s []uint16) {
		list = append(list, string(utf16.Decode(s)))
	}
	dict.Search(utf16.Encode([]rune("し")), fn)
	for _, w := range list {
		if w == "市" {
			return
		}
	}
	t.Error()
}

func TestCompactDictionary_3(t *testing.T) {
	f, err := os.Open("../migemo-compact-dict")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	dict := migemo.NewCompactDictionary(buf)
	list := []string{}
	fn := func(s []uint16) {
		list = append(list, string(utf16.Decode(s)))
	}
	dict.PredictiveSearch(utf16.Encode([]rune("し")), fn)
	for _, w := range list {
		if w == "三味" {
			return
		}
	}
	t.Error()
}
