package mapper

import (
	"log"
	"regexp"
	"strings"
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
		w := reg.ReplaceAllString(w, "")

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
