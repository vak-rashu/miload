package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// backend server 1
type msg struct {
	M string `json:"msg"`
}

type heathCheck struct {
	Beat string `json:"beat"`
}

func main() {
	http.HandleFunc("/2", GetMsgHandler)
	http.HandleFunc("/heart-beat", GetHeartBeat)
	log.Fatal(http.ListenAndServe(":8011", nil))
}

func GetMsgHandler(w http.ResponseWriter, r *http.Request) {
	sendMsg := msg{}
	if r.Method == "GET" {
		sendMsg.M = "Hello from Backend 2"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sendMsg)
	} else {
		http.Error(w, "This method is prohibited", http.StatusMethodNotAllowed)
	}
}

func GetHeartBeat(w http.ResponseWriter, r *http.Request) {
	beat := heathCheck{}
	beat.Beat = "HI"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(beat)
}
