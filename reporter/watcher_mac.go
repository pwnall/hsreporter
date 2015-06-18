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

// reportLine sends the line information over the channel.
func (l *LogWatcher) reportLine(line []byte) {
  if len(line) == 0 || line[0] != byte('[') {
    // Skip lines that don't start with a [
    return
  }

  // NOTE: We copy the slice because its underlying buffer is the line buffer,
  //       which changes often.
  // TODO(pwnall): Consider cutting slices from large pools.
  lineCopy := make([]byte, len(line))
  copy(lineCopy, line)
  l.logLines <- lineCopy
}
