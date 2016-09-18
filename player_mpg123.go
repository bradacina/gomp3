// +build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/kr/pty"
)

type Player struct {
	pty     *os.File
	volume  int
	mutex   *sync.RWMutex
	closed  bool
	paused  bool
	stopped bool
}

type PlayerStatus struct {
	Volume  int
	Closed  bool
	Paused  bool
	Stopped bool
	// todo: play progress bar
	// todo: current idv3 tag of song
}

var player *Player

func init() {
	player = newPlayer()
}

func ThePlayer() *Player {
	return player
}

func newPlayer() *Player {
	log.Println("Created mpg123 instance")
	cmd := exec.Command("mpg123", "-R", "--keep-open")
	p, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}

	return &Player{
		pty:    p,
		mutex:  &sync.RWMutex{},
		volume: 0,
		closed: false}
}

func (p *Player) Close() {
	if p.closed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pty.Write([]byte("q\r\n"))
	p.closed = true
	p.stopped = true
	p.paused = false
	p.volume = 0
}

func (p *Player) setVolume() {
	p.pty.Write([]byte("v " + stringify(p.volume) + "\r\n"))
}

func (p *Player) ChangeVolume(delta int) {
	if p.closed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.volume += delta
	p.setVolume()
}

func (p *Player) Status() *PlayerStatus {

	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return &PlayerStatus{
		Volume:  p.volume,
		Closed:  p.closed,
		Paused:  p.paused,
		Stopped: p.stopped}
}

func (p *Player) PlaySong(file string) {
	if p.closed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.pty.Write([]byte("l " + file + "\r\n"))
}

func (p *Player) Pause() {
	if p.closed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.pty.Write([]byte("p\r\n"))
}

func (p *Player) Stop() {
	if p.closed {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.pty.Write([]byte("s\r\n"))
}

func stringify(v int) string {
	return strconv.FormatInt(int64(v), 10)
}
