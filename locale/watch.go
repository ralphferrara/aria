package locale

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

//||------------------------------------------------------------------------------------------------||
//|| WatchRendered: Watch `.rendered` directory for changes and reload
//||------------------------------------------------------------------------------------------------||

func WatchRendered(dir string) error {
	renderedDir := filepath.Join(dir, ".rendered")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					log.Printf("[locale] change detected: %s, reloading...", event.Name)
					if err := LoadRendered(dir); err != nil {
						log.Printf("[locale] reload failed: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[locale] watcher error: %v", err)
			}
		}
	}()

	return watcher.Add(renderedDir)
}
