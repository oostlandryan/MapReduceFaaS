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
	// TODO Change to batch writes
	for k, v := range inverseIndex {
		_, err = client.Collection("ryoost-mapreduce").Doc("InverseIndex").Collection("Words").Doc(k).Set(ctx, v)
		if err != nil {
			fmt.Printf("An error has occurred: %s", err)
			fmt.Fprint(w, "Inverse Index Couldn't be stored")
			return
		}
	}
	fmt.Fprint(w, "Inverse Index Built")
}
