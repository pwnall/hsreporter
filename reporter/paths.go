package reporter

import (
  "os"
  "path"
  "runtime"
)

// DefaultConfigFile returns the path to Hearthstone's logging config file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultConfigFile() string {
  var configRoot string
  if runtime.GOOS == "windows" {
    configRoot = os.Getenv("LOCALAPPDATA")
  } else {
    // OSX.
    configRoot = path.Join(os.Getenv("HOME"), "Library", "Preferences")
  }
  return path.Join(configRoot, "Blizzard", "Hearthstone", "log.config")
}

// DefaultLogPath returns the path to Hearthstone's logging output file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultLogFile() string {
  if runtime.GOOS == "windows" {
    var programFiles string
    if runtime.GOARCH == "amd64" {
      programFiles = path.Join("C:", "Program Files (x86)")
    } else {
      programFiles = path.Join("C:", "Program Files")
    }
    return path.Join(programFiles, "Hearthstone", "Hearthstone_data",
        "output_log.txt")
  }

  // OSX.
  return path.Join(os.Getenv("HOME"), "Library", "Logs", "Unity", "Player.log")
}
