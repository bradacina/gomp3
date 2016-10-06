package main

import (
	"log"
	"os"
	"path"
)

func initLog(folder string) {
	err := os.MkdirAll(folder, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)
	logPath := path.Join(folder, "gomp3.log")
	f, err := os.Create(logPath)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
}

func main() {
	c := getConfig()
	initLog(c.LogPath)
	p := NewPlayer(c.Mp3Location)

	StartWeb(p, c.WebappPath, c.Port)
}
