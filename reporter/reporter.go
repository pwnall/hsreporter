// Package reporter uploads Hearthstone's logging output to a HTTP endpoint.
//
// It includes logic for configuring the game's logging mechanism, watching the
// logging output file, and uploading the logging output in real time, as it is
// written to the file.
package reporter

// Configuration for the log uploader.
type Config struct {
  // Path to Hearthstone's logging config file.
  ConfigFile string
  // Path to Hearthstone's logging output file.
  LogFile string
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
  // Log watcher.
  Watcher LogWatcher
}

// Sets up the logger's state.
//
// It returns any error encountered.
// The caller must have set up the logger's configuration.
func (s *State) Init() error {
  if err := s.Watcher.Init(s.Config.LogFile); err != nil {
    return err
  }

  s.Uploader.Init(s.Config.ServerUrl, s.Config.ServerToken,
      s.Watcher.LogLines())
  if err := s.Uploader.FetchConfig(); err != nil {
    return err
  }

  return nil
}

// Writes Hearthstone's log configuration and touches its log file.
func (s *State) ConfigLogging() error {
  if err := WriteConfigFile(s.Config.ConfigFile,
      s.Uploader.ServerConfig.Categories); err != nil {
    return err
  }
  if err := TouchLogFile(s.Config.LogFile); err != nil {
    return err
  }
  return nil
}
