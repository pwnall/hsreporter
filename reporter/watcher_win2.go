// +build !arm

package reporter

import (
  "fmt"
  "path/filepath"
  fsnotify "github.com/pwnall/fsnotify"
)

// OS-dependent log watcher data structures.
type SysWatcher struct {
  // Filesystem notifications client.
  fsWatcher *fsnotify.Watcher

  // The master log watcher logic.
  logWatcher *LogWatcher
}

// Init initializes the OS-dependent log watcher data structures.
func (w *SysWatcher) Init(logWatcher *LogWatcher) error {
  w.logWatcher = logWatcher
  var err error
  w.fsWatcher, err = fsnotify.NewWatcher()
  if err != nil {
    return err
  }
  return nil
}

// Start kicks off the OS-depenent filesystem event listener.
func (w *SysWatcher) Start(logFile string) error {
  return w.fsWatcher.Add(filepath.Dir(logFile))
}

// Stop stops the OS-dependent filesystem listener.
func (w *SysWatcher) Stop(logFile string) error {
  return w.fsWatcher.Remove(filepath.Dir(logFile))
}

// ListenLoop repeatedly listens for filesystem events and acts on them.
func (w *SysWatcher) ListenLoop() {
  for {
    select {
    case fsEvent := <- w.fsWatcher.Events:
      fmt.Printf("FS event: %v\n", fsEvent)
      if err := w.logWatcher.handleWrite(); err != nil {
        w.logWatcher.errors <- err
      }
    case fsError := <- w.fsWatcher.Errors:
      w.logWatcher.errors <- fsError
    case command := <- w.logWatcher.commands:
      if command == 1 {
        return
      }
    }
  }
}

