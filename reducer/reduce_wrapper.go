package reducer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReduceHttp(w http.ResponseWriter, r *http.Request) {
	// ensure it's a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	input := struct {
		MapResult []mrTuple `json:"mapresult"`
	}{}
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}
	// Reduce input
	output := reduce(input.MapResult)

	w.Header().Set("Content-Type", "application/json")
	resultStruct := struct {
		ReduceResult []mrTuple `json:"reduceresult"`
	}{output}
	json.NewEncoder(w).Encode(resultStruct)
}
