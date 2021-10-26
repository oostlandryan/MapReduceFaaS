package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

const projectID string = "cloud-computing-327315"

type files_struct struct {
	Files []string
}

// Google Cloud Function to call the map function on
func MapHttp(w http.ResponseWriter, r *http.Request) {
	// ensure it's a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// decode response
	decoder := json.NewDecoder(r.Body)
	var f files_struct
	err := decoder.Decode(&f)
	if err != nil {
		panic(err)
	}

	// Connect to Firestore DB
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Get results of given files
	var result []mrTuple
	c := client.Collection("files")
	for _, f := range f.Files {
		d := c.Doc(f)
		dsnap, err := d.Get(ctx)
		if err != nil {
			log.Fatalln(err)
			fmt.Fprint(w, "error1 ")
		}
		s := fmt.Sprintf("%v", dsnap.Data()["text"])
		words := WordCount(f, s)
		result = append(result, words...)
	}

	// Return results

	w.Header().Set("Content-Type", "application/json")
	resultStruct := struct {
		MapResult []mrTuple `json:"mapresult"`
	}{result}
	json.NewEncoder(w).Encode(resultStruct)
}
