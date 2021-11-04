package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

// Here are some of the most downloaded books on Project Gutenberg, not counting those exceeding firestore's 1MB document limit
var bookURLs = map[string]string{
	"Frankenstein":                                "https://www.gutenberg.org/files/84/84-0.txt",
	"Pride and Prejudice":                         "https://www.gutenberg.org/files/1342/1342-0.txt",
	"The Legend of Sleepy Hollow":                 "https://www.gutenberg.org/files/41/41-0.txt",
	"Alice's Adventures in Wonderland":            "https://www.gutenberg.org/files/11/11-0.txt",
	"Dracula":                                     "https://www.gutenberg.org/files/345/345-0.txt",
	"The Scarlet Letter":                          "https://www.gutenberg.org/cache/epub/25344/pg25344.txt",
	"A Christmas Carol":                           "https://www.gutenberg.org/cache/epub/46/pg46.txt",
	"The Adventures of Sherlock Holmes":           "https://www.gutenberg.org/files/1661/1661-0.txt",
	"The Yellow Wallpaper":                        "https://www.gutenberg.org/files/1952/1952-0.txt",
	"The Picture of Dorian Gray":                  "https://www.gutenberg.org/files/174/174-0.txt",
	"A Tale of Two Cities":                        "https://www.gutenberg.org/files/98/98-0.txt",
	"The Strange Case of Dr. Jekyll And Mr. Hyde": "https://www.gutenberg.org/files/43/43-0.txt",
	"The Great Gatsby":                            "https://www.gutenberg.org/cache/epub/64317/pg64317.txt",
	"A Doll's House":                              "https://www.gutenberg.org/files/2542/2542-0.txt",
	"A Modest Proposal":                           "https://www.gutenberg.org/files/1080/1080-0.txt",
	"Metamorphosis":                               "https://www.gutenberg.org/files/5200/5200-0.txt",
	"The Prince":                                  "https://www.gutenberg.org/files/1232/1232-0.txt",
	"Heart of Darkness":                           "https://www.gutenberg.org/files/219/219-0.txt",
	"The Odyssey":                                 "https://www.gutenberg.org/cache/epub/1727/pg1727.txt",
	"Grimms' Fairy Tales":                         "https://www.gutenberg.org/files/2591/2591-0.txt",
	"Beowulf":                                     "https://www.gutenberg.org/files/16328/16328-0.txt",
	"The Adventures of Tom Sawyer":                "https://www.gutenberg.org/files/74/74-0.txt",
	"Emma":                                        "https://www.gutenberg.org/files/158/158-0.txt",
	"The Communist Manifesto":                     "https://www.gutenberg.org/cache/epub/61/pg61.txt",
	"Anthem":                                      "https://www.gutenberg.org/files/1250/1250-0.txt",
}

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

	for title, url := range bookURLs {
		// Download file from Project Gutenberg
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)

		// Add file to firebase
		_, err = client.Collection("ryoost-mapreduce").Doc(title).Set(ctx, map[string]interface{}{
			"text": sb,
		})
		if err != nil {
			log.Fatalf("Failed %s\n", title)
		}
	}
	log.Printf("%d Files Uploaded\n", len(bookURLs))
}
