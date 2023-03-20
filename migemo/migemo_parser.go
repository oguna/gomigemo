package migemo

type MigemoParser struct {
	Query  string
	Cursor int
}

func NewMigemoParser(query string) MigemoParser {
	return MigemoParser{
		Query:  query,
		Cursor: 0,
	}
}

func (parser *MigemoParser) Next() string {
	// カーソルが終端なら終了
	if len(parser.Query) <= parser.Cursor {
		return ""
	}
	// 空白をスキップ
	for parser.Cursor < len(parser.Query) && parser.Query[parser.Cursor] == 0x20 {
		parser.Cursor++
	}
	// スキップした結果カーソルが終端に達していないか確認
	if len(parser.Query) <= parser.Cursor {
		return ""
	}
	// 単語の先頭文字の種類で場合分け
	start := parser.Cursor
	c := parser.Query[parser.Cursor]
	if 0x41 <= c && c <= 0x5a {
		// 大文字なら、大文字または小文字が続くまで
		if parser.Cursor+1 < len(parser.Query) {
			nextChar := parser.Query[parser.Cursor+1]
			if 0x41 <= nextChar && nextChar <= 0x5a {
				for 0x41 <= nextChar && nextChar <= 0x5a {
					parser.Cursor++
					if parser.Cursor+1 < len(parser.Query) {
						nextChar = parser.Query[parser.Cursor+1]
					} else {
						break
					}
				}
			} else if 0x61 <= nextChar && nextChar <= 0x7a {
				for 0x61 <= nextChar && nextChar <= 0x7a {
					parser.Cursor++
					if parser.Cursor+1 < len(parser.Query) {
						nextChar = parser.Query[parser.Cursor+1]
					} else {
						break
					}
				}
			}
		}
		parser.Cursor++
		return parser.Query[start:parser.Cursor]
	} else if 0x61 <= c && c <= 0x7a {
		// 小文字なら、小文字が続くまで
		if parser.Cursor+1 < len(parser.Query) {
			nextChar := parser.Query[parser.Cursor+1]
			for 0x61 <= nextChar && nextChar <= 0x7a {
				parser.Cursor++
				if parser.Cursor+1 < len(parser.Query) {
					nextChar = parser.Query[parser.Cursor+1]
				} else {
					break
				}
			}
		}
		parser.Cursor++
		return parser.Query[start:parser.Cursor]
	} else {
		// それ以外なら、空白に至るまで
		nextChar := uint8(0)
		if parser.Cursor+1 < len(parser.Query) {
			nextChar = parser.Query[parser.Cursor+1]
		} else {
			parser.Cursor++
			return parser.Query[start:parser.Cursor]
		}
		for nextChar != 0x20 {
			parser.Cursor++
			if parser.Cursor+1 < len(parser.Query) {
				nextChar = parser.Query[parser.Cursor+1]
			} else {
				break
			}
		}
		parser.Cursor++
		return parser.Query[start:parser.Cursor]
	}
}
