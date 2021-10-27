package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type mrTuple struct {
	WordFile string `json:"worldfile"`
	Count    int    `json:"count"`
}

func main() {
	f := []string{"PrideandPrejudice"}
	intermediate := mapCloud(f)
	result := reduceCloud(intermediate)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count < result[j].Count
	})
	log.Println(result)
}

func mapCloud(files []string) []mrTuple {
	mapFuncUrl := "https://us-central1-cloud-computing-327315.cloudfunctions.net/MapHttp"

	j := struct {
		Files []string `json:"Files"`
	}{files}

	postBody, _ := json.Marshal(j)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(mapFuncUrl, "application/json", responseBody)
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
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(reduceFuncUrl, "application/json", responseBody)
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
