package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/oguna/gomigemo/migemo"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	d := flag.String("d", "migemo-compact-dict", "Use a file <dict> for dictionary.")
	w := flag.String("w", "", "Expand a <word> and soon exit.")
	q := flag.Bool("q", false, "Show no message except results.")
	v := flag.Bool("v", false, "Use vim style regexp.")
	e := flag.Bool("e", false, "Use emacs style regexp.")
	n := flag.Bool("n", false, "Don't use newline match.")
	//p := flag.Int("p", 0, "<port> number for HTTP server.")

	flag.Parse()

	var regex_operator *migemo.RegexOperator
	if *v {
		if *n {
			regex_operator = migemo.NewRegexOperator("\\|", "\\%(", "\\)", "[", "]", "")
		} else {
			regex_operator = migemo.NewRegexOperator("\\|", "\\%(", "\\)", "[", "]", "\\_s*")
		}
	} else if *e {
		if *n {
			regex_operator = migemo.NewRegexOperator("\\|", "\\(", "\\)", "[", "]", "")
		} else {
			regex_operator = migemo.NewRegexOperator("\\|", "\\(", "\\)", "[", "]", "\\s-*")
		}
	} else {
		regex_operator = migemo.NewRegexOperator("|", "(", ")", "[", "]", "")
	}

	f, err := os.Open(*d)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	dict := migemo.NewCompactDictionary(buf)

	/*
		// セキュリティソフトにバックドアと誤検知されるため、一時的に機能を削除
		if *p > 0 {
			http.HandleFunc("/migemo", func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					v := r.URL.Query()
					q, ok := v["q"]
					if !ok {
						return
					}
					for _, a := range q {
						regex := migemo.Query(a, dict, regex_operator)
						fmt.Fprintln(w, regex)
					}
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			})
			err := http.ListenAndServe(":"+strconv.Itoa(*p), nil)
			if err != nil {
				log.Fatal(err)
			}
		} else
	*/
	if len(*w) == 0 {
		stdin := bufio.NewScanner(os.Stdin)
		if !*q {
			fmt.Print("QUERY: ")
		}
		for stdin.Scan() {
			s := stdin.Text()
			if len(s) == 0 {
				break
			}
			r := migemo.Query(s, dict, regex_operator)
			if !*q {
				r = "PATTERN: " + r
			}
			fmt.Println(r)
			if !*q {
				fmt.Print("QUERY: ")
			}
		}
	} else {
		r := migemo.Query(*w, dict, regex_operator)
		fmt.Println(r)
	}
}
