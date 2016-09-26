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

type loadSongMessage struct {
	Song string
}

var p *Player

func init() {
	http.Handle("/", http.FileServer(http.Dir("./webapp")))
	http.Handle("/volumeup",
		httpChain.NewChainWithFunc(checkPlayer).
			NextFunc(checkPost).
			NextFunc(volumeUp))
	http.Handle("/volumedown",
		httpChain.NewChainWithFunc(checkPlayer).
			NextFunc(checkPost).
			NextFunc(volumeDown))
	http.Handle("/stop",
		httpChain.NewChainWithFunc(checkPost).
			NextFunc(checkPlayer).
			NextFunc(stop))
	http.Handle("/togglepause",
		httpChain.NewChainWithFunc(checkPost).
			NextFunc(checkPlayer).
			NextFunc(togglePause))
	http.Handle("/togglelooping",
		httpChain.NewChainWithFunc(checkPost).
			NextFunc(checkPlayer).
			NextFunc(toggleLooping))
	http.Handle("/list",
		httpChain.NewChainWithFunc(checkGet).
			NextFunc(checkPlayer).
			NextFunc(list))
	http.Handle("/load",
		httpChain.NewChainWithFunc(checkPost).
			NextFunc(checkPlayer).
			NextFunc(load))
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
		internalError(writer)
		return
	}

	writeContentType(writer)
	writen, err := writer.Write(encoded)
	if err != nil {
		log.Println("Error when writing status", err)
		internalError(writer)
		return
	}

	if writen != len(encoded) {
		log.Println(
			"Not all status data was written? Writen:", writen,
			"Len:", len(encoded))
	}
}

func volumeUp(writer http.ResponseWriter, request *http.Request) {
	p.ChangeVolume(10)
}

func volumeDown(writer http.ResponseWriter, request *http.Request) {
	p.ChangeVolume(-10)
}

func stop(writer http.ResponseWriter, request *http.Request) {
	p.Stop()
}

func togglePause(writer http.ResponseWriter, request *http.Request) {
	p.TogglePause()
}

func toggleLooping(writer http.ResponseWriter, request *http.Request) {
	p.ToggleLooping()
}

func load(writer http.ResponseWriter, request *http.Request) {
	buf := make([]byte, request.ContentLength)

	request.Body.Read(buf)

	msg := loadSongMessage{}
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		log.Println("error unmarshaling loadSongMessage", err)
		internalError(writer)
	}

	err = p.LoadSong(msg.Song)
	if err != nil {
		log.Println("error when loading song", err)
		internalError(writer)
	}
}

func list(writer http.ResponseWriter, request *http.Request) {
	s, err := p.ListSongs()
	if err != nil {
		log.Println("Error when retrieving list", err)
		internalError(writer)
		return
	}

	encoded, err := json.Marshal(s)
	if err != nil {
		log.Println("Error when encoding list of songs", err)
		internalError(writer)
		return
	}

	writeContentType(writer)

	writen, err := writer.Write(encoded)
	if err != nil {
		log.Println("Error when writing list of songs", err)
		internalError(writer)
		return
	}

	if writen != len(encoded) {
		log.Println(
			"Not all list of songs data was written? Writen:", writen,
			"Len:", len(encoded))
	}
}

func checkPlayer(writer http.ResponseWriter, request *http.Request) {
	if p == nil {
		httpChain.BreakChain(request)
		writeError(writer, http.StatusInternalServerError, "no mp3 player active")
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

func writeContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/json")
}

func internalError(writer http.ResponseWriter) {
	http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
}

func writeError(writer http.ResponseWriter, errorCode int, message string) {
	encoded, err := json.Marshal(&errorMessage{isError: true, message: message})
	if err != nil {
		log.Println("error encoding errorReponse", err)
	}

	writer.WriteHeader(errorCode)
	writeContentType(writer)
	n, err := writer.Write(encoded)
	if err != nil {
		log.Println("error writing errorResponse", err)
	}

	if n != len(encoded) {
		log.Println("not all bytes of the errorResponse were written")
	}
}
