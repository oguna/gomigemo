package migemo

import (
	"sort"
	"strings"
)

type romanEntry struct {
	roman    string
	hiragana string
	remain   int
	index    uint32
}

// RomajiProcessor は、ローマ字を処理する構造体
type RomajiProcessor struct {
	entries []*romanEntry
	indexes []uint32
}

func newRomanEntry(roman string, hiragana string, remain int) *romanEntry {
	return &romanEntry{
		roman:    roman,
		hiragana: hiragana,
		remain:   remain,
		index:    calculateIndex(roman),
	}
}

// NewRomajiProcessor は、RomajiProcessorを初期化する
func NewRomajiProcessor() *RomajiProcessor {
	var entries = []*romanEntry{
		newRomanEntry("-", "ー", 0),
		newRomanEntry("~", "〜", 0),
		newRomanEntry(".", "。", 0),
		newRomanEntry(",", "、", 0),
		newRomanEntry("z/", "・", 0),
		newRomanEntry("z.", "…", 0),
		newRomanEntry("z,", "‥", 0),
		newRomanEntry("zh", "←", 0),
		newRomanEntry("zj", "↓", 0),
		newRomanEntry("zk", "↑", 0),
		newRomanEntry("zl", "→", 0),
		newRomanEntry("z-", "〜", 0),
		newRomanEntry("z[", "『", 0),
		newRomanEntry("z]", "』", 0),
		newRomanEntry("[", "「", 0),
		newRomanEntry("]", "」", 0),
		newRomanEntry("va", "ゔぁ", 0),
		newRomanEntry("vi", "ゔぃ", 0),
		newRomanEntry("vu", "ゔ", 0),
		newRomanEntry("ve", "ゔぇ", 0),
		newRomanEntry("vo", "ゔぉ", 0),
		newRomanEntry("vya", "ゔゃ", 0),
		newRomanEntry("vyi", "ゔぃ", 0),
		newRomanEntry("vyu", "ゔゅ", 0),
		newRomanEntry("vye", "ゔぇ", 0),
		newRomanEntry("vyo", "ゔょ", 0),
		newRomanEntry("qq", "っ", 1),
		newRomanEntry("vv", "っ", 1),
		newRomanEntry("ll", "っ", 1),
		newRomanEntry("xx", "っ", 1),
		newRomanEntry("kk", "っ", 1),
		newRomanEntry("gg", "っ", 1),
		newRomanEntry("ss", "っ", 1),
		newRomanEntry("zz", "っ", 1),
		newRomanEntry("jj", "っ", 1),
		newRomanEntry("tt", "っ", 1),
		newRomanEntry("dd", "っ", 1),
		newRomanEntry("hh", "っ", 1),
		newRomanEntry("ff", "っ", 1),
		newRomanEntry("bb", "っ", 1),
		newRomanEntry("pp", "っ", 1),
		newRomanEntry("mm", "っ", 1),
		newRomanEntry("yy", "っ", 1),
		newRomanEntry("rr", "っ", 1),
		newRomanEntry("ww", "っ", 1),
		newRomanEntry("www", "w", 2),
		newRomanEntry("cc", "っ", 1),
		newRomanEntry("kya", "きゃ", 0),
		newRomanEntry("kyi", "きぃ", 0),
		newRomanEntry("kyu", "きゅ", 0),
		newRomanEntry("kye", "きぇ", 0),
		newRomanEntry("kyo", "きょ", 0),
		newRomanEntry("gya", "ぎゃ", 0),
		newRomanEntry("gyi", "ぎぃ", 0),
		newRomanEntry("gyu", "ぎゅ", 0),
		newRomanEntry("gye", "ぎぇ", 0),
		newRomanEntry("gyo", "ぎょ", 0),
		newRomanEntry("sya", "しゃ", 0),
		newRomanEntry("syi", "しぃ", 0),
		newRomanEntry("syu", "しゅ", 0),
		newRomanEntry("sye", "しぇ", 0),
		newRomanEntry("syo", "しょ", 0),
		newRomanEntry("sha", "しゃ", 0),
		newRomanEntry("shi", "し", 0),
		newRomanEntry("shu", "しゅ", 0),
		newRomanEntry("she", "しぇ", 0),
		newRomanEntry("sho", "しょ", 0),
		newRomanEntry("zya", "じゃ", 0),
		newRomanEntry("zyi", "じぃ", 0),
		newRomanEntry("zyu", "じゅ", 0),
		newRomanEntry("zye", "じぇ", 0),
		newRomanEntry("zyo", "じょ", 0),
		newRomanEntry("tya", "ちゃ", 0),
		newRomanEntry("tyi", "ちぃ", 0),
		newRomanEntry("tyu", "ちゅ", 0),
		newRomanEntry("tye", "ちぇ", 0),
		newRomanEntry("tyo", "ちょ", 0),
		newRomanEntry("cha", "ちゃ", 0),
		newRomanEntry("chi", "ち", 0),
		newRomanEntry("chu", "ちゅ", 0),
		newRomanEntry("che", "ちぇ", 0),
		newRomanEntry("cho", "ちょ", 0),
		newRomanEntry("cya", "ちゃ", 0),
		newRomanEntry("cyi", "ちぃ", 0),
		newRomanEntry("cyu", "ちゅ", 0),
		newRomanEntry("cye", "ちぇ", 0),
		newRomanEntry("cyo", "ちょ", 0),
		newRomanEntry("dya", "ぢゃ", 0),
		newRomanEntry("dyi", "ぢぃ", 0),
		newRomanEntry("dyu", "ぢゅ", 0),
		newRomanEntry("dye", "ぢぇ", 0),
		newRomanEntry("dyo", "ぢょ", 0),
		newRomanEntry("tsa", "つぁ", 0),
		newRomanEntry("tsi", "つぃ", 0),
		newRomanEntry("tse", "つぇ", 0),
		newRomanEntry("tso", "つぉ", 0),
		newRomanEntry("tha", "てゃ", 0),
		newRomanEntry("thi", "てぃ", 0),
		newRomanEntry("t'i", "てぃ", 0),
		newRomanEntry("thu", "てゅ", 0),
		newRomanEntry("the", "てぇ", 0),
		newRomanEntry("tho", "てょ", 0),
		newRomanEntry("t'yu", "てゅ", 0),
		newRomanEntry("dha", "でゃ", 0),
		newRomanEntry("dhi", "でぃ", 0),
		newRomanEntry("d'i", "でぃ", 0),
		newRomanEntry("dhu", "でゅ", 0),
		newRomanEntry("dhe", "でぇ", 0),
		newRomanEntry("dho", "でょ", 0),
		newRomanEntry("d'yu", "でゅ", 0),
		newRomanEntry("twa", "とぁ", 0),
		newRomanEntry("twi", "とぃ", 0),
		newRomanEntry("twu", "とぅ", 0),
		newRomanEntry("twe", "とぇ", 0),
		newRomanEntry("two", "とぉ", 0),
		newRomanEntry("t'u", "とぅ", 0),
		newRomanEntry("dwa", "どぁ", 0),
		newRomanEntry("dwi", "どぃ", 0),
		newRomanEntry("dwu", "どぅ", 0),
		newRomanEntry("dwe", "どぇ", 0),
		newRomanEntry("dwo", "どぉ", 0),
		newRomanEntry("d'u", "どぅ", 0),
		newRomanEntry("nya", "にゃ", 0),
		newRomanEntry("nyi", "にぃ", 0),
		newRomanEntry("nyu", "にゅ", 0),
		newRomanEntry("nye", "にぇ", 0),
		newRomanEntry("nyo", "にょ", 0),
		newRomanEntry("hya", "ひゃ", 0),
		newRomanEntry("hyi", "ひぃ", 0),
		newRomanEntry("hyu", "ひゅ", 0),
		newRomanEntry("hye", "ひぇ", 0),
		newRomanEntry("hyo", "ひょ", 0),
		newRomanEntry("bya", "びゃ", 0),
		newRomanEntry("byi", "びぃ", 0),
		newRomanEntry("byu", "びゅ", 0),
		newRomanEntry("bye", "びぇ", 0),
		newRomanEntry("byo", "びょ", 0),
		newRomanEntry("pya", "ぴゃ", 0),
		newRomanEntry("pyi", "ぴぃ", 0),
		newRomanEntry("pyu", "ぴゅ", 0),
		newRomanEntry("pye", "ぴぇ", 0),
		newRomanEntry("pyo", "ぴょ", 0),
		newRomanEntry("fa", "ふぁ", 0),
		newRomanEntry("fi", "ふぃ", 0),
		newRomanEntry("fu", "ふ", 0),
		newRomanEntry("fe", "ふぇ", 0),
		newRomanEntry("fo", "ふぉ", 0),
		newRomanEntry("fya", "ふゃ", 0),
		newRomanEntry("fyu", "ふゅ", 0),
		newRomanEntry("fyo", "ふょ", 0),
		newRomanEntry("hwa", "ふぁ", 0),
		newRomanEntry("hwi", "ふぃ", 0),
		newRomanEntry("hwe", "ふぇ", 0),
		newRomanEntry("hwo", "ふぉ", 0),
		newRomanEntry("hwyu", "ふゅ", 0),
		newRomanEntry("mya", "みゃ", 0),
		newRomanEntry("myi", "みぃ", 0),
		newRomanEntry("myu", "みゅ", 0),
		newRomanEntry("mye", "みぇ", 0),
		newRomanEntry("myo", "みょ", 0),
		newRomanEntry("rya", "りゃ", 0),
		newRomanEntry("ryi", "りぃ", 0),
		newRomanEntry("ryu", "りゅ", 0),
		newRomanEntry("rye", "りぇ", 0),
		newRomanEntry("ryo", "りょ", 0),
		newRomanEntry("n'", "ん", 0),
		newRomanEntry("nn", "ん", 0),
		newRomanEntry("n", "ん", 0),
		newRomanEntry("xn", "ん", 0),
		newRomanEntry("a", "あ", 0),
		newRomanEntry("i", "い", 0),
		newRomanEntry("u", "う", 0),
		newRomanEntry("wu", "う", 0),
		newRomanEntry("e", "え", 0),
		newRomanEntry("o", "お", 0),
		newRomanEntry("xa", "ぁ", 0),
		newRomanEntry("xi", "ぃ", 0),
		newRomanEntry("xu", "ぅ", 0),
		newRomanEntry("xe", "ぇ", 0),
		newRomanEntry("xo", "ぉ", 0),
		newRomanEntry("la", "ぁ", 0),
		newRomanEntry("li", "ぃ", 0),
		newRomanEntry("lu", "ぅ", 0),
		newRomanEntry("le", "ぇ", 0),
		newRomanEntry("lo", "ぉ", 0),
		newRomanEntry("lyi", "ぃ", 0),
		newRomanEntry("xyi", "ぃ", 0),
		newRomanEntry("lye", "ぇ", 0),
		newRomanEntry("xye", "ぇ", 0),
		newRomanEntry("ye", "いぇ", 0),
		newRomanEntry("ka", "か", 0),
		newRomanEntry("ki", "き", 0),
		newRomanEntry("ku", "く", 0),
		newRomanEntry("ke", "け", 0),
		newRomanEntry("ko", "こ", 0),
		newRomanEntry("xka", "ヵ", 0),
		newRomanEntry("xke", "ヶ", 0),
		newRomanEntry("lka", "ヵ", 0),
		newRomanEntry("lke", "ヶ", 0),
		newRomanEntry("ga", "が", 0),
		newRomanEntry("gi", "ぎ", 0),
		newRomanEntry("gu", "ぐ", 0),
		newRomanEntry("ge", "げ", 0),
		newRomanEntry("go", "ご", 0),
		newRomanEntry("sa", "さ", 0),
		newRomanEntry("si", "し", 0),
		newRomanEntry("su", "す", 0),
		newRomanEntry("se", "せ", 0),
		newRomanEntry("so", "そ", 0),
		newRomanEntry("ca", "か", 0),
		newRomanEntry("ci", "し", 0),
		newRomanEntry("cu", "く", 0),
		newRomanEntry("ce", "せ", 0),
		newRomanEntry("co", "こ", 0),
		newRomanEntry("qa", "くぁ", 0),
		newRomanEntry("qi", "くぃ", 0),
		newRomanEntry("qu", "く", 0),
		newRomanEntry("qe", "くぇ", 0),
		newRomanEntry("qo", "くぉ", 0),
		newRomanEntry("kwa", "くぁ", 0),
		newRomanEntry("kwi", "くぃ", 0),
		newRomanEntry("kwu", "くぅ", 0),
		newRomanEntry("kwe", "くぇ", 0),
		newRomanEntry("kwo", "くぉ", 0),
		newRomanEntry("gwa", "ぐぁ", 0),
		newRomanEntry("gwi", "ぐぃ", 0),
		newRomanEntry("gwu", "ぐぅ", 0),
		newRomanEntry("gwe", "ぐぇ", 0),
		newRomanEntry("gwo", "ぐぉ", 0),
		newRomanEntry("za", "ざ", 0),
		newRomanEntry("zi", "じ", 0),
		newRomanEntry("zu", "ず", 0),
		newRomanEntry("ze", "ぜ", 0),
		newRomanEntry("zo", "ぞ", 0),
		newRomanEntry("ja", "じゃ", 0),
		newRomanEntry("ji", "じ", 0),
		newRomanEntry("ju", "じゅ", 0),
		newRomanEntry("je", "じぇ", 0),
		newRomanEntry("jo", "じょ", 0),
		newRomanEntry("jya", "じゃ", 0),
		newRomanEntry("jyi", "じぃ", 0),
		newRomanEntry("jyu", "じゅ", 0),
		newRomanEntry("jye", "じぇ", 0),
		newRomanEntry("jyo", "じょ", 0),
		newRomanEntry("ta", "た", 0),
		newRomanEntry("ti", "ち", 0),
		newRomanEntry("tu", "つ", 0),
		newRomanEntry("tsu", "つ", 0),
		newRomanEntry("te", "て", 0),
		newRomanEntry("to", "と", 0),
		newRomanEntry("da", "だ", 0),
		newRomanEntry("di", "ぢ", 0),
		newRomanEntry("du", "づ", 0),
		newRomanEntry("de", "で", 0),
		newRomanEntry("do", "ど", 0),
		newRomanEntry("xtu", "っ", 0),
		newRomanEntry("xtsu", "っ", 0),
		newRomanEntry("ltu", "っ", 0),
		newRomanEntry("ltsu", "っ", 0),
		newRomanEntry("na", "な", 0),
		newRomanEntry("ni", "に", 0),
		newRomanEntry("nu", "ぬ", 0),
		newRomanEntry("ne", "ね", 0),
		newRomanEntry("no", "の", 0),
		newRomanEntry("ha", "は", 0),
		newRomanEntry("hi", "ひ", 0),
		newRomanEntry("hu", "ふ", 0),
		newRomanEntry("fu", "ふ", 0),
		newRomanEntry("he", "へ", 0),
		newRomanEntry("ho", "ほ", 0),
		newRomanEntry("ba", "ば", 0),
		newRomanEntry("bi", "び", 0),
		newRomanEntry("bu", "ぶ", 0),
		newRomanEntry("be", "べ", 0),
		newRomanEntry("bo", "ぼ", 0),
		newRomanEntry("pa", "ぱ", 0),
		newRomanEntry("pi", "ぴ", 0),
		newRomanEntry("pu", "ぷ", 0),
		newRomanEntry("pe", "ぺ", 0),
		newRomanEntry("po", "ぽ", 0),
		newRomanEntry("ma", "ま", 0),
		newRomanEntry("mi", "み", 0),
		newRomanEntry("mu", "む", 0),
		newRomanEntry("me", "め", 0),
		newRomanEntry("mo", "も", 0),
		newRomanEntry("xya", "ゃ", 0),
		newRomanEntry("lya", "ゃ", 0),
		newRomanEntry("ya", "や", 0),
		newRomanEntry("wyi", "ゐ", 0),
		newRomanEntry("xyu", "ゅ", 0),
		newRomanEntry("lyu", "ゅ", 0),
		newRomanEntry("yu", "ゆ", 0),
		newRomanEntry("wye", "ゑ", 0),
		newRomanEntry("xyo", "ょ", 0),
		newRomanEntry("lyo", "ょ", 0),
		newRomanEntry("yo", "よ", 0),
		newRomanEntry("ra", "ら", 0),
		newRomanEntry("ri", "り", 0),
		newRomanEntry("ru", "る", 0),
		newRomanEntry("re", "れ", 0),
		newRomanEntry("ro", "ろ", 0),
		newRomanEntry("xwa", "ゎ", 0),
		newRomanEntry("lwa", "ゎ", 0),
		newRomanEntry("wa", "わ", 0),
		newRomanEntry("wi", "うぃ", 0),
		newRomanEntry("we", "うぇ", 0),
		newRomanEntry("wo", "を", 0),
		newRomanEntry("wha", "うぁ", 0),
		newRomanEntry("whi", "うぃ", 0),
		newRomanEntry("whu", "う", 0),
		newRomanEntry("whe", "うぇ", 0),
		newRomanEntry("who", "うぉ", 0),
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].index < entries[j].index })
	var indexes = make([]uint32, len(entries))
	for i := range indexes {
		indexes[i] = entries[i].index
	}
	return &RomajiProcessor{
		entries,
		indexes,
	}
}

func calculateIndex(roman string) uint32 {
	return _calculateIndex(roman, 0, 4)
}

func _calculateIndex(roman string, start int, end int) uint32 {
	var result uint32 = 0
	for i := 0; i < 4; i++ {
		index := i + start
		var c uint8 = 0
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

// RomajiPredictiveResult は、RomajiProcessorの結果を返す構造体
type RomajiPredictiveResult struct {
	Prefix   string
	Suffixes []string
}

// RomajiToHiragana は、ローマ字からひらがなに変換する
func (processor *RomajiProcessor) RomajiToHiragana(romaji string) string {
	if len(romaji) == 0 {
		return ""
	}
	var hiragana strings.Builder
	var start = 0
	var end = 1
	for start < len(romaji) {
		var lastFound = -1
		var lower = 0
		var upper = len(processor.indexes)
		for upper-lower > 1 && end <= len(romaji) {
			var lowerKey = _calculateIndex(romaji, start, end)
			lower = binarySearch(processor.indexes, lower, upper, lowerKey)
			if lower >= 0 {
				lastFound = lower
			} else {
				lower = -lower - 1
			}
			var upperKey = lowerKey + (1 << (32 - 8*(end-start)))
			upper = binarySearch(processor.indexes, lower, upper, upperKey)
			if upper < 0 {
				upper = -upper - 1
			}
			end = end + 1
		}
		if lastFound >= 0 {
			var entry = processor.entries[lastFound]
			hiragana.WriteString(entry.hiragana)
			start = start + len(entry.roman) - entry.remain
			end = start + 1
		} else {
			hiragana.WriteByte(romaji[start])
			start = start + 1
			end = start + 1
		}
	}
	return hiragana.String()
}

func (processor *RomajiProcessor) findRomanEntryPredicatively(roman string, offset int) []romanEntry {
	var startIndex = 0
	var endIndex = len(processor.indexes)
	for i := 0; i < 4; i++ {
		if len(roman) <= offset+i {
			break
		}
		var startKey = _calculateIndex(roman, offset, offset+i+1)
		startIndex = binarySearch(processor.indexes, startIndex, endIndex, startKey)
		if startIndex >= 0 {
		} else {
			startIndex = -startIndex - 1
		}
		var endKey = startKey + (1 << (24 - 8*i))
		endIndex = binarySearch(processor.indexes, startIndex, endIndex, endKey)
		if endIndex < 0 {
			endIndex = -endIndex - 1
		}
		if endIndex-startIndex == 1 {
			return []romanEntry{*processor.entries[startIndex]}
		}
	}
	var result = []romanEntry{}
	for i := startIndex; i < endIndex; i++ {
		result = append(result, *processor.entries[i])
	}
	return result
}

// RomajiToHiraganaPredictively は、ローマ字からひらがなに変換する。末尾の未確定部分は予測した結果を返す
func (processor *RomajiProcessor) RomajiToHiraganaPredictively(romaji string) RomajiPredictiveResult {
	if len(romaji) == 0 {
		return RomajiPredictiveResult{
			Prefix:   "",
			Suffixes: []string{""},
		}
	}
	var hiragana strings.Builder
	var start = 0
	var end = 1
	for start < len(romaji) {
		var lastFound = -1
		var lower = 0
		var upper = len(processor.indexes)
		for upper-lower > 1 && end <= len(romaji) {
			var lowerKey = _calculateIndex(romaji, start, end)
			lower = binarySearch(processor.indexes, lower, upper, lowerKey)
			if lower >= 0 {
				lastFound = lower
			} else {
				lower = -lower - 1
			}
			var upperKey = lowerKey + (1 << (32 - 8*(end-start)))
			upper = binarySearch(processor.indexes, lower, upper, upperKey)
			if upper < 0 {
				upper = -upper - 1
			}
			end++
		}
		if end > len(romaji) && upper-lower > 1 {
			var set = make(map[string]struct{})
			for i := lower; i < upper; i++ {
				var re = processor.entries[i]
				if re.remain > 0 {
					var set2 = processor.findRomanEntryPredicatively(romaji, end-1-re.remain)
					for _, re2 := range set2 {
						if re2.remain == 0 {
							set[re.hiragana+re2.hiragana] = struct{}{}
						}
					}
				} else {
					set[re.hiragana] = struct{}{}
				}
			}
			array := make([]string, 0, len(set))
			for k := range set {
				array = append(array, k)
			}
			return RomajiPredictiveResult{
				Prefix:   hiragana.String(),
				Suffixes: array,
			}
		}
		if lastFound >= 0 {
			var entry = processor.entries[lastFound]
			hiragana.WriteString(entry.hiragana)
			start = start + len(entry.roman) - entry.remain
			end = start + 1
		} else {
			hiragana.WriteByte(romaji[start])
			start++
			end = start + 1
		}
	}
	return RomajiPredictiveResult{
		Prefix:   hiragana.String(),
		Suffixes: []string{""},
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
