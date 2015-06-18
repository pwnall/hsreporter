package reporter

import (
  "fmt"
  "os"
  "path"
)

// WriteConfigFile overwrites Hearthstone's logging config file.
//
// It returns any error encountered.
// The file receives the configuration necessary for the uploader.
func WriteConfigFile(configFile string, logCategories []string) error {
  err := os.MkdirAll(path.Dir(configFile), 0755)
  if err != nil {
    return err
  }
  file, err := os.OpenFile(configFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
      0644)
  if err != nil {
    return err
  }
  defer file.Close()
  for _, category := range logCategories {
    _, err := fmt.Fprintf(file,
        "[%s]\nLogLevel=1\nFilePrinting=false\n" +
        "ConsolePrinting=true\nScreenPrinting=true\n",
        category)
    if err != nil {
      return err
    }
  }
  return nil
}

// TouchLogFile opens Hearthstone's logging file.
//
// It returns any error encountered.
// Hearthstone never re-creates the log file if it is already there. Touching
// the file ensures that it will always be there when the watcher looks for it,
// which greatly reduces the watcher's complexity.
func TouchLogFile(logFile string) error {
  err := os.MkdirAll(path.Dir(logFile), 0755)
  if err != nil {
    return err
  }
  file, err := os.OpenFile(logFile, os.O_WRONLY | os.O_CREATE, 0644)
  if err != nil {
    return err
  }
  return file.Close()
}
