package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var listener net.Listener
var done chan bool

func StartTcp(player *Player) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error when listening on port 8081", err)
	}

	log.Println("Listening on TCP port :8081")
	done = make(chan bool)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(
				"Error when accepting a TCP connection",
				err)
		}

		go handleConn(conn, player)
	}
}

func Close() {
	listener.Close()
}

func handleConn(conn net.Conn, p *Player) {

	read := make(chan string, 5)

	go readConn(conn, read)

	writeStatus(conn, p.Status())

	for {
		select {
		case data := <-read:
			if !handleRead(conn, p, data) {
				// todo: find a better way to handle disconnects from client
				log.Println("Client disconnected. Stoping handleConn()")
				return
			}
		case <-time.After(5 * time.Second):
			writeStatus(conn, p.Status())
		}
	}
}

func handleRead(conn net.Conn, p *Player, data string) bool {
	tokens := strings.Split(data, " ")

	if len(tokens) < 1 {
		return true
	}

	cmd := tokens[0]

	switch cmd {
	case "shutdown":
		close(done)
		return false
	case "closed":
		return false

	case "+":
		// vol up
		p.ChangeVolume(10)
	case "-":
		p.ChangeVolume(-10)
	case "l":
		if len(tokens) != 2 {
			return true
		}
		p.LoadSong(tokens[1])
	case "p":
		p.TogglePause()
	case "s":
		p.Stop()
	}
	writeStatus(conn, p.Status())

	return true
}

func writeStatus(conn net.Conn, status *PlayerStatus) {
	statusBytes := toBytes(status)
	n, err := conn.Write(statusBytes)
	if err != nil {
		log.Println("Error when writting to conn", err)
	}

	if n != len(statusBytes) {
		log.Println("Warning! Not all the bytes of the status were sent")
	}
}

func toBytes(status *PlayerStatus) []byte {
	var b bytes.Buffer
	b.WriteString("Status\r\n")

	vol := strconv.FormatInt(int64(status.Volume), 10)

	b.WriteString("File: " + filepath.Join(status.Folder, status.Filename))

	b.WriteString("\r\n")

	b.WriteString("Volume: " + vol)

	b.WriteString("\tPaused: " + strconv.FormatBool(status.IsPaused))

	b.WriteString("\tLoaded: " + strconv.FormatBool(status.Loaded))

	b.WriteString("\tLooping: " + strconv.FormatBool(status.Looping))

	b.WriteString("\r\n")

	return b.Bytes()
}

func readConn(conn net.Conn, data chan<- string) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString(byte('\n'))
		data <- strings.TrimSuffix(line, "\r\n")

		if err != nil {
			data <- "closed"
			log.Println("Error when reading from connection", err)
			return
		}
	}
}
