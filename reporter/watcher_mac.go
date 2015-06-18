// +build !windows

package reporter

import (
  "os"
)

// listenLoop repeatedly listens for filesystem events and acts on them.
func (l *LogWatcher) listenLoop() {
  for {
    select {
    case <- l.fsWatcher.Events:
      if err := l.handleWrite(); err != nil {
        l.errors <- err
      }
    case fsError := <- l.fsWatcher.Errors:
      l.errors <- fsError
    case command := <- l.commands:
      if command == 1 {
        break
      }
    }
  }
}

// handleWrite is called when the log file is updated.
func (l *LogWatcher) handleWrite() error {
  var err error
  if l.log == nil {
    l.log, err = os.OpenFile(l.logFile, os.O_RDONLY, 0644)
    if err != nil {
      return err
    }
  }

  return l.tailLog()
}
