package mapper

import (
	"log"
	"regexp"
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
		reg, err := regexp.Compile("[^a-zA-Z0-9']+")
		if err != nil {
			log.Fatal(err)
		}
		w := reg.ReplaceAllString(w, " ")
		ws := strings.Split(w, " ")
		for _, v := range ws {
			v = strings.TrimFunc(v, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
			tup := mrTuple{
				WordFile: v + "_" + fileName,
				Count:    1,
			}
			if v != "" {
				result = append(result, tup)
			}
		}
	}
	return result
}
