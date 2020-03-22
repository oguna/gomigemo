package migemo

import (
	"sort"
	"unicode/utf16"
)

type RomanEntry struct {
	roman    []uint16
	hiragana []uint16
	remain   int
	index    uint32
}

type RomajiProcessor struct {
	entries []*RomanEntry
	indexes []uint32
}

func NewRomanEntry(roman string, hiragana string, remain int) *RomanEntry {
	return &RomanEntry{
		roman:    utf16.Encode([]rune(roman)),
		hiragana: utf16.Encode([]rune(hiragana)),
		remain:   remain,
		index:    calculateIndex(utf16.Encode([]rune(roman))),
	}
}

func NewRomajiProcessor() *RomajiProcessor {
	var entries = []*RomanEntry{
		NewRomanEntry("-", "ー", 0),
		NewRomanEntry("~", "〜", 0),
		NewRomanEntry(".", "。", 0),
		NewRomanEntry(",", "、", 0),
		NewRomanEntry("z/", "・", 0),
		NewRomanEntry("z.", "…", 0),
		NewRomanEntry("z,", "‥", 0),
		NewRomanEntry("zh", "←", 0),
		NewRomanEntry("zj", "↓", 0),
		NewRomanEntry("zk", "↑", 0),
		NewRomanEntry("zl", "→", 0),
		NewRomanEntry("z-", "〜", 0),
		NewRomanEntry("z[", "『", 0),
		NewRomanEntry("z]", "』", 0),
		NewRomanEntry("[", "「", 0),
		NewRomanEntry("]", "」", 0),
		NewRomanEntry("va", "ゔぁ", 0),
		NewRomanEntry("vi", "ゔぃ", 0),
		NewRomanEntry("vu", "ゔ", 0),
		NewRomanEntry("ve", "ゔぇ", 0),
		NewRomanEntry("vo", "ゔぉ", 0),
		NewRomanEntry("vya", "ゔゃ", 0),
		NewRomanEntry("vyi", "ゔぃ", 0),
		NewRomanEntry("vyu", "ゔゅ", 0),
		NewRomanEntry("vye", "ゔぇ", 0),
		NewRomanEntry("vyo", "ゔょ", 0),
		NewRomanEntry("qq", "っ", 1),
		NewRomanEntry("vv", "っ", 1),
		NewRomanEntry("ll", "っ", 1),
		NewRomanEntry("xx", "っ", 1),
		NewRomanEntry("kk", "っ", 1),
		NewRomanEntry("gg", "っ", 1),
		NewRomanEntry("ss", "っ", 1),
		NewRomanEntry("zz", "っ", 1),
		NewRomanEntry("jj", "っ", 1),
		NewRomanEntry("tt", "っ", 1),
		NewRomanEntry("dd", "っ", 1),
		NewRomanEntry("hh", "っ", 1),
		NewRomanEntry("ff", "っ", 1),
		NewRomanEntry("bb", "っ", 1),
		NewRomanEntry("pp", "っ", 1),
		NewRomanEntry("mm", "っ", 1),
		NewRomanEntry("yy", "っ", 1),
		NewRomanEntry("rr", "っ", 1),
		NewRomanEntry("ww", "っ", 1),
		NewRomanEntry("www", "w", 2),
		NewRomanEntry("cc", "っ", 1),
		NewRomanEntry("kya", "きゃ", 0),
		NewRomanEntry("kyi", "きぃ", 0),
		NewRomanEntry("kyu", "きゅ", 0),
		NewRomanEntry("kye", "きぇ", 0),
		NewRomanEntry("kyo", "きょ", 0),
		NewRomanEntry("gya", "ぎゃ", 0),
		NewRomanEntry("gyi", "ぎぃ", 0),
		NewRomanEntry("gyu", "ぎゅ", 0),
		NewRomanEntry("gye", "ぎぇ", 0),
		NewRomanEntry("gyo", "ぎょ", 0),
		NewRomanEntry("sya", "しゃ", 0),
		NewRomanEntry("syi", "しぃ", 0),
		NewRomanEntry("syu", "しゅ", 0),
		NewRomanEntry("sye", "しぇ", 0),
		NewRomanEntry("syo", "しょ", 0),
		NewRomanEntry("sha", "しゃ", 0),
		NewRomanEntry("shi", "し", 0),
		NewRomanEntry("shu", "しゅ", 0),
		NewRomanEntry("she", "しぇ", 0),
		NewRomanEntry("sho", "しょ", 0),
		NewRomanEntry("zya", "じゃ", 0),
		NewRomanEntry("zyi", "じぃ", 0),
		NewRomanEntry("zyu", "じゅ", 0),
		NewRomanEntry("zye", "じぇ", 0),
		NewRomanEntry("zyo", "じょ", 0),
		NewRomanEntry("tya", "ちゃ", 0),
		NewRomanEntry("tyi", "ちぃ", 0),
		NewRomanEntry("tyu", "ちゅ", 0),
		NewRomanEntry("tye", "ちぇ", 0),
		NewRomanEntry("tyo", "ちょ", 0),
		NewRomanEntry("cha", "ちゃ", 0),
		NewRomanEntry("chi", "ち", 0),
		NewRomanEntry("chu", "ちゅ", 0),
		NewRomanEntry("che", "ちぇ", 0),
		NewRomanEntry("cho", "ちょ", 0),
		NewRomanEntry("cya", "ちゃ", 0),
		NewRomanEntry("cyi", "ちぃ", 0),
		NewRomanEntry("cyu", "ちゅ", 0),
		NewRomanEntry("cye", "ちぇ", 0),
		NewRomanEntry("cyo", "ちょ", 0),
		NewRomanEntry("dya", "ぢゃ", 0),
		NewRomanEntry("dyi", "ぢぃ", 0),
		NewRomanEntry("dyu", "ぢゅ", 0),
		NewRomanEntry("dye", "ぢぇ", 0),
		NewRomanEntry("dyo", "ぢょ", 0),
		NewRomanEntry("tsa", "つぁ", 0),
		NewRomanEntry("tsi", "つぃ", 0),
		NewRomanEntry("tse", "つぇ", 0),
		NewRomanEntry("tso", "つぉ", 0),
		NewRomanEntry("tha", "てゃ", 0),
		NewRomanEntry("thi", "てぃ", 0),
		NewRomanEntry("t'i", "てぃ", 0),
		NewRomanEntry("thu", "てゅ", 0),
		NewRomanEntry("the", "てぇ", 0),
		NewRomanEntry("tho", "てょ", 0),
		NewRomanEntry("t'yu", "てゅ", 0),
		NewRomanEntry("dha", "でゃ", 0),
		NewRomanEntry("dhi", "でぃ", 0),
		NewRomanEntry("d'i", "でぃ", 0),
		NewRomanEntry("dhu", "でゅ", 0),
		NewRomanEntry("dhe", "でぇ", 0),
		NewRomanEntry("dho", "でょ", 0),
		NewRomanEntry("d'yu", "でゅ", 0),
		NewRomanEntry("twa", "とぁ", 0),
		NewRomanEntry("twi", "とぃ", 0),
		NewRomanEntry("twu", "とぅ", 0),
		NewRomanEntry("twe", "とぇ", 0),
		NewRomanEntry("two", "とぉ", 0),
		NewRomanEntry("t'u", "とぅ", 0),
		NewRomanEntry("dwa", "どぁ", 0),
		NewRomanEntry("dwi", "どぃ", 0),
		NewRomanEntry("dwu", "どぅ", 0),
		NewRomanEntry("dwe", "どぇ", 0),
		NewRomanEntry("dwo", "どぉ", 0),
		NewRomanEntry("d'u", "どぅ", 0),
		NewRomanEntry("nya", "にゃ", 0),
		NewRomanEntry("nyi", "にぃ", 0),
		NewRomanEntry("nyu", "にゅ", 0),
		NewRomanEntry("nye", "にぇ", 0),
		NewRomanEntry("nyo", "にょ", 0),
		NewRomanEntry("hya", "ひゃ", 0),
		NewRomanEntry("hyi", "ひぃ", 0),
		NewRomanEntry("hyu", "ひゅ", 0),
		NewRomanEntry("hye", "ひぇ", 0),
		NewRomanEntry("hyo", "ひょ", 0),
		NewRomanEntry("bya", "びゃ", 0),
		NewRomanEntry("byi", "びぃ", 0),
		NewRomanEntry("byu", "びゅ", 0),
		NewRomanEntry("bye", "びぇ", 0),
		NewRomanEntry("byo", "びょ", 0),
		NewRomanEntry("pya", "ぴゃ", 0),
		NewRomanEntry("pyi", "ぴぃ", 0),
		NewRomanEntry("pyu", "ぴゅ", 0),
		NewRomanEntry("pye", "ぴぇ", 0),
		NewRomanEntry("pyo", "ぴょ", 0),
		NewRomanEntry("fa", "ふぁ", 0),
		NewRomanEntry("fi", "ふぃ", 0),
		NewRomanEntry("fu", "ふ", 0),
		NewRomanEntry("fe", "ふぇ", 0),
		NewRomanEntry("fo", "ふぉ", 0),
		NewRomanEntry("fya", "ふゃ", 0),
		NewRomanEntry("fyu", "ふゅ", 0),
		NewRomanEntry("fyo", "ふょ", 0),
		NewRomanEntry("hwa", "ふぁ", 0),
		NewRomanEntry("hwi", "ふぃ", 0),
		NewRomanEntry("hwe", "ふぇ", 0),
		NewRomanEntry("hwo", "ふぉ", 0),
		NewRomanEntry("hwyu", "ふゅ", 0),
		NewRomanEntry("mya", "みゃ", 0),
		NewRomanEntry("myi", "みぃ", 0),
		NewRomanEntry("myu", "みゅ", 0),
		NewRomanEntry("mye", "みぇ", 0),
		NewRomanEntry("myo", "みょ", 0),
		NewRomanEntry("rya", "りゃ", 0),
		NewRomanEntry("ryi", "りぃ", 0),
		NewRomanEntry("ryu", "りゅ", 0),
		NewRomanEntry("rye", "りぇ", 0),
		NewRomanEntry("ryo", "りょ", 0),
		NewRomanEntry("n'", "ん", 0),
		NewRomanEntry("nn", "ん", 0),
		NewRomanEntry("n", "ん", 0),
		NewRomanEntry("xn", "ん", 0),
		NewRomanEntry("a", "あ", 0),
		NewRomanEntry("i", "い", 0),
		NewRomanEntry("u", "う", 0),
		NewRomanEntry("wu", "う", 0),
		NewRomanEntry("e", "え", 0),
		NewRomanEntry("o", "お", 0),
		NewRomanEntry("xa", "ぁ", 0),
		NewRomanEntry("xi", "ぃ", 0),
		NewRomanEntry("xu", "ぅ", 0),
		NewRomanEntry("xe", "ぇ", 0),
		NewRomanEntry("xo", "ぉ", 0),
		NewRomanEntry("la", "ぁ", 0),
		NewRomanEntry("li", "ぃ", 0),
		NewRomanEntry("lu", "ぅ", 0),
		NewRomanEntry("le", "ぇ", 0),
		NewRomanEntry("lo", "ぉ", 0),
		NewRomanEntry("lyi", "ぃ", 0),
		NewRomanEntry("xyi", "ぃ", 0),
		NewRomanEntry("lye", "ぇ", 0),
		NewRomanEntry("xye", "ぇ", 0),
		NewRomanEntry("ye", "いぇ", 0),
		NewRomanEntry("ka", "か", 0),
		NewRomanEntry("ki", "き", 0),
		NewRomanEntry("ku", "く", 0),
		NewRomanEntry("ke", "け", 0),
		NewRomanEntry("ko", "こ", 0),
		NewRomanEntry("xka", "ヵ", 0),
		NewRomanEntry("xke", "ヶ", 0),
		NewRomanEntry("lka", "ヵ", 0),
		NewRomanEntry("lke", "ヶ", 0),
		NewRomanEntry("ga", "が", 0),
		NewRomanEntry("gi", "ぎ", 0),
		NewRomanEntry("gu", "ぐ", 0),
		NewRomanEntry("ge", "げ", 0),
		NewRomanEntry("go", "ご", 0),
		NewRomanEntry("sa", "さ", 0),
		NewRomanEntry("si", "し", 0),
		NewRomanEntry("su", "す", 0),
		NewRomanEntry("se", "せ", 0),
		NewRomanEntry("so", "そ", 0),
		NewRomanEntry("ca", "か", 0),
		NewRomanEntry("ci", "し", 0),
		NewRomanEntry("cu", "く", 0),
		NewRomanEntry("ce", "せ", 0),
		NewRomanEntry("co", "こ", 0),
		NewRomanEntry("qa", "くぁ", 0),
		NewRomanEntry("qi", "くぃ", 0),
		NewRomanEntry("qu", "く", 0),
		NewRomanEntry("qe", "くぇ", 0),
		NewRomanEntry("qo", "くぉ", 0),
		NewRomanEntry("kwa", "くぁ", 0),
		NewRomanEntry("kwi", "くぃ", 0),
		NewRomanEntry("kwu", "くぅ", 0),
		NewRomanEntry("kwe", "くぇ", 0),
		NewRomanEntry("kwo", "くぉ", 0),
		NewRomanEntry("gwa", "ぐぁ", 0),
		NewRomanEntry("gwi", "ぐぃ", 0),
		NewRomanEntry("gwu", "ぐぅ", 0),
		NewRomanEntry("gwe", "ぐぇ", 0),
		NewRomanEntry("gwo", "ぐぉ", 0),
		NewRomanEntry("za", "ざ", 0),
		NewRomanEntry("zi", "じ", 0),
		NewRomanEntry("zu", "ず", 0),
		NewRomanEntry("ze", "ぜ", 0),
		NewRomanEntry("zo", "ぞ", 0),
		NewRomanEntry("ja", "じゃ", 0),
		NewRomanEntry("ji", "じ", 0),
		NewRomanEntry("ju", "じゅ", 0),
		NewRomanEntry("je", "じぇ", 0),
		NewRomanEntry("jo", "じょ", 0),
		NewRomanEntry("jya", "じゃ", 0),
		NewRomanEntry("jyi", "じぃ", 0),
		NewRomanEntry("jyu", "じゅ", 0),
		NewRomanEntry("jye", "じぇ", 0),
		NewRomanEntry("jyo", "じょ", 0),
		NewRomanEntry("ta", "た", 0),
		NewRomanEntry("ti", "ち", 0),
		NewRomanEntry("tu", "つ", 0),
		NewRomanEntry("tsu", "つ", 0),
		NewRomanEntry("te", "て", 0),
		NewRomanEntry("to", "と", 0),
		NewRomanEntry("da", "だ", 0),
		NewRomanEntry("di", "ぢ", 0),
		NewRomanEntry("du", "づ", 0),
		NewRomanEntry("de", "で", 0),
		NewRomanEntry("do", "ど", 0),
		NewRomanEntry("xtu", "っ", 0),
		NewRomanEntry("xtsu", "っ", 0),
		NewRomanEntry("ltu", "っ", 0),
		NewRomanEntry("ltsu", "っ", 0),
		NewRomanEntry("na", "な", 0),
		NewRomanEntry("ni", "に", 0),
		NewRomanEntry("nu", "ぬ", 0),
		NewRomanEntry("ne", "ね", 0),
		NewRomanEntry("no", "の", 0),
		NewRomanEntry("ha", "は", 0),
		NewRomanEntry("hi", "ひ", 0),
		NewRomanEntry("hu", "ふ", 0),
		NewRomanEntry("fu", "ふ", 0),
		NewRomanEntry("he", "へ", 0),
		NewRomanEntry("ho", "ほ", 0),
		NewRomanEntry("ba", "ば", 0),
		NewRomanEntry("bi", "び", 0),
		NewRomanEntry("bu", "ぶ", 0),
		NewRomanEntry("be", "べ", 0),
		NewRomanEntry("bo", "ぼ", 0),
		NewRomanEntry("pa", "ぱ", 0),
		NewRomanEntry("pi", "ぴ", 0),
		NewRomanEntry("pu", "ぷ", 0),
		NewRomanEntry("pe", "ぺ", 0),
		NewRomanEntry("po", "ぽ", 0),
		NewRomanEntry("ma", "ま", 0),
		NewRomanEntry("mi", "み", 0),
		NewRomanEntry("mu", "む", 0),
		NewRomanEntry("me", "め", 0),
		NewRomanEntry("mo", "も", 0),
		NewRomanEntry("xya", "ゃ", 0),
		NewRomanEntry("lya", "ゃ", 0),
		NewRomanEntry("ya", "や", 0),
		NewRomanEntry("wyi", "ゐ", 0),
		NewRomanEntry("xyu", "ゅ", 0),
		NewRomanEntry("lyu", "ゅ", 0),
		NewRomanEntry("yu", "ゆ", 0),
		NewRomanEntry("wye", "ゑ", 0),
		NewRomanEntry("xyo", "ょ", 0),
		NewRomanEntry("lyo", "ょ", 0),
		NewRomanEntry("yo", "よ", 0),
		NewRomanEntry("ra", "ら", 0),
		NewRomanEntry("ri", "り", 0),
		NewRomanEntry("ru", "る", 0),
		NewRomanEntry("re", "れ", 0),
		NewRomanEntry("ro", "ろ", 0),
		NewRomanEntry("xwa", "ゎ", 0),
		NewRomanEntry("lwa", "ゎ", 0),
		NewRomanEntry("wa", "わ", 0),
		NewRomanEntry("wi", "うぃ", 0),
		NewRomanEntry("we", "うぇ", 0),
		NewRomanEntry("wo", "を", 0),
		NewRomanEntry("wha", "うぁ", 0),
		NewRomanEntry("whi", "うぃ", 0),
		NewRomanEntry("whu", "う", 0),
		NewRomanEntry("whe", "うぇ", 0),
		NewRomanEntry("who", "うぉ", 0),
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].index < entries[j].index })
	var indexes = make([]uint32, len(entries))
	for i, _ := range indexes {
		indexes[i] = entries[i].index
	}
	return &RomajiProcessor{
		entries,
		indexes,
	}
}

func calculateIndex(roman []uint16) uint32 {
	return _calculateIndex(roman, 0, 4)
}

func _calculateIndex(roman []uint16, start int, end int) uint32 {
	var result uint32 = 0
	for i := 0; i < 4; i++ {
		index := i + start
		var c uint16 = 0
		if index < len(roman) && index < end {
			c = roman[index]
		}
		result = result | uint32(c)
		if i < 3 {
			result = result << 8
		}
	}
	return result
}

type RomajiPredictiveResult struct {
	Prefix   []uint16
	Suffixes [][]uint16
}

func (this *RomajiProcessor) RomajiToHiragana(romaji []uint16) []uint16 {
	if len(romaji) == 0 {
		return []uint16{}
	}
	var hiragana = []uint16{}
	var start = 0
	var end = 1
	for start < len(romaji) {
		var lastFound = -1
		var lower = 0
		var upper = len(this.indexes)
		for upper-lower > 1 && end <= len(romaji) {
			var lowerKey = _calculateIndex(romaji, start, end)
			lower = binarySearch(this.indexes, lower, upper, lowerKey)
			if lower >= 0 {
				lastFound = lower
			} else {
				lower = -lower - 1
			}
			var upperKey = lowerKey + (1 << (32 - 8*(end-start)))
			upper = binarySearch(this.indexes, lower, upper, upperKey)
			if upper < 0 {
				upper = -upper - 1
			}
			end = end + 1
		}
		if lastFound >= 0 {
			var entry = this.entries[lastFound]
			hiragana = append(hiragana, entry.hiragana...)
			start = start + len(entry.roman) - entry.remain
			end = start + 1
		} else {
			hiragana = append(hiragana, romaji[start])
			start = start + 1
			end = start + 1
		}
	}
	return hiragana
}

func (this *RomajiProcessor) findRomanEntryPredicatively(roman []uint16, offset int) []RomanEntry {
	var startIndex = 0
	var endIndex = len(this.indexes)
	for i := 0; i < 4; i++ {
		if len(roman) <= offset+i {
			break
		}
		var startKey = _calculateIndex(roman, offset, offset+i+1)
		startIndex = binarySearch(this.indexes, startIndex, endIndex, startKey)
		if startIndex >= 0 {
		} else {
			startIndex = -startIndex - 1
		}
		var endKey = startKey + (1 << (24 - 8*i))
		endIndex = binarySearch(this.indexes, startIndex, endIndex, endKey)
		if endIndex < 0 {
			endIndex = -endIndex - 1
		}
		if endIndex-startIndex == 1 {
			return []RomanEntry{*this.entries[startIndex]}
		}
	}
	var result = []RomanEntry{}
	for i := startIndex; i < endIndex; i++ {
		result = append(result, *this.entries[i])
	}
	return result
}

func (this *RomajiProcessor) RomajiToHiraganaPredictively(romaji []uint16) RomajiPredictiveResult {
	if len(romaji) == 0 {
		return RomajiPredictiveResult{
			Prefix:   []uint16{},
			Suffixes: [][]uint16{{}},
		}
	}
	var hiragana = []uint16{}
	var start = 0
	var end = 1
	for start < len(romaji) {
		var lastFound = -1
		var lower = 0
		var upper = len(this.indexes)
		for upper-lower > 1 && end <= len(romaji) {
			var lowerKey = _calculateIndex(romaji, start, end)
			lower = binarySearch(this.indexes, lower, upper, lowerKey)
			if lower >= 0 {
				lastFound = lower
			} else {
				lower = -lower - 1
			}
			var upperKey = lowerKey + (1 << (32 - 8*(end-start)))
			upper = binarySearch(this.indexes, lower, upper, upperKey)
			if upper < 0 {
				upper = -upper - 1
			}
			end++
		}
		if end > len(romaji) && upper-lower > 1 {
			var set = [][]uint16{}
			for i := lower; i < upper; i++ {
				var re = this.entries[i]
				if re.remain > 0 {
					var set2 = this.findRomanEntryPredicatively(romaji, end-1-re.remain)
					for _, re2 := range set2 {
						if re2.remain == 0 {
							set = append(set, append(re.hiragana, re2.hiragana...))
						}
					}
				} else {
					set = append(set, re.hiragana)
				}
			}
			return RomajiPredictiveResult{
				Prefix:   hiragana,
				Suffixes: set,
			}
		}
		if lastFound >= 0 {
			var entry = this.entries[lastFound]
			hiragana = append(hiragana, entry.hiragana...)
			start = start + len(entry.roman) - entry.remain
			end = start + 1
		} else {
			hiragana = append(hiragana, romaji[start])
			start++
			end = start + 1
		}
	}
	return RomajiPredictiveResult{
		Prefix:   hiragana,
		Suffixes: [][]uint16{{}},
	}
}

func binarySearch(a []uint32, fromIndex int, toIndex int, key uint32) int {
	var low = fromIndex
	var high = toIndex - 1
	for low <= high {
		var mid = (low + high) >> 1
		var midVal = a[mid]

		if midVal < key {
			low = mid + 1
		} else if midVal > key {
			high = mid - 1
		} else {
			return mid
		}
	}
	return -(low + 1)
}
