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

func main() {
	http.HandleFunc("/2", GetMsgHandler)
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
