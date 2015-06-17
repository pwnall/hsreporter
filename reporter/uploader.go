package reporter

import (
  "bytes"
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strings"
)

// The JSON response returned by a GET request to the HTTP endpoint.
type ServerConfig struct {
  Categories []string
}

// The logic for uploading logging output to a HTTP endpoint.
type Uploader struct {
  // The logging configuration requested by the HTTP endpoint.
  ServerConfig ServerConfig

  // The HTTP endpoint's URL.
  url string
  // The Authorization HTTP header value.
  authHeader string
  // The X-HsReport-Sess HTTP header value.
  sessHeader string
  // Source for Hearthstone's logging output.
  logLines <-chan []byte
  // Sink for HTTP errors.
  errors chan error
  // http.Client instance used for all communication with the HTTP endpoint.
  httpClient http.Client
}

// Init sets up the uploader's initial state.
func (u *Uploader) Init(serverUrl string, serverToken string,
    logLines <-chan []byte) {
  u.logLines = logLines
  u.url = serverUrl
  u.authHeader = "Token " + serverToken
  u.errors = make(chan error, 5)
}

// Errors returns a channel that receives upload errors.
func(u *Uploader) Errors() <-chan error {
  return u.errors
}

// FetchConfig obtains logging configuration data from the server.
//
// It returns any error encountered.
func (u *Uploader) FetchConfig() error {
  // NOTE: It'd be more natural to generate the upload session ID in Init().
  //       However, that'd require having Init() return an error. Generating a
  //       new session ID whenever we get the logging configuration seems
  //       reasonable enough.
  sessionBytes := make([]byte, 16)
  if _, err := rand.Read(sessionBytes); err != nil {
    return err
  }
  u.sessHeader = strings.TrimRight(
      base64.URLEncoding.EncodeToString(sessionBytes), "=")

  request, err := http.NewRequest("GET", u.url, nil)
  if err != nil {
    return err
  }
  request.Header.Add("Authorization", u.authHeader)
  request.Header.Add("X-HsReport-Sess", u.sessHeader)

  response, err := u.httpClient.Do(request)
  if err != nil {
    return err
  }

  jsonBytes, err := ioutil.ReadAll(response.Body)
  response.Body.Close()
  if err != nil {
    return err
  }

  if err = json.Unmarshal(jsonBytes, &u.ServerConfig); err != nil {
    return err
  }

  return nil
}

// Start starts uploading Hearthstone logging information to the HTTP endpoint.
func (u *Uploader) Start() error {
  go u.uploadLoop()
  return nil
}

// uploadLoop reads Hearthstone's logging output and posts it to the server.
func (u *Uploader) uploadLoop() {
  buffer := bytes.Buffer{}
  for {
    firstLine := <- u.logLines
    buffer.Write(firstLine)

    // Batch all available lines in the same request.
    batchingLoop: for {
      select {
      case line := <- u.logLines:
        buffer.Write(line)
      default:
        break batchingLoop
      }
    }

    request, err := http.NewRequest("POST", u.url, &buffer)
    if err != nil {
      u.errors <- err
      continue
    }
    request.Header.Add("Authorization", u.authHeader)
    request.Header.Add("Content-Type", "application/octet-stream")
    request.Header.Add("X-HsReport-Sess", u.sessHeader)

    postSucceeded := false
    for attemptsLeft := 3; attemptsLeft > 0; attemptsLeft -= 1 {
      response, err := u.httpClient.Do(request)
      if err == nil {
        response.Body.Close()
        postSucceeded = true
        break
      }
      u.errors <- err
    }
    if postSucceeded {
      buffer.Reset()
    }
  }
}
