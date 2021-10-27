package reducer

type mrTuple struct {
	WordFile string `json:"worldfile"`
	Count    int    `json:"count"`
}

func reduce(input []mrTuple) []mrTuple {
	var resultMap = make(map[string]int)
	// Sum the word counts
	for _, t := range input {
		if _, ok := resultMap[t.WordFile]; ok {
			resultMap[t.WordFile] = resultMap[t.WordFile] + 1
		} else {
			resultMap[t.WordFile] = 1
		}
	}

	// Turn map into list and return
	var result []mrTuple
	for k, v := range resultMap {
		result = append(result, mrTuple{WordFile: k, Count: v})
	}

	return result
}
