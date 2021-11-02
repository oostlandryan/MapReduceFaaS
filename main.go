package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

type mrTuple struct {
	WordFile string `json:"worldfile"`
	Count    int    `json:"count"`
}

func main() {
	inverseIndex([]string{"Frankenstein",
		"Pride and Prejudice",
		"The Game of Go",
		"The Legend of Sleepy Hollow",
		"Alice's Adventures in Wonderland",
		"Dracula",
		"The Scarlet Letter",
		"A Christmas Carol",
		"The Adventures of Sherlock Holmes",
		"The Yellow Wallpaper",
		"The Picture of Dorian Gray",
		"A Tale of Two Cities",
		"The Strange Case of Dr. Jekyll And Mr. Hyde",
		"The Great Gatsby"}, 14, 100)
}

/*
inverseIndex creates an inverse index of the given files (as long as they're in the firestore)
using m mappers and r reducers
*/
func inverseIndex(files []string, m int, r int) {
	// if m or r is less than 1, set it to 1
	if m < 1 {
		m = 1
	}
	if r < 1 {
		r = 1
	}

	// if m is greater than length of files, set it to length of files
	if m > len(files) {
		m = len(files)
	}

	// Partition files into m slices
	fSlices := make([][]string, m)
	for i, f := range files {
		i = i % m
		fSlices[i] = append(fSlices[i], f)
	}
	mStart := time.Now()
	// Concurrently run m map functions
	mapChan := make(chan []mrTuple)
	for i := 0; i < m; i++ {
		go func(s int) {
			mapChan <- mapCloud(fSlices[s])
		}(i)
	}

	/*
		Collect Results of m map functions
		NOTE: This acts as the barrier, because result := <-mapChan will wait until there is a value
		in mapChan that can be assigned to result.
	*/
	var mapResult []mrTuple
	for i := 0; i < m; i++ {
		result := <-mapChan
		mapResult = append(mapResult, result...)
	}
	mElapsed := time.Since(mStart)
	// Sort mapResult
	sort.Slice(mapResult, func(i, j int) bool {
		return mapResult[i].WordFile < mapResult[j].WordFile
	})

	// Find indexes where wordfile changes
	indexes := []int{}
	prev := mapResult[0].WordFile
	for i := 1; i < len(mapResult); i++ {
		if mapResult[i].WordFile != prev {
			indexes = append(indexes, i)
			prev = mapResult[i].WordFile
		}
	}

	// Partition mapResult to r reducers
	rSlices := make([][]mrTuple, r)
	rSlices[0] = mapResult[:indexes[0]]
	for i := 1; i < len(indexes); i++ {
		reducer := i % r
		rSlices[reducer] = append(rSlices[reducer], mapResult[indexes[i-1]:indexes[i]]...)
	}
	rStart := time.Now()
	// Concurrently run r reduce functions
	reduceChan := make(chan []mrTuple)
	for i := 0; i < r; i++ {
		go func(s int) {
			reduceChan <- reduceCloud(rSlices[s])
		}(i)
	}

	// Collect results of r reduce functions
	var reduceResult []mrTuple
	for i := 0; i < r; i++ {
		result := <-reduceChan
		reduceResult = append(reduceResult, result...)
	}
	rElapsed := time.Since(rStart)
	fmt.Println(reduceResult)
	log.Println(len(reduceResult))
	fmt.Printf("Map Time: %s\n", mElapsed)
	fmt.Printf("Reduce Time: %s\n", rElapsed)

}

/*
mapCloud
*/
func mapCloud(files []string) []mrTuple {
	mapFuncUrl := "https://us-central1-cloud-computing-327315.cloudfunctions.net/MapHttp"

	j := struct {
		Files []string `json:"Files"`
	}{files}

	postBody, _ := json.Marshal(j)
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(mapFuncUrl, "application/json", requestBody)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	output := struct {
		MapResult []mrTuple `json:"mapresult"`
	}{}
	json.Unmarshal(body, &output)

	return output.MapResult
}

func reduceCloud(input []mrTuple) []mrTuple {
	reduceFuncUrl := "https://us-central1-cloud-computing-327315.cloudfunctions.net/ReduceHttp"

	j := struct {
		MapResult []mrTuple `json:"mapresult"`
	}{input}

	postBody, _ := json.Marshal(j)
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(reduceFuncUrl, "application/json", requestBody)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	output := struct {
		ReduceResult []mrTuple `json:"reduceresult"`
	}{}
	json.Unmarshal(body, &output)

	return output.ReduceResult
}