package searchindex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type SearchRequest struct {
	Term string `json:"term"`
}

func RyoostSearchHttp(w http.ResponseWriter, r *http.Request) {
	// ensure it's a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	var sr SearchRequest
	err := decoder.Decode(&sr)
	if err != nil {
		panic(err)
	}

	// Get results of search
	results := SearchIndex(sr.Term)

	// order by count
	type tcTuple struct {
		Title string `json:"title"`
		Count int    `json:"count"`
	}
	var sortedResults []tcTuple
	for k, v := range results {
		sortedResults = append(sortedResults, tcTuple{k, v})
	}
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Count > sortedResults[j].Count
	})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resultStruct := struct {
		Results []tcTuple `json:"results"`
	}{sortedResults}
	json.NewEncoder(w).Encode(resultStruct)
}
