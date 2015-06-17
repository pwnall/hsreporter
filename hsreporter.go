package main

import (
  "flag"
  "fmt"
  "github.com/pwnall/hsreporter/reporter"
  "os"
)

var logger reporter.State

func main() {
  flag.StringVar(&logger.Config.ServerToken, "token",
      "", "Token for authenticating to the HTTP endpoint")
  flag.StringVar(&logger.Config.ServerUrl, "server",
      "https://histone.herokuapp.com/hsreporter.json",
      "HTTP endpoint that receives logging information")
  flag.StringVar(&logger.Config.ConfigFile, "log-config",
      reporter.DefaultConfigFile(),
      "Path to Hearthstone's logging configuration file")
  flag.StringVar(&logger.Config.LogFile, "log-file",
      reporter.DefaultLogFile(),
      "Path to Hearthstone's logging output file")
  flag.Parse()

  if err := logger.Init(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Printf("Logging config: %s\n", logger.Config.ConfigFile)
  fmt.Printf("Logging output: %s\n", logger.Config.LogFile)
  fmt.Printf("Logging categories: %v\n",
      logger.Uploader.ServerConfig.Categories)

  if err := logger.ConfigLogging(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if err := logger.Uploader.Start(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if err := logger.Watcher.Start(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  uploadErrors := logger.Uploader.Errors()
  watchErrors := logger.Watcher.Errors()
  for {
    select {
    case uploadErr := <- uploadErrors:
      fmt.Printf("Upload error: %v\n", uploadErr)
    case watchErr := <- watchErrors:
      fmt.Printf("Watch error: %v\n", watchErr)
    }
  }
}
