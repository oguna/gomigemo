package migemo

import (
	"strings"
	"sort"
)

type SimpleDictionary struct {
	keys []string;
	values []string;
}

type KeyValuePair struct {
	key string;
	value string;
}

func BuildSimpleDictionary(file string) *SimpleDictionary {
	var lines = strings.Split(file, "\n")
	var keyValuePairs = []KeyValuePair{}
	for i := 0; i < len(lines); i++ {
		var line = lines[i]
		if (!strings.HasPrefix(line, ";") && len(line) != 0) {
			var semicolonPos = strings.Index(line, "\t");
			var key = line[:semicolonPos]
			var value = line[:semicolonPos + 1]
			keyValuePairs = append(keyValuePairs, KeyValuePair {
				key: key,
				value: value,
			})
		}
	}
	sort.Slice(keyValuePairs, func(i, j int) bool {
		var left = keyValuePairs[i].key
		var right = keyValuePairs[j].key
		var minlen = len(left)
		if (minlen > len(right)) {
			minlen = len(right)
		}
		for k := 0; k < minlen; k++ {
			if left[k] == right[k] {
				continue
			} else if left[k] < right[k] {
				return true
			} else {
				return false
			}
		}
		return len(left) < len(right)
	})
	var keys = make([]string, len(keyValuePairs))
	for i, v := range keyValuePairs {
		keys[i] = v.key
	}
	var values = make([]string, len(keyValuePairs))
	for i, v := range keyValuePairs {
		values[i] = v.value
	}
	return &SimpleDictionary {
		keys: keys,
		values: values,
	}
}

func (this *SimpleDictionary) PredictiveSearch(hiragana string) []string {
	if len(hiragana) == 0 {
		return []string {}
	}
	var hiraganaRune = []rune(hiragana)
	var stop = string(append(hiraganaRune[:len(hiragana) - 1], hiraganaRune[len(hiragana) - 1]));
	var startPos = binarySearchString(this.keys, 0, len(this.keys), hiragana);
	if (startPos < 0) {
		startPos = -(startPos + 1);
	}
	var endPos = binarySearchString(this.keys, 0, len(this.keys), stop);
	if (endPos < 0) {
		endPos = -(endPos + 1);
	}
	var result = []string {};
	for i := startPos; i < endPos; i++ {
		for _, j := range strings.Split(this.values[i], "\t") {
			result = append(result, j);
		}
	}
	return result;
}

func binarySearchString(a []string, fromIndex int, toIndex int, key string) int {
	var low = fromIndex;
	var high = toIndex - 1;
	for (low <= high) {
		var mid = (low + high) >> 1;
		var midVal = a[mid];
		if (midVal < key) {
			low = mid + 1;
		} else if (midVal > key) {
			high = mid - 1;
		} else {
			return mid;
		}
	}
	return -(low + 1);
}