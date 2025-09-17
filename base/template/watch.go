package template

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

// ||------------------------------------------------------------------------------------------------||
// || WatchTemplates: Watch template directory for changes
// ||------------------------------------------------------------------------------------------------||

func WatchTemplates(dir string, reload func()) error {
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
					log.Printf("[template] change detected: %s, reloading...", event.Name)
					reload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("[template] watcher error: %v", err)
			}
		}
	}()

	return watcher.Add(dir)
}
