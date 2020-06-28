package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
		param1 := r.URL.Query().Get("msg")
		if param1 != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{ "status": "OK", "status_code":"1", "message": "` + param1 + `" }`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{status:"OK", "status_code":"1", message: "Si llegaste acá, ya sabes que hacer."}`))
		return
	case "POST":
		type Person struct {
			Nombre        string `json:"Nombre"`
			Departamento  string `json:"Departamento"`
			Edad          int    `json:"Edad"`
			FormaContagio string `json:"Forma de contagio"`
			Estado        string `json:"Estado"`
		}
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"status":"FAILED","status_code":"0","message":"El cuerpo del mensaje no tiene el formato correcto."}`, http.StatusBadRequest)
			return
		}
		Nombre := p.Nombre
		Departamento := p.Departamento
		Edad := p.Edad
		FormaContagio := p.FormaContagio
		Estado := p.Estado
		mensaje := `{"Nombre":"` + Nombre + `", "Departamento":"` + Departamento + `", "Edad":` + strconv.Itoa(Edad) + `, "FormaContagio":"` + FormaContagio + `", "Estado":"` + Estado + `"}`

		if err := nc.Publish("proyecto2", []byte(mensaje)); err != nil {
			http.Error(w, `{"status":"FAILED","status_code":"0","message":"No se logró publicar en el canal. ¿Está dosponible el servidor de NATS?"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{status:"OK", "status_code":"1", "data:" {"Elemento previo": ` + last + `,"Enviando": ` + mensaje + `}}`))
	default:
		http.Error(w, `{"status":"FAILED","status_code":"0","message":"Opps, solamente se soportan los métodos GET y POST."}`, http.StatusBadRequest)
		return
	}
}
