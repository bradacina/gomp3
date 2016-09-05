package main

import (
	"bytes"
	"log"
	"net"
	"strconv"
	"time"
)

var listener net.Listener

func StartTcp() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error when listening on port 8081", err)
	}

	log.Println("Listening on TCP port :8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(
				"Error when accepting a TCP connection",
				err)
		}

		go handleConn(conn)
	}
}

func Close() {
	listener.Close()
}

func handleConn(conn net.Conn) {
	p := ThePlayer()

	read := make(chan []byte, 5)

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

func handleRead(conn net.Conn, p *Player, data []byte) bool {
	cmd := string(data)

	switch cmd {
	case "closed":
		return false

	case "+":
		// vol up
		p.ChangeVolume(10)
	case "-":
		p.ChangeVolume(-10)
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
	b.WriteString("Status\n")

	vol := strconv.FormatInt(int64(status.Volume), 10)

	b.WriteString("Volume: " + vol)

	b.WriteString("\tPaused: " + strconv.FormatBool(status.Paused))

	b.WriteString("\tStoped: " + strconv.FormatBool(status.Stopped))

	// todo: write name of currently playing song
	// todo: write percentage song completion

	b.WriteString("\n")

	return b.Bytes()
}

func readConn(conn net.Conn, data chan<- []byte) {
	buf := make([]byte, 10)
	for {
		n, err := conn.Read(buf)
		if n != 0 {
			data <- buf[:n]
		}

		if err != nil {
			data <- []byte("closed")
			log.Println("Error when reading from connection", err)
			return
		}
	}
}
