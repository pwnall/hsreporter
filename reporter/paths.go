package reporter

import (
  "os"
  "path/filepath"
)

// DefaultConfigFile returns the path to Hearthstone's logging config file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultConfigFile() string {
  // Windows attempt.
  localAppData := os.Getenv("LOCALAPPDATA")
  _, err := os.Stat(filepath.Join(localAppData, "Blizzard", "Hearthstone"))
  if localAppData != "" && err == nil {
    return filepath.Join(localAppData, "Blizzard", "Hearthstone", "log.config")
  }

  // OSX attempt.
  homeDir := os.Getenv("HOME")
  configRoot := filepath.Join(homeDir, "Library", "Preferences")
  _, err = os.Stat(filepath.Join(configRoot, "Blizzard", "Hearthstone"))
  if homeDir != "" && err == nil  {
    return filepath.Join(configRoot, "Blizzard", "Hearthstone", "log.config")
  }

  // Failed to find the default path.
  return ""
}

// DefaultGameLogFile returns the path to Hearthstone's game logging file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultGameLogFile() string {
  // Windows attempts.
  for _, programDir := range []string{"Program Files (x86)", "Program Files"} {
    dataDir := filepath.Join("C:", programDir, "Hearthstone",
        "Hearthstone_data")
    if _, err := os.Stat(dataDir); err == nil {
      return filepath.Join(dataDir, "output_log.txt")
    }
  }

  // OSX attempt.
  homeDir := os.Getenv("HOME")
  _, err := os.Stat(filepath.Join(homeDir, "Library"))
  if homeDir != "" && err == nil {
    return filepath.Join(homeDir, "Library", "Logs", "Unity", "Player.log")
  }

  // Failed to find a default path.
  return ""
}

// DefaultNetLogFile returns the path to Hearthstone's network logging file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultNetLogFile() string {
  // Windows attempts.
  for _, programDir := range []string{"Program Files (x86)", "Program Files"} {
    dataDir := filepath.Join("C:", programDir, "Hearthstone")
    if _, err := os.Stat(dataDir); err == nil {
      return filepath.Join(dataDir, "ConnectLog.txt")
    }
  }

  // OSX attempt.
  appDir := "/Applications/Hearthstone"
  if _, err := os.Stat(appDir); err == nil {
    return filepath.Join(appDir, "ConnectLog.txt")
  }

  // Failed to find a default path.
  return ""
}
