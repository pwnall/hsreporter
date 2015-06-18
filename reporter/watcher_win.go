// +build windows

package reporter

import (
  "os"
  "time"
)

// listenLoop repeatedly listens for filesystem events and acts on them.
func (l *LogWatcher) listenLoop() {
  // Hearthstone doesn't explicitly flush its logging output file. That seems
  // to make it impossible to detect changes to the file in a timely manner on
  // Windows. Therefore, on Windows, we are forced to degrade to polling.
  //
  // This problem is confirmed by another Hearthstone tracker project.
  // https://github.com/stevschmid/track-o-bot/blob/master/src/HearthstoneLogWatcher.cpp
  pollingTicks := time.NewTicker(time.Millisecond * 500).C

  for {
    select {
    //case <- l.fsWatcher.Events:
    //  if err := l.handleWrite(); err != nil {
    //    l.errors <- err
    //  }
    case <- pollingTicks:
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

  err = l.tailLog()

  // On Windows, Hearthstone can't truncate its log file if we keep it open.
  // Therefore, we must open and close it on every operation.
  if l.log != nil {
    l.log.Close()
    l.log = nil
  }

  return err
}
