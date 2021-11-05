package isbuilt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
)

const projectID string = "cloud-computing-327315"

func IndexBuiltHttp(w http.ResponseWriter, r *http.Request) {
	status := is_index_built()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resultStruct := struct {
		Status string `json:"status"`
	}{strconv.FormatBool(status)}
	json.NewEncoder(w).Encode(resultStruct)
}

func is_index_built() bool {
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

	// Get document for the given word
	_, err = client.Collection("ryoost-mapreduce").Doc("InverseIndex").Collection("Words").Doc("the").Get(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}
