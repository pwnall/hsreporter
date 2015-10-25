// Package reporter uploads Hearthstone's logging output to a HTTP endpoint.
//
// It includes logic for configuring the game's logging mechanism, watching the
// game and network logging output files, and uploading the logging output in
// real time, as it is written to the file.
package reporter

// Configuration for the log uploader.
type Config struct {
  // Path to Hearthstone's logging config file.
  ConfigFile string
  // Path to Hearthstone's game logging output file.
  GameLogFile string
  // Path to Hearthstone's network logging output file.
  NetLogFile string
  // HTTP endpoint that receives filtered game logging output.
  ServerUrl string
  // Token used to authenticate to the HTTP endpoint.
  ServerToken string
}

// The log uploader's state.
type State struct {
  Config Config
  // HTTP data uploader.
  Uploader Uploader
  // Game log watcher.
  GameLogWatcher LogWatcher
  // Network log watcher.
  NetLogWatcher LogWatcher
}

// Sets up the logger's state.
//
// It returns any error encountered.
// The caller must have set up the logger's configuration.
func (s *State) Init() error {
  logLines := make(chan []byte, 1024)

  // The game log has a lot of useless lines, and all the useful lines start
  // with the category marker [, so we use line filtering.
  err := s.GameLogWatcher.Init(s.Config.GameLogFile, true, logLines)
  if err != nil {
    return err
  }
  // The network log has very few lines, and the category marker [ is output
  // after the current date. Filtering would be difficult to implement, and is
  // unnecessary, so we just upload everything.
  err = s.NetLogWatcher.Init(s.Config.NetLogFile, false, logLines)
  if err != nil {
    return err
  }
  // The server always needs the full network log, because its beginning
  // contains region information.
  s.NetLogWatcher.ReportExistingData()

  s.Uploader.Init(s.Config.ServerUrl, s.Config.ServerToken, logLines)
  if err := s.Uploader.FetchConfig(); err != nil {
    return err
  }

  return nil
}

// Writes Hearthstone's log configuration and touches its log files.
func (s *State) ConfigLogging() error {
  if err := WriteConfigFile(s.Config.ConfigFile,
      s.Uploader.ServerConfig.Categories); err != nil {
    return err
  }
  if err := TouchLogFile(s.Config.GameLogFile); err != nil {
    return err
  }
  if err := TouchLogFile(s.Config.NetLogFile); err != nil {
    return err
  }
  return nil
}
