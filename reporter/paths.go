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

// DefaultLogPath returns the path to Hearthstone's logging output file.
//
// It returns the expected file path, assuming a standard game installation.
func DefaultLogFile() string {
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
