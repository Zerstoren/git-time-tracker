package tracker

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"git-time-tracker.com/internal"
	"github.com/fsnotify/fsnotify"
)

func Notify(wg *sync.WaitGroup, repository internal.Repository, fn func(event fsnotify.Event)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		debounce := time.NewTicker(internal.DEBOUNCE_INTERVAL)
		var lastEvent fsnotify.Event

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("error: watch event channel closed")
					return
				}

				if event.Has(fsnotify.Create) {
					st, err := os.Stat(event.Name)
					if err == nil {
						if st.IsDir() && !slices.Contains(watcher.WatchList(), event.Name) {
							watcher.Add(event.Name)
						}
					}
				}

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					for _, exclude := range repository.Exclude {
						if strings.Contains(event.Name, exclude) {
							if internal.DEBUG {
								log.Println("exclude: ", event.Name, exclude)
							}
							continue
						}
					}

					lastEvent = event
					debounce.Reset(internal.DEBOUNCE_INTERVAL)
				}

			case <-debounce.C:
				if lastEvent.Name == "" {
					continue
				}

				if internal.DEBUG {
					log.Println("debounce event: ", lastEvent.Name)
				}

				fn(lastEvent)
				lastEvent = fsnotify.Event{}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			}
		}
	}()

	err = watcher.Add(repository.Path)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir(repository.Path, func(path string, d fs.DirEntry, err error) error {
		for _, exclude := range repository.Exclude {
			if strings.Contains(path, exclude) {
				return nil
			}
		}

		if d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	wg.Wait()

	return nil
}
