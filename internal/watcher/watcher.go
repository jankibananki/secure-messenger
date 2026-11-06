package watcher

import (
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Start begins monitoring a directory for newly created files.
func Start(
	directory string,
	onFile func(string),
	onError func(error),
) (func(), error) {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.Add(directory); err != nil {
		w.Close()
		return nil, err
	}

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}

				if !event.Has(fsnotify.Create) {
					continue
				}

				info, err := os.Lstat(event.Name)
				if err != nil || !info.Mode().IsRegular() {
					continue
				}

				if onFile != nil {
					onFile(event.Name)
				}

			case err, ok := <-w.Errors:
				if !ok {
					return
				}

				if onError != nil {
					onError(err)
				}
			}
		}
	}()

	var once sync.Once

	stop := func() {
		once.Do(func() {
			w.Close()
		})
	}

	return stop, nil
}
