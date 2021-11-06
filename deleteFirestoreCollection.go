package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

const projectID string = "cloud-computing-327315"

func main() {
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

	// Delete Inverted Index firebase
	ref := client.Collection("ryoost-mapreduce").Doc("InverseIndex").Collection("Words")
	err = deleteCollection(ref, ctx, client)
	if err != nil {
		fmt.Println("Error Deleting Inverted Index: ", err)
	} else {
		fmt.Println("Inverted Index Deleted")
	}

	ref = client.Collection("ryoost-mapreduce")
	err = deleteCollection(ref, ctx, client)
	if err != nil {
		fmt.Println("Error Deleting Corpus: ", err)
	} else {
		fmt.Println("Corpus Deleted")
	}

}

func deleteCollection(ref *firestore.CollectionRef, ctx context.Context, client *firestore.Client) error {
	for {
		iter := ref.Limit(500).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			break
		}
		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
