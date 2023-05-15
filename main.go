package main

import (
	"flag"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var (
	fileargs = flag.String("files", "", "csv of files to watch")
	helpargs = flag.String("help", "./gofim -files /etc/passwd", "Specify how to run gofim")
)

func watcher(fPath string) {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	// Start listening for events.
	// closure style!
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Not okay event!", event)
					return
				}
				log.Println("event:", event)

				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)

				} else if event.Has(fsnotify.Create) {
					log.Println("Created file:", event.Name)

				} else if event.Has(fsnotify.Remove) {
					log.Println("Removed file:", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Println("Adding " + fPath + "to watch for...")
	err = watcher.Add(fPath)
	if err != nil {
		panic(err)
	}

	<-make(chan struct{})
}

func main() {

	flag.Parse()
	for _, f := range strings.Split(*fileargs, ",") {
		watcher(f)
	}
}
