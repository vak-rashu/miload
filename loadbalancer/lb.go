package main

import (
	"encoding/json"
	"io"
	"log"
	http "net/http"
)

// what miload is?
// distributes clients requests load to the network servers
// ensure high avaibility by sending requests to servers that are online
// provides the flexibility to add or subtract servers as demand dictates

// step 1
// create a server that starts up, listens for incoming requests
// and sends them to another server to process

type msg struct {
	M string `json:"msg"`
}

type pointer struct {
	i int
}

var serverSlice []string = []string{
	"http://localhost:8010/1",
	"http://localhost:8011/2",
	"http://localhost:8012/3",
}

func (I *pointer) roundRobin([]string) string {
	if I.i == 0 {
		return serverSlice[I.i]
	}
	if I.i == len(serverSlice)-1 {
		I.i = 0
	}
	I.i += 1

	return serverSlice[I.i]
}

func main() {
	http.HandleFunc("/", callBackendHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

var v = pointer{i: 0}

func callBackendHandler(w http.ResponseWriter, r *http.Request) {
	server := v.roundRobin(serverSlice)

	sendMsg := msg{}
	backendMsg, err := callBackend(server)
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
	}
	sendMsg.M = string(backendMsg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendMsg.M)
}

func callBackend(server string) ([]byte, error) {
	resp, err := http.Get(server)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, nil

}
