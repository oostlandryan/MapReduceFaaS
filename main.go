package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type mrTuple struct {
	WordFile string `json:"worldfile"`
	Count    int    `json:"count"`
}

func main() {
	f := []string{"testfile"}
	log.Println(mapCloud(f))
}

func mapCloud(files []string) []mrTuple {
	reduceFuncUrl := "https://us-central1-cloud-computing-327315.cloudfunctions.net/MapHttp"

	j := struct {
		Files []string `json:"Files"`
	}{files}

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
		MapResult []mrTuple `json:"mapresult"`
	}{}
	json.Unmarshal(body, &output)

	return output.MapResult
}

func reduceCloud() {

}
