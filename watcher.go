package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"slices"
	"strings"
)

func Watch(p string, reload func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	htmlRef := HtmlCrossReference(p)
	filesToWatch := []string{p}
	for _, s := range htmlRef.scripts {
		filesToWatch = append(filesToWatch, s)
	}
	for _, s := range htmlRef.styles {
		filesToWatch = append(filesToWatch, s)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					// tmp file saved
					if !strings.HasSuffix(event.Name, "~") {
						log.Println("modified file:", event.Name)
					}
					if slices.Contains(filesToWatch, event.Name) {
						log.Println("reload:", event.Name)
						reload()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, pp := range filesToWatch {
		watcher.Add(path.Dir(pp))
		// TODO: Stop watching when reference changed
		log.Println("start watching: ", pp)
	}
	<-make(chan struct{})
}
