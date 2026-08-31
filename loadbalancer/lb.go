package main

import (
	"encoding/json"
	"io"
	"log"
	http "net/http"
	"net/url"
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

type request struct {
	body   string
	method string
	url    *url.URL
}

func main() {
	http.HandleFunc("/", callBackendHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func callBackendHandler(w http.ResponseWriter, r *http.Request) {
	sendMsg := msg{}
	backendMsg, err := callBackend()
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
	}
	sendMsg.M = string(backendMsg)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendMsg.M)
}

func callBackend() ([]byte, error) {

	resp, err := http.Get("http://localhost:8010/msg")
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, nil

}
