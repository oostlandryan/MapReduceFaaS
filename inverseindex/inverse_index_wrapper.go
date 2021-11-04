package inverseindex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
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

	// Connect to Firestore
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

	// Store inverse index in Firestore, 500 words at a time
	count := 0
	var batch *firestore.WriteBatch
	for k, v := range inverseIndex {
		if count%500 == 0 {
			batch = client.Batch()
		}
		docRef := client.Collection("ryoost-mapreduce").Doc("InverseIndex").Collection("Words").Doc(k)
		batch.Set(docRef, v)
		if count%500 == 499 {
			_, err := batch.Commit(ctx)
			if err != nil {
				log.Printf("Error on batch write: %s", err)
				fmt.Fprint(w, "Inverse Index not stored")
				os.Exit(1)
			}
		}
		count++
	}
	fmt.Fprint(w, "Inverse Index Built")
}
