package mapper

import (
	"strings"
	"unicode"
)

type mrTuple struct {
	WordFile string `json:"worldfile"`
	Count    int    `json:"count"`
}

func WordCount(fileName string, text string) []mrTuple {
	var result []mrTuple
	words := strings.Fields(text)
	for _, w := range words {
		w = strings.TrimRightFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		w = strings.TrimLeftFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		tup := mrTuple{
			WordFile: w + "_" + fileName,
			Count:    1,
		}
		if w != "" {
			result = append(result, tup)
		}

	}
	return result
}
