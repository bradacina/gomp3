package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/veandco/go-sdl2/sdl_mixer"
)

const (
	getVolume = -1
)

func init() {
	err := mix.Init(mix.INIT_MP3)
	if err != nil {
		log.Fatal(err)
	}

	err = mix.OpenAudio(mix.DEFAULT_FREQUENCY,
		mix.DEFAULT_FORMAT,
		mix.DEFAULT_CHANNELS,
		mix.DEFAULT_CHUNKSIZE)
	if err != nil {
		log.Fatal(err)
	}
}

var ErrNoFileLoaded = errors.New("player: No file was loaded")

type Player struct {
	folder  string
	music   *mix.Music
	loaded  bool
	looping bool
	file    string
	lock    sync.RWMutex
}

type PlayerStatus struct {
	Folder   string
	Loaded   bool
	Looping  bool
	IsPaused bool
	Filename string
	Volume   int
}

func NewPlayer(folder string) *Player {
	return &Player{folder: folder}
}

func (p *Player) Status() *PlayerStatus {
	s := &PlayerStatus{}

	p.lock.RLock()
	defer p.lock.RUnlock()

	s.Folder = p.folder
	s.Loaded = p.loaded
	s.Looping = p.looping
	s.Filename = p.file

	s.IsPaused = mix.PausedMusic()
	s.Volume = mix.VolumeMusic(getVolume)

	return s
}

func (p *Player) ChangeVolume(delta int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	vol := mix.VolumeMusic(getVolume)

	vol = clamp(vol + delta)

	mix.VolumeMusic(vol)
}

func clamp(vol int) int {
	if vol < 0 {
		return 0
	}

	if vol > mix.MAX_VOLUME {
		return mix.MAX_VOLUME
	}

	return vol
}

func (p *Player) TogglePause() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.music == nil {
		return
	}

	if p.loaded {
		if !mix.PausedMusic() {
			log.Print("music is playing")
			mix.PauseMusic()
		} else {
			log.Print("music is not playing")
			mix.ResumeMusic()
		}
	} else {
		p.Play()
	}
}

func (p *Player) Play() error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.music == nil {
		return ErrNoFileLoaded
	}

	loops := 1
	if p.looping {
		loops = -1
	}

	err := p.music.Play(loops)
	if err == nil {
		p.loaded = true
	}
	return err
}

func (p *Player) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()

	mix.HaltMusic()
	if p.music != nil {
		p.music.Free()
		p.music = nil
		p.file = ""
		p.loaded = false
	}
}

func (p *Player) Close() {
	mix.CloseAudio()
}

func (p *Player) LoadSong(file string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.loaded = false
	music, err := mix.LoadMUS(filepath.Join(p.folder, file))
	if err == nil {
		p.music = music
		p.file = file
	}
	return err
}

func (p *Player) ToggleLooping() {
	p.looping = !p.looping
}

func (p *Player) ListSongs() ([]string, error) {
	file, err := os.Open(p.folder)
	if err != nil {
		return nil, err
	}

	files, err := file.Readdir(-1)
	var filesOnly []string

	for _, item := range files {
		if !item.IsDir() && filepath.Ext(item.Name()) == ".mp3" {
			filesOnly = append(filesOnly, item.Name())
		}
	}

	return filesOnly, nil
}
