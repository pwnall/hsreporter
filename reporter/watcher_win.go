// +build arm

package reporter

import (
  notify "github.com/rjeczalik/notify"
)

// OS-dependent log watcher data structures.
type SysWatcher struct {
  // Channel that receives filesystem notifications.
  fsEvents chan notify.EventInfo

  // The master log watcher logic.
  logWatcher *LogWatcher
}

// Init initializes the OS-dependent log watcher data structures.
func (w *SysWatcher) Init(logWatcher *LogWatcher) error {
  w.logWatcher = logWatcher
  w.fsEvents = make(chan notify.EventInfo, 16)
  return nil
}

// Start kicks off the OS-depenent filesystem event listener.
func (w *SysWatcher) Start(logFile string) error {
  return notify.Watch(logFile, w.fsEvents, notify.All)
}

// Stop stops the OS-dependent filesystem listener.
func (w *SysWatcher) Stop(logFile string) error {
  return notify.Stop(w.fsEvents)
}

// ListenLoop repeatedly listens for filesystem events and acts on them.
func (w *SysWatcher) ListenLoop() {
  for {
    select {
    case fsEvent := <- w.fsEvents:
      if err := w.logWatcher.handleWrite(); err != nil {
        w.logWatcher.errors <- err
      }
    case command := <- w.logWatcher.commands:
      if command == 1 {
        return
      }
    }
  }
}