package inverseindex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

const projectID string = "cloud-computing-327315"

type CreateIndexRequest struct {
	Files    []string
	Mappers  int
	Reducers int
}

func CreateIndexHttp(w http.ResponseWriter, r *http.Request) {
	// ensure it's a post request
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// decode request
	decoder := json.NewDecoder(r.Body)
	var cir CreateIndexRequest
	err := decoder.Decode(&cir)
	if err != nil {
		panic(err)
	}

	// Create inverse index
	inverseIndex := inverseIndex(cir.Files, cir.Mappers, cir.Reducers)

	// Connect to Firebase
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
	// Store inverse index in Firestore
	_, err = client.Collection("ryoost-mapreduce").Doc("InverseIndex").Set(ctx, inverseIndex)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		fmt.Printf("An error has occurred: %s", err)
		fmt.Fprint(w, "Inverse Index Couldn't be stored")
	} else {
		fmt.Fprint(w, "Inverse Index Built")
	}

}
