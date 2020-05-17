package migemo

var han2zen = map[rune]rune{
	'!':  '！',
	'"':  '”',
	'#':  '＃',
	'$':  '＄',
	'%':  '％',
	'&':  '＆',
	'\'': '’',
	'(':  '（',
	')':  '）',
	'*':  '＊',
	'+':  '＋',
	',':  '，',
	'-':  '－',
	'.':  '．',
	'/':  '／',
	'0':  '０',
	'1':  '１',
	'2':  '２',
	'3':  '３',
	'4':  '４',
	'5':  '５',
	'6':  '６',
	'7':  '７',
	'8':  '８',
	'9':  '９',
	':':  '：',
	';':  '；',
	'<':  '＜',
	'=':  '＝',
	'>':  '＞',
	'?':  '？',
	'@':  '＠',
	'A':  'Ａ',
	'B':  'Ｂ',
	'C':  'Ｃ',
	'D':  'Ｄ',
	'E':  'Ｅ',
	'F':  'Ｆ',
	'G':  'Ｇ',
	'H':  'Ｈ',
	'I':  'Ｉ',
	'J':  'Ｊ',
	'K':  'Ｋ',
	'L':  'Ｌ',
	'M':  'Ｍ',
	'N':  'Ｎ',
	'O':  'Ｏ',
	'P':  'Ｐ',
	'Q':  'Ｑ',
	'R':  'Ｒ',
	'S':  'Ｓ',
	'T':  'Ｔ',
	'U':  'Ｕ',
	'V':  'Ｖ',
	'W':  'Ｗ',
	'X':  'Ｘ',
	'Y':  'Ｙ',
	'Z':  'Ｚ',
	'[':  '［',
	'\\': '￥',
	']':  '］',
	'^':  '＾',
	'_':  '＿',
	'`':  '‘',
	'a':  'ａ',
	'b':  'ｂ',
	'c':  'ｃ',
	'd':  'ｄ',
	'e':  'ｅ',
	'f':  'ｆ',
	'g':  'ｇ',
	'h':  'ｈ',
	'i':  'ｉ',
	'j':  'ｊ',
	'k':  'ｋ',
	'l':  'ｌ',
	'm':  'ｍ',
	'n':  'ｎ',
	'o':  'ｏ',
	'p':  'ｐ',
	'q':  'ｑ',
	'r':  'ｒ',
	's':  'ｓ',
	't':  'ｔ',
	'u':  'ｕ',
	'v':  'ｖ',
	'w':  'ｗ',
	'x':  'ｘ',
	'y':  'ｙ',
	'z':  'ｚ',
	'{':  '｛',
	'|':  '｜',
	'}':  '｝',
	'~':  '～',
	'｡':  '。',
	'｢':  '「',
	'｣':  '」',
	'､':  '、',
	'･':  '・',
	'ｦ':  'ヲ',
	'ｧ':  'ァ',
	'ｨ':  'ィ',
	'ｩ':  'ゥ',
	'ｪ':  'ェ',
	'ｫ':  'ォ',
	'ｬ':  'ャ',
	'ｭ':  'ュ',
	'ｮ':  'ョ',
	'ｯ':  'ッ',
	'ｰ':  'ー',
	'ｱ':  'ア',
	'ｲ':  'イ',
	'ｳ':  'ウ',
	'ｴ':  'エ',
	'ｵ':  'オ',
	'ｶ':  'カ',
	'ｷ':  'キ',
	'ｸ':  'ク',
	'ｹ':  'ケ',
	'ｺ':  'コ',
	'ｻ':  'サ',
	'ｼ':  'シ',
	'ｽ':  'ス',
	'ｾ':  'セ',
	'ｿ':  'ソ',
	'ﾀ':  'タ',
	'ﾁ':  'チ',
	'ﾂ':  'ツ',
	'ﾃ':  'テ',
	'ﾄ':  'ト',
	'ﾅ':  'ナ',
	'ﾆ':  'ニ',
	'ﾇ':  'ヌ',
	'ﾈ':  'ネ',
	'ﾉ':  'ノ',
	'ﾊ':  'ハ',
	'ﾋ':  'ヒ',
	'ﾌ':  'フ',
	'ﾍ':  'ヘ',
	'ﾎ':  'ホ',
	'ﾏ':  'マ',
	'ﾐ':  'ミ',
	'ﾑ':  'ム',
	'ﾒ':  'メ',
	'ﾓ':  'モ',
	'ﾔ':  'ヤ',
	'ﾕ':  'ユ',
	'ﾖ':  'ヨ',
	'ﾗ':  'ラ',
	'ﾘ':  'リ',
	'ﾙ':  'ル',
	'ﾚ':  'レ',
	'ﾛ':  'ロ',
	'ﾜ':  'ワ',
	'ﾝ':  'ン',
	'ﾞ':  '゛',
	'ﾟ':  '゜',
}

var zen2han = map[rune]string{
	'！': "!",
	'”': "\"",
	'＃': "#",
	'＄': "$",
	'％': "%",
	'＆': "&",
	'’': "'",
	'（': "(",
	'）': ")",
	'＊': "*",
	'＋': "+",
	'，': ",",
	'－': "-",
	'．': ".",
	'／': "/",
	'０': "0",
	'１': "1",
	'２': "2",
	'３': "3",
	'４': "4",
	'５': "5",
	'６': "6",
	'７': "7",
	'８': "8",
	'９': "9",
	'：': ":",
	'；': ";",
	'＜': "<",
	'＝': "=",
	'＞': ">",
	'？': "?",
	'＠': "@",
	'Ａ': "A",
	'Ｂ': "B",
	'Ｃ': "C",
	'Ｄ': "D",
	'Ｅ': "E",
	'Ｆ': "F",
	'Ｇ': "G",
	'Ｈ': "H",
	'Ｉ': "I",
	'Ｊ': "J",
	'Ｋ': "K",
	'Ｌ': "L",
	'Ｍ': "M",
	'Ｎ': "N",
	'Ｏ': "O",
	'Ｐ': "P",
	'Ｑ': "Q",
	'Ｒ': "R",
	'Ｓ': "S",
	'Ｔ': "T",
	'Ｕ': "U",
	'Ｖ': "V",
	'Ｗ': "W",
	'Ｘ': "X",
	'Ｙ': "Y",
	'Ｚ': "Z",
	'［': "[",
	'￥': "\\",
	'］': "]",
	'＾': "^",
	'＿': "_",
	'‘': "`",
	'ａ': "a",
	'ｂ': "b",
	'ｃ': "c",
	'ｄ': "d",
	'ｅ': "e",
	'ｆ': "f",
	'ｇ': "g",
	'ｈ': "h",
	'ｉ': "i",
	'ｊ': "j",
	'ｋ': "k",
	'ｌ': "l",
	'ｍ': "m",
	'ｎ': "n",
	'ｏ': "o",
	'ｐ': "p",
	'ｑ': "q",
	'ｒ': "r",
	'ｓ': "s",
	'ｔ': "t",
	'ｕ': "u",
	'ｖ': "v",
	'ｗ': "w",
	'ｘ': "x",
	'ｙ': "y",
	'ｚ': "z",
	'｛': "{",
	'｜': "|",
	'｝': "}",
	'～': "~",
	'。': "｡",
	'「': "｢",
	'」': "｣",
	'、': "､",
	'・': "･",
	'ヲ': "ｦ",
	'ァ': "ｧ",
	'ィ': "ｨ",
	'ゥ': "ｩ",
	'ェ': "ｪ",
	'ォ': "ｫ",
	'ャ': "ｬ",
	'ュ': "ｭ",
	'ョ': "ｮ",
	'ッ': "ｯ",
	'ー': "ｰ",
	'ア': "ｱ",
	'イ': "ｲ",
	'ウ': "ｳ",
	'エ': "ｴ",
	'オ': "ｵ",
	'カ': "ｶ",
	'キ': "ｷ",
	'ク': "ｸ",
	'ケ': "ｹ",
	'コ': "ｺ",
	'サ': "ｻ",
	'シ': "ｼ",
	'ス': "ｽ",
	'セ': "ｾ",
	'ソ': "ｿ",
	'タ': "ﾀ",
	'チ': "ﾁ",
	'ツ': "ﾂ",
	'テ': "ﾃ",
	'ト': "ﾄ",
	'ナ': "ﾅ",
	'ニ': "ﾆ",
	'ヌ': "ﾇ",
	'ネ': "ﾈ",
	'ノ': "ﾉ",
	'ハ': "ﾊ",
	'ヒ': "ﾋ",
	'フ': "ﾌ",
	'ヘ': "ﾍ",
	'ホ': "ﾎ",
	'マ': "ﾏ",
	'ミ': "ﾐ",
	'ム': "ﾑ",
	'メ': "ﾒ",
	'モ': "ﾓ",
	'ヤ': "ﾔ",
	'ユ': "ﾕ",
	'ヨ': "ﾖ",
	'ラ': "ﾗ",
	'リ': "ﾘ",
	'ル': "ﾙ",
	'レ': "ﾚ",
	'ロ': "ﾛ",
	'ワ': "ﾜ",
	'ン': "ﾝ",
	'゛': "ﾞ",
	'゜': "ﾟ",
	'ヴ': "ｳﾞ",
	'ガ': "ｶﾞ",
	'ギ': "ｷﾞ",
	'グ': "ｸﾞ",
	'ゲ': "ｹﾞ",
	'ゴ': "ｺﾞ",
	'ザ': "ｻﾞ",
	'ジ': "ｼﾞ",
	'ズ': "ｽﾞ",
	'ゼ': "ｾﾞ",
	'ゾ': "ｿﾞ",
	'ダ': "ﾀﾞ",
	'ヂ': "ﾁﾞ",
	'ヅ': "ﾂﾞ",
	'デ': "ﾃﾞ",
	'ド': "ﾄﾞ",
	'バ': "ﾊﾞ",
	'ビ': "ﾋﾞ",
	'ブ': "ﾌﾞ",
	'ベ': "ﾍﾞ",
	'ボ': "ﾎﾞ",
	'パ': "ﾊﾟ",
	'ピ': "ﾋﾟ",
	'プ': "ﾌﾟ",
	'ペ': "ﾍﾟ",
	'ポ': "ﾎﾟ",
}

// ConvertHan2Zen は、半角から全角へ文字列を変更する
func ConvertHan2Zen(source string) string {
	var sb = []rune{}
	for _, c := range []rune(source) {
		a, ok := han2zen[c]
		if !ok {
			sb = append(sb, c)
		} else {
			sb = append(sb, a)
		}
	}
	return string(sb)
}

// ConvertZen2Han は、全角から半角へ文字列を変更する
func ConvertZen2Han(source string) string {
	var sb = []rune{}
	for _, c := range []rune(source) {
		a, ok := zen2han[c]
		if !ok {
			sb = append(sb, c)
		} else {
			sb = append(sb, []rune(a)...)
		}
	}
	return string(sb)
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
