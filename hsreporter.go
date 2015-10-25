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
  flag.StringVar(&logger.Config.GameLogFile, "game-log-file",
      reporter.DefaultGameLogFile(),
      "Path to Hearthstone's game logging output file")
  flag.StringVar(&logger.Config.NetLogFile, "net-log-file",
      reporter.DefaultNetLogFile(),
      "Path to Hearthstone's network logging output file")
  flag.Parse()

  if err := logger.Init(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  fmt.Printf("Logging config: %s\n", logger.Config.ConfigFile)
  fmt.Printf("Game log: %s\n", logger.Config.GameLogFile)
  fmt.Printf("Network log: %s\n", logger.Config.NetLogFile)
  fmt.Printf("Logging categories: %v\n",
      logger.Uploader.ServerConfig.Categories)
  fmt.Printf("Uploading old log data: %v\n",
      logger.Uploader.ServerConfig.ExistingData)

  if err := logger.ConfigLogging(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  // NOTE: The uploader must start before any watcher, so it can drain the
  //       watchers' output channel. Otherwise, a watcher can deadlock in
  //       Start() while producing old log data.
  if err := logger.Uploader.Start(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  // NOTE: The network log watcher must start first, so that we uploads region
  //       information to the server before we start uploading game chunks.
  if err := logger.NetLogWatcher.Start(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if logger.Uploader.ServerConfig.ExistingData {
    logger.GameLogWatcher.ReportExistingData()
  }
  if err := logger.GameLogWatcher.Start(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  uploadErrors := logger.Uploader.Errors()
  gameLogWatchErrors := logger.GameLogWatcher.Errors()
  netLogWatchErrors := logger.NetLogWatcher.Errors()
  for {
    select {
    case uploadErr := <- uploadErrors:
      fmt.Printf("Upload error: %v\n", uploadErr)
    case watchErr := <- gameLogWatchErrors:
      fmt.Printf("Game log watch error: %v\n", watchErr)
    case watchErr := <- netLogWatchErrors:
      fmt.Printf("Net log watch error: %v\n", watchErr)
    }
  }
}
