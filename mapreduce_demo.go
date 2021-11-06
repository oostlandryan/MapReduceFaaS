package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const searchFuncUrl string = "https://us-central1-cloud-computing-327315.cloudfunctions.net/RyoostSearchHttp"
const createFuncUrl string = "https://us-central1-cloud-computing-327315.cloudfunctions.net/RyoostCreateIndexHttp"

func main() {
	var word string
	flag.StringVar(&word, "word", "the", "Word to search")
	var m int
	flag.IntVar(&m, "m", 25, "Number of mapper functions to use")
	var r int
	flag.IntVar(&r, "r", 30, "Number of reducer functions to use")
	var index bool
	flag.BoolVar(&index, "createindex", true, "'true' or 'false' create a new inverted index")
	flag.Parse()

	if index {
		fmt.Println(createIndexCloud(m, r))
	}

	fmt.Println(searchCloud(word))
}

type tcTuple struct {
	Title string `json:"title"`
	Count int    `json:"count"`
}

func createIndexCloud(m, r int) string {
	f := []string{"Frankenstein", "Pride and Prejudice", "The Legend of Sleepy Hollow", "Alice's Adventures in Wonderland", "Dracula", "The Scarlet Letter", "A Christmas Carol", "The Adventures of Sherlock Holmes", "The Yellow Wallpaper", "The Picture of Dorian Gray", "A Tale of Two Cities", "The Strange Case of Dr. Jekyll And Mr. Hyde", "The Great Gatsby", "A Doll's House", "A Modest Proposal", "Metamorphosis", "The Prince", "Heart of Darkness", "The Odyssey", "Grimms' Fairy Tales", "Beowulf", "The Adventures of Tom Sawyer", "Emma", "The Communist Manifesto", "Anthem"}
	j := struct {
		Files    []string
		Mappers  int
		Reducers int
	}{f, m, r}

	postBody, _ := json.Marshal(j)
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(createFuncUrl, "application/json", requestBody)
	if err != nil {
		log.Fatalln("3", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("4", err)
	}
	return string(body)
}

func searchCloud(word string) string {
	j := struct {
		Term string
	}{word}

	postBody, _ := json.Marshal(j)
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(searchFuncUrl, "application/json", requestBody)
	if err != nil {
		log.Fatalln("1", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("2", err)
	}
	return string(body)
}
