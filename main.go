package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"runtime/pprof"
)

var logFile *os.File

func closeLog() {
	logFile.Close()
}

func initLog(folder string) {
	err := os.MkdirAll(folder, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags)
	logPath := path.Join(folder, "gomp3.log")
	logFile, err = os.Create(logPath)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
}

func watchForInterrupts() {
	interrupt := make(chan os.Signal, 100)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	for _ = range interrupt {
		pprof.StopCPUProfile()
		closeLog()
		os.Exit(1)
	}
}

func main() {
	c := getConfig()
	initLog(c.LogPath)
	defer closeLog()

	profPath := path.Join(c.LogPath, "gomp3.prof")
	f, err := os.Create(profPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go watchForInterrupts()

	p := NewPlayer(c.Mp3Location)

	StartWeb(p, c.WebappPath, c.Port)
}
