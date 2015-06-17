package reporter

import (
  "os"
  "path/filepath"
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
    configRoot = filepath.Join(os.Getenv("HOME"), "Library", "Preferences")
  }
  return filepath.Join(configRoot, "Blizzard", "Hearthstone", "log.config")
}

// DefaultLogPath returns the path to Hearthstone's logging output file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultLogFile() string {
  if runtime.GOOS == "windows" {
    var programFiles string
    if runtime.GOARCH == "amd64" {
      programFiles = filepath.Join("C:", "Program Files (x86)")
    } else {
      programFiles = filepath.Join("C:", "Program Files")
    }
    return filepath.Join(programFiles, "Hearthstone", "Hearthstone_data",
        "output_log.txt")
  }

  // OSX.
  return filepath.Join(os.Getenv("HOME"), "Library", "Logs", "Unity",
      "Player.log")
}
