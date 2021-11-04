package searchindex

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const projectID string = "cloud-computing-327315"

// Search index returns the map of books to counts for the given word
func SearchIndex(word string) map[string]int {
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
	dsnap, err := client.Collection("ryoost-mapreduce").Doc("InverseIndex").Collection("Words").Doc(word).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return map[string]int{}
	} else if err != nil {
		log.Fatalln(err)
	}
	result := make(map[string]int)
	for k, v := range dsnap.Data() {
		result[k] = int(v.(int64))
	}
	return result
}
