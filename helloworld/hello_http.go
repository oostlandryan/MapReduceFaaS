package helloworld

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

// HelloHTTP is a HTTP Cloud Function with a request parameter
func HelloHTTP(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello, World!\n")
		return
	}
	if d.Name == "" {
		fmt.Fprint(w, "Hello, World!\n")
		return
	}
	fmt.Fprintf(w, "Hello, %s!\n", html.EscapeString(d.Name))
}
