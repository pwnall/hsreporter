package reporter

import (
  "bytes"
  "fmt"
  "os"
)

type LogWatcher struct {
  // OS-dependent functionality.
  sys SysWatcher
  // Path to the log file that will be watched.
  logFile string
  // Sink for the lines written to the log file.
  logLines chan []byte
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
}

// Init sets up the filesystem watcher.
func (l *LogWatcher) Init(logFile string) error {
  l.logFile = logFile
  l.logLines = make(chan []byte, 1024)
  l.readOffset = -1
  l.lineBuffer = make([]byte, 4096)[:0]
  l.errors = make(chan error, 5)
  return l.sys.Init(l)
}

// LogLines returns the channel that produces Hearthstone's logging output.
func (l *LogWatcher) LogLines() <-chan []byte {
  return l.logLines
}

// Errors returns the channel for errors encountered while watching the log.
func (l *LogWatcher) Errors() <-chan error {
  return l.errors
}

// Start spawns a goroutine that listens for log-related filesystem events.
func (l *LogWatcher) Start() error {
  if err := l.handleWrite(); err != nil {
    return err
  }
  if err := l.sys.Start(l.logFile); err != nil {
    return err
  }
  go l.sys.ListenLoop()
  return nil
}

// Stop causes the filesystem listener to break out of its loop.
func (l *LogWatcher) Stop() error {
  if err := l.sys.Stop(l.logFile); err != nil {
    return err
  }
  l.commands <- 1
  return nil
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

  fmt.Printf("Reading from %d to %d\n", l.readOffset, logSize)

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
