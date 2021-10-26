package mapper

import (
	"strings"
	"unicode"
)

type mrTuple struct {
	wordFile string
	count    int
}

func WordCount(fileName string, text string) []mrTuple {
	var result []mrTuple
	words := strings.Fields(text)
	for _, w := range words {
		w = strings.TrimRightFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		tup := mrTuple{
			wordFile: w + "_" + fileName,
			count:    1,
		}
		result = append(result, tup)
	}
	return result
}
