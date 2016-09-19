package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bradacina/mp3player/internal/httpChain"
)

type errorMessage struct {
	isError bool
	message string
}

var p *Player

func init() {
	http.Handle("/", http.FileServer(http.Dir("./webapp")))
	http.Handle("/volume",
		httpChain.
			NewChainWithFunc(checkPlayer).
			NextFunc(volume))
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/togglepause", togglePause)
	http.HandleFunc("/togglelooping", toggleLooping)
	http.HandleFunc("/list", list)
	http.HandleFunc("/load", load)
	http.Handle("/status",
		httpChain.NewChainWithFunc(checkGet).
			NextFunc(checkPlayer).
			NextFunc(status))
}

func StartWeb(player *Player) {
	p = player
	log.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot startup server on :8080", err)
	}
}

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

func togglePause(writer http.ResponseWriter, request *http.Request) {

}

func toggleLooping(writer http.ResponseWriter, request *http.Request) {

}

func load(writer http.ResponseWriter, request *http.Request) {

}

func list(writer http.ResponseWriter, request *http.Request) {

}

func checkPlayer(writer http.ResponseWriter, request *http.Request) {
	if p == nil {
		httpChain.BreakChain(request)
		encoded, err := json.Marshal(&errorMessage{isError: true, message: "no mp3 player active"})
		if err != nil {
			log.Println("error encoding errorReponse", err)
		}

		writer.Write(encoded)
	}
}

func checkGet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		httpChain.BreakChain(request)
		http.NotFound(writer, request)
	}
}

func checkPost(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		httpChain.BreakChain(request)
		http.NotFound(writer, request)
	}
}
