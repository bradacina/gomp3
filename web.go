package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.Handle("/", http.FileServer(http.Dir("./webapp")))
	http.HandleFunc("/volume", volume)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/pause", pause)
	http.HandleFunc("/load", load)

	p = ThePlayer()
}

func StartWeb() {
	log.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot startup server on :8080", err)
	}
}

var p *Player

func status(writer http.ResponseWriter, request *http.Request) {
	s := p.Status()
	encoded, err := json.Marshal(s)
	if err != nil {
		log.Println("Error when encoding status", err)
		return
	}

	writen, err := writer.Write(encoded)
	if err != nil {
		log.Println("Error when writing status", err)
		return
	}

	if writen != len(encoded) {
		log.Println(
			"Not all status data was written? Writen:", writen,
			"Len:", len(encoded))
	}
}

func volume(writer http.ResponseWriter, request *http.Request) {

}

func stop(writer http.ResponseWriter, request *http.Request) {

}

func pause(writer http.ResponseWriter, request *http.Request) {

}

func load(writer http.ResponseWriter, request *http.Request) {

}
