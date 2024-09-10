package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// Watch html file and its references change
func Watch(htmlFilePath string, reload func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	htmlRef, err := HtmlCrossReference(htmlFilePath)
	if err != nil {
		log.Println("Format file path failed")
		return err
	}

	h, err := filepath.Abs(htmlFilePath)
	if err != nil {
		return err
	}
	filesToWatch := []string{h}
	for _, s := range append(htmlRef.scripts, htmlRef.styles...) {
		p, err := filepath.Abs(s)
		if err != nil {
			log.Printf("Format file path failed %s\n", s)
			return err
		}
		filesToWatch = append(filesToWatch, p)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					// tmp file saved
					if !strings.HasSuffix(event.Name, "~") {
						log.Println("modified file:", event.Name)
					}
					p, _ := filepath.Abs(event.Name)
					if slices.Contains(filesToWatch, p) {
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
	return nil
}
