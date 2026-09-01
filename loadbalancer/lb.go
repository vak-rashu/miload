package main

import (
	"encoding/json"
	"io"
	"log"
	http "net/http"
	"slices"
	"time"
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

// list of backend server
const backend1 string = "http://localhost:8010/1"
const backend2 string = "http://localhost:8011/2"
const backend3 string = "http://localhost:8012/3"

// list of backend server health check points
const checkHealth1 string = "http://localhost:8010/heart-beat"
const checkHealth2 string = "http://localhost:8011/heart-beat"
const checkHealth3 string = "http://localhost:8012/heart-beat"

// round robin slice
var serverSlice []string = []string{
	backend1,
	backend2,
	backend3,
}

// map to check health of backend servers
var serverMap map[string]string = map[string]string{
	checkHealth1: backend1,
	checkHealth2: backend2,
	checkHealth3: backend3,
}

// a round robin algo takes the input: a slice
func (I *pointer) roundRobin([]string) string {
	if I.i == 0 {
		return serverSlice[I.i]
	}
	if I.i == len(serverSlice) {
		I.i = 0
	}

	return serverSlice[I.i]
}

func main() {
	go checkHeartBeat()

	http.HandleFunc("/", callBackendHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

var I = &pointer{i: 0}

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
	serverURL := I.roundRobin(serverSlice)
	I.i += 1
	resp, err := http.Get(serverURL)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, nil

}

func checkHealth(serverURL string) ([]byte, error) {
	resp, err := http.Get(serverURL)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, nil
}

var stoppedServer string

func checkHeartBeat() {
	s := time.Duration(2 * time.Second)
	for url := range serverMap {
		_, err := checkHealth(url)

		if err != nil {
			stoppedServer = serverMap[url]
			i := slices.Index(serverSlice, stoppedServer)
			serverSlice = slices.Concat(serverSlice[:i], serverSlice[i:])
		}

		if hasServer := slices.Contains(serverSlice, stoppedServer); hasServer == false {
			serverSlice = append(serverSlice, stoppedServer)
		}
	}

	time.Sleep(s)
}
