package main

import (
	"encoding/json"
	"fmt"
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
	s []string
}

// list of backend server
const backend1 string = "http://localhost:8010/1"
const backend2 string = "http://localhost:8011/2"
const backend3 string = "http://localhost:8012/3"

// list of backend server health check points
const checkHealth1 string = "http://localhost:8010/heart-beat1"
const checkHealth2 string = "http://localhost:8011/heart-beat2"
const checkHealth3 string = "http://localhost:8012/heart-beat3"

// round robin slice
var serverSlice []string = []string{
	backend1,
	backend2,
	backend3,
}

// map to check health of backend servers
var serverMap map[string]string = map[string]string{
	backend1: checkHealth1,
	backend2: checkHealth2,
	backend3: checkHealth3,
}

// a round robin algo takes the input: a slice
// func (I *pointer) roundRobin(c chan string) {
// 	currInd := I.i
// 	// use modulus to loop again and again
// 	// for the same range of i values
// 	I.i = (I.i + 1) % len(I.s)
// 	serverURL := I.s[currInd]
// 	c <- serverURL

// 	// return I.s[currInd]
// }

func (I *pointer) roundRobin() string {
	currInd := I.i
	// use modulus to loop again and again
	// for the same range of i values
	I.i = (I.i + 1) % len(I.s)

	return I.s[currInd]
}

var I = &pointer{i: 0, s: serverSlice}

func main() {

	go callHeartBeat()

	http.HandleFunc("/", callBackendHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func callBackendHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		sendMsg := msg{}
		backendMsg, err := callBackend()
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
		}
		sendMsg.M = string(backendMsg)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sendMsg.M)
	} else {
		http.Error(w, "This method is prohibited", http.StatusMethodNotAllowed)
	}

}

func callBackend() ([]byte, error) {
	serverURL := I.roundRobin()
	resp, err := http.Get(serverURL)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, nil

}

// may be the response needs to be mapped out to a json
func checkHeartBeat(server string) ([]byte, error) {
	resp, err := http.Get(server)
	if err != nil {
		return []byte{}, fmt.Errorf("error making response:%v", err)
	}

	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading response:%v", err)
	}
	return r, nil
}

var cTime = make(chan time.Duration)

func getTickerTime() {
	cTime <- time.Duration(time.Second * 2)
}

func callHeartBeat() {
	for {
		// time.Sleep(time.Duration(time.Second * 2))
		getTickerTime()
		select {
		case <-cTime:
			for serverURL, healthURL := range serverMap {
				_, err := checkHeartBeat(healthURL)
				if err != nil {
					ind := slices.Index(I.s, serverURL)
					if ind == -1 {
						continue
					}
					I.s = append(I.s[:ind], I.s[ind+1:]...)
				} else {
					if hasActive := slices.Contains(I.s, serverURL); hasActive != true {
						I.s = append(I.s, serverURL)
					}
				}
			}
		}
		// time.NewTicker(time.Duration(time.Second * 2))

	}
}
