package migemo

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"unicode/utf16"
)

func BuildDictionaryFromMigemoDictFile(fp *os.File) *CompactDictionary {
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
	for i := 1; i <= keyTrie.Size(); i++ {
		key := keyTrie.GetKey(uint32(i))
		mappingBitList.Add(false)
		values, ok := dict[string(utf16.Decode(key))]
		if ok {
			for j := 0; j < len(values); j++ {
				mappingBitList.Add(true)
				mapping[mappingIndex] = uint32(valueTrie.Get(utf16.Encode([]rune(values[j]))))
				mappingIndex++
			}
		}
	}
	mappingBitVector := NewBitVector(mappingBitList.Words, uint32(mappingBitList.Size))

	return &CompactDictionary{
		keyTrie:          keyTrie,
		valueTrie:        valueTrie,
		mapping:          mapping,
		mappingBitVector: mappingBitVector,
	}
}
