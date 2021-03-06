package reporter

import (
  "bytes"
  "os"
  fsnotify "gopkg.in/fsnotify.v1"
)

type LogWatcher struct {
  // Path to the log file that will be watched.
  logFile string
  // Filesystem notifications client.
  fsWatcher *fsnotify.Watcher
  // Sink for the lines written to the log file.
  logLines chan<- []byte
  // Tells the log watch loop when to stop.
  commands chan int
  // Sink for errors encountered by the log file watching loop.
  errors chan error
  // Filesystem handle for the log file.
  log *os.File
  // The number of bytes already read from the log file.
  readOffset int64
  // The buffer used to read from the file.
  lineBuffer []byte
  // True if lines that don't start with [ should be discarded.
  filterLines bool
}

// Init sets up the filesystem watcher.
func (l *LogWatcher) Init(logFile string, filterLines bool,
    logLines chan<- []byte) error {
  l.logFile = logFile
  l.filterLines = filterLines
  l.logLines = logLines

  var err error
  l.fsWatcher, err = fsnotify.NewWatcher()
  if err != nil {
    return err
  }

  l.readOffset = -1
  l.lineBuffer = make([]byte, 4096)[:0]
  l.errors = make(chan error, 5)
  return nil
}

// Errors returns the channel for errors encountered while watching the log.
func (l *LogWatcher) Errors() <-chan error {
  return l.errors
}

// ReportExistingData configures the watcher to dump the initial file contents.
//
// By default, a log watcher only reports data written to the watched file
// after Start is called.
func (l *LogWatcher) ReportExistingData() {
  if l.readOffset == -1 {
    l.readOffset = 0
  }
}

// Start spawns a goroutine that listens for log-related filesystem events.
func (l *LogWatcher) Start() error {
  if err := l.handleWrite(); err != nil {
    return err
  }
  if err := l.fsWatcher.Add(l.logFile); err != nil {
    return err
  }
  go l.listenLoop()
  return nil
}

// Stop causes the filesystem listener to break out of its loop.
func (l *LogWatcher) Stop() error {
  if err := l.fsWatcher.Remove(l.logFile); err != nil {
    return err
  }
  l.commands <- 1
  return nil
}

// tailLog reads the newly appended data from the log file.
func (l *LogWatcher) tailLog() error {
  fileInfo, err := l.log.Stat()
  if err != nil {
    return err
  }

  logSize := fileInfo.Size()
  if logSize < l.readOffset {
    // The log file was truncated.
    l.readOffset = 0
  } else if l.readOffset == -1 {
    // The watcher is just getting started.
    l.readOffset = logSize
  }

  for l.readOffset < logSize {
    readSize := logSize - l.readOffset
    bufferOffset := len(l.lineBuffer)
    bufferCapacity := cap(l.lineBuffer) - bufferOffset
    if readSize > int64(bufferCapacity) {
      readSize = int64(bufferCapacity)
    }

    readBuffer := l.lineBuffer[bufferOffset : bufferOffset + int(readSize)]
    bytesRead, err := l.log.ReadAt(readBuffer, l.readOffset)
    if err != nil {
      return err
    }
    l.readOffset += int64(bytesRead)
    l.lineBuffer = l.lineBuffer[0 : bufferOffset + bytesRead]

    l.sliceLines(bufferOffset)
  }
  return nil
}

// sliceLines removes complete lines from the read buffer.
// "bufferOffset
func (l *LogWatcher) sliceLines(bufferOffset int) {
  lineStart := 0
  for {
    readBuffer := l.lineBuffer[bufferOffset:]
    relativeIndex := bytes.IndexByte(readBuffer, byte('\n'))
    if relativeIndex == -1 {
      break
    }
    newlineIndex := relativeIndex + bufferOffset
    l.reportLine(l.lineBuffer[lineStart : newlineIndex + 1])

    bufferOffset = newlineIndex + 1
    lineStart = bufferOffset
  }

  if lineStart == 0 {
    return
  }
  bufferOffset = len(l.lineBuffer) - lineStart
  copy(l.lineBuffer[0:bufferOffset], l.lineBuffer[lineStart:])
  l.lineBuffer = l.lineBuffer[0:bufferOffset]
}
