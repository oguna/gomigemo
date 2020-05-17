package migemo

import (
	"bufio"
	"io"
	"sort"
	"strings"
	"unicode/utf16"
)

// BuildDictionaryFromMigemoDictFile は、ファイルからCompactDictionaryを読み込む
func BuildDictionaryFromMigemoDictFile(fp io.Reader) *CompactDictionary {
	scanner := bufio.NewScanner(fp)
	dict := make(map[string][]string)
	keys := make([]string, 0, 1024)
	values := make(map[string]struct{}, 1024)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") || len(line) == 0 {
			continue
		}
		columns := strings.Split(line, "\t")
		key := columns[0]
		var skip = false
		for _, c := range utf16.Encode([]rune(key)) {
			if encode(c) == 0 {
				println("skip this word: ", key)
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		keys = append(keys, key)
		for _, w := range columns[1:] {
			values[w] = struct{}{}
		}
		dict[key] = columns[1:]
	}

	// build key trie
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	keysUtf16 := make([][]uint16, len(keys))
	for i := 0; i < len(keys); i++ {
		keysUtf16[i] = utf16.Encode([]rune(keys[i]))
	}
	keyTrie, _, keyErr := BuildLoudsTrie(keysUtf16)
	if keyErr != nil {
		panic(keyErr)
	}

	// build value trie
	valuesUtf16 := make([][]uint16, 0, len(dict))
	for k := range values {
		valuesUtf16 = append(valuesUtf16, utf16.Encode([]rune(k)))
	}
	sort.Slice(valuesUtf16, func(i, j int) bool { return CompareUtf16String(valuesUtf16[i], valuesUtf16[j]) < 0 })
	valueTrie, _, valueTrieErr := BuildLoudsTrie(valuesUtf16)
	if valueTrieErr != nil {
		panic(valueTrieErr)
	}

	// build mapping from key trie to value trie
	mappingCount := 0
	for _, v := range dict {
		mappingCount += len(v)
	}
	mapping := make([]uint32, mappingCount)
	mappingIndex := 0
	mappingBitList := NewBitList()
	key := make([]uint16, 0, 16)
	for i := 1; i <= keyTrie.Size(); i++ {
		key = key[:0]
		keyTrie.ReverseLookup(uint32(i), &key)
		mappingBitList.Add(false)
		values, ok := dict[string(utf16.Decode(key))]
		if ok {
			for j := 0; j < len(values); j++ {
				mappingBitList.Add(true)
				mapping[mappingIndex] = uint32(valueTrie.Lookup(utf16.Encode([]rune(values[j]))))
				mappingIndex++
			}
		}
	}
	mappingBitVector := NewBitVector(mappingBitList.Words, uint32(mappingBitList.Size))

	return &CompactDictionary{
		keyTrie:           keyTrie,
		valueTrie:         valueTrie,
		mapping:           mapping,
		mappingBitVector:  mappingBitVector,
		hasMappingBitList: createHasMappingBitList(mappingBitVector),
	}
}

// ExtractTail は、文字列の配列から分岐のない末尾(TAIL)を抽出する
func ExtractTail(words []string) []uint32 {
	tails := make([]uint32, len(words))
	for i := 0; i < len(words); i++ {
		prevWord := []rune{}
		if i != 0 {
			prevWord = []rune(words[i-1])
		}
		currentWord := []rune(words[i])
		nextWord := []rune{}
		if i != len(words)-1 {
			nextWord = []rune(words[i+1])
		}
		cursor := 0
		for true {
			prevChar := rune(0)
			currentChar := rune(0)
			nextChar := rune(0)
			if cursor < len(prevWord) {
				prevChar = prevWord[cursor]
			}
			if cursor < len(currentWord) {
				currentChar = currentWord[cursor]
			}
			if cursor < len(nextWord) {
				nextChar = nextWord[cursor]
			}
			if prevChar == 0 && currentChar == 0 && nextChar == 0 {
				break
			}
			if prevChar != currentChar && currentChar != nextChar {
				break
			}
			cursor++
		}
		if cursor+1 < len(currentWord) {
			tails[i] = uint32(len(currentWord)) - uint32(cursor) - 1
		} else {
			tails[i] = 0
		}
	}
	return tails
}
