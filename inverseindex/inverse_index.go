package inverseindex

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
	chunkSize := int(len(files) / m)
	for i := 0; i < m-1; i++ {
		fSlices[i] = files[i*chunkSize : (i+1)*chunkSize]
	}
	fSlices[m-1] = files[(m-1)*chunkSize:]

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

	// Sort mapResult

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
