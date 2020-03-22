package main

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"github.com/oguna/gomigemo/migemo"
)

func main() {
	f, err := os.Open("migemo-compact-dict")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	dict := migemo.NewCompactDictionary(buf)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		s := stdin.Text()
		if len(s) == 0 {
			break
		}
		r := migemo.Query(s, dict)
		fmt.Println(r)
	}
}