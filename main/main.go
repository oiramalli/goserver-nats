package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	nats "github.com/nats-io/nats.go"
)

var nc, _ = nats.Connect("nats://nats:4222") // nats://localhost:4222

func main() {
	http.HandleFunc("/", rootHandler)
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)
	// Simple Async Subscriber
	nc.Subscribe("proyecto2", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	log.Println("Listening on :8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	type RData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	keys, _ := r.URL.Query()["msg"]
	nc.Publish("proyecto2", []byte(keys[0]))
	w.Header().Set("Content-Type", "application/json")
	s := `{ "status": "OK", "message": "` + keys[0] + `" }`
	data := &RData{}
	err := json.Unmarshal([]byte(s), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(data)
	w.Write(js)
}
