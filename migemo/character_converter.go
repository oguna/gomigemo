package migemo

import "strings"

var fullwidthPunctuationList = [4]uint8{0x02, 0x0c, 0x0d, 0x01}
var fullwidthKatakanaList = [59]uint8{0xFB, 0xF2, 0xA1, 0xA3, 0xA5, 0xA7, 0xA9, 0xE3, 0xE5, 0xE7, 0xC3, 0xFC, 0xA2, 0xA4, 0xA6, 0xA8, 0xAA, 0xAB, 0xAD, 0xAF, 0xB1, 0xB3, 0xB5, 0xB7, 0xB9, 0xBB, 0xBD, 0xBF, 0xC1, 0xC4, 0xC6, 0xC8, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD2, 0xD5, 0xD8, 0xDB, 0xDE, 0xDF, 0xE0, 0xE1, 0xE2, 0xE4, 0xE6, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEF, 0xF3, 0x99, 0x9A}
var dakuonBits = [6]uint16{0b0101000000000000, 0b0101010101010101, 0b0000001010100101, 0b0001001001001001, 0, 0b0000000000010000}
var handakuonBits = [6]uint16{0, 0, 0, 0b0010010010010010, 0, 0}
var halfwidthKatakanaList = [96]uint16{0, 0x67, 0x71, 0x68, 0x72, 0x69, 0x73, 0x6A, 0x74, 0x6B, 0x75, 0x76, 0x76, 0x77, 0x77, 0x78, 0x78, 0x79, 0x79, 0x7A, 0x7A, 0x7B, 0x7B, 0x7C, 0x7C, 0x7D, 0x7D, 0x7E, 0x7E, 0x7F, 0x7F, 0x80, 0x80, 0x81, 0x81, 0x6F, 0x82, 0x82, 0x83, 0x83, 0x84, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8A, 0x8A, 0x8B, 0x8B, 0x8B, 0x8C, 0x8C, 0x8C, 0x8D, 0x8D, 0x8D, 0x8E, 0x8E, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x6C, 0x94, 0x6D, 0x95, 0x6E, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0, 0x9C, 0, 0, 0x66, 0x9D, 0x73, 0, 0, 0, 0, 0, 0, 0x65, 0x70, 0, 0, 0}

// ConvertHan2Zen は、半角から全角へ文字列を変更する
func ConvertHan2Zen(source string) string {
	var sb strings.Builder
	sb.Grow(len(source))
	for _, c := range []rune(source) {
		if c == 0x20 {
			// 半角スペース(migemoの仕様上スペースは入力されないが、関数の汎用性のため)
			c = 0x3000
		} else if 0x21 <= c && c <= 0x7e {
			// 半角ASCII印字可能文字
			c = 0xFF01 - 0x21 + c
		} else if 0xFF61 <= c && c <= 0xFF64 {
			// 半角句読点(｡｢｣､)
			c = rune(fullwidthPunctuationList[c-0xFF61]) + 0x3000
		} else if 0xFF65 <= c && c <= 0xFF9F {
			// 半角カタカナ
			c = rune(fullwidthKatakanaList[c-0xFF65]) + 0x3000
		}
		sb.WriteRune(c)
	}
	return sb.String()
}

// ConvertZen2Han は、全角から半角へ文字列を変更する
func ConvertZen2Han(source string) string {
	var sb strings.Builder
	sb.Grow(len(source))
	for _, c := range []rune(source) {
		if 0xFF01 <= c && c <= 0xFF5E {
			// 全角ASCII印字可能文字
			c = c - 0xFF00 + 0x0020
			sb.WriteRune(c)
		} else if 0x3000 <= c && c <= 0x300F {
			// 全角スペースと句読点
			if c == 0x3000 {
				c = 0x20
			} else if c == 0x3002 {
				c = 0xFF61
			} else if c == 0x300C {
				c = 0xFF62
			} else if c == 0x300D {
				c = 0xFF63
			} else if c == 0x3001 {
				c = 0xFF64
			}
			sb.WriteRune(c)
		} else if 0x3099 == c {
			c = '\uFF9E' // 全角カタカナ濁音
			sb.WriteRune(c)
		} else if 0x309A == c {
			c = '\uFF9F' // 全角カタカナ半濁音
			sb.WriteRune(c)
		} else if 0x30A0 <= c && c <= 0x30FF {
			// 全角カタカナ
			x1 := uint8((c >> 0) & 0xF)
			x2 := uint8((c >> 4) & 0xF)
			isDakuon := (dakuonBits[x2-10] & (uint16(1) << x1)) > 0
			isHandakuon := (handakuonBits[x2-10] & (uint16(1) << x1)) > 0
			c_ := halfwidthKatakanaList[c-0x30A0]
			if c_ != 0 {
				c = rune(c_) + 0xFF00
			}
			sb.WriteRune(c)
			if isDakuon {
				sb.WriteRune('ﾞ')
			}
			if isHandakuon {
				sb.WriteRune('ﾟ')
			}
		}
	}
	return sb.String()
}

// ConvertHira2Kata は、ひらがなからカタカナへ文字列を変更する
func ConvertHira2Kata(source string) string {
	var sb = []rune(source)
	for i, c := range sb {
		if 'ぁ' <= c && c <= 'ん' {
			sb[i] = rune(c - 'ぁ' + 'ァ')
		}
	}
	return string(sb)
}
