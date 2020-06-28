package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	nats "github.com/nats-io/nats.go"
)

var nc, _ = nats.Connect("nats://nats:4222") // nats://localhost:4222
var last = "{}"

func main() {
	http.HandleFunc("/", rootHandler)
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)
	// Simple Async Subscriber
	nc.Subscribe("proyecto2", func(m *nats.Msg) {
		last = string(m.Data)
	})
	log.Println("Listening on :8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		type RData struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		param1 := r.URL.Query().Get("msg")

		if param1 != "" {
			w.Header().Set("Content-Type", "application/json")
			s := `{ "status": "OK", "message": "` + param1 + `" }`
			data := &RData{}
			err := json.Unmarshal([]byte(s), data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			js, err := json.Marshal(data)
			w.Write(js)
			return
		}
		w.Write([]byte("Si llegaste acá, ya sabes que hacer."))
		return
	case "POST":
		if err := r.ParseMultipartForm(0); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		Nombre := r.FormValue("Nombre")
		Departamento := r.FormValue("Departamento")
		Edad := r.FormValue("Edad")
		FormaContagio := r.FormValue("FormaContagio")
		Estado := r.FormValue("Estado")
		mensaje := `{"Nombre":"` + Nombre + `, "Departamento":"` + Departamento + `, "Edad":"` + Edad + `, "FormaContagio":"` + FormaContagio + `, "Estado":"` + Estado + `"}`
		nc.Publish("proyecto2", []byte(mensaje))
		w.Write([]byte("Elemento previo: " + last + ". Enviando: " + mensaje))
	default:
		w.Write([]byte("Sorry, only GET and POST methods are supported."))
	}
}
