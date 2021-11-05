package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
)

const projectID string = "cloud-computing-327315"

func main() {
	result := is_index_built()
	fmt.Println(result)
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
