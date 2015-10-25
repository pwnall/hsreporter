# hsreporter

hsreporter is a [Go](https://golang.org/) command-line tool and library that
uploads [Hearthstone](http://us.battle.net/hearthstone/)'s logging output to a
Web server.


## Usage

Obtain your API token from the Hearthstone analytics application that you use,
and pass it to run hsreporter. Your analytics Web application should give you a
command line for hsreporter that contains all necessary information. The
command should look as follows.

```bash
hsreporter -token xxxxxxxxxx -server https://my.tracker.com/hsreporter.json
```

For best results, restart Hearthstone after starting hstracker. A restart is
absolutely required the first time you run the tool, so Hearthstone can pick up
configuration changes.

hstracker must run for the entire duration of a game. Stopping and restarting
hsreporter during a game will render that game's report invalid.

To avoid interference, do not run other Hearthstone tracking software at the
same time. The following trackers are known to interfere with hsreporter.

* [Track-o-Bot](https://trackobot.com/)
* [HearthStats](http://hearthstats.net/)
* [Hearthstone Tracker](http://hearthstonetracker.com/)


## Protocol

hsreporter only sends HTTP requests to the provided server URL. No URL
derivation is performed. The server must be able to process `GET` and `POST`
requests at the same URL.

hsreporter requests include the following headers:

* `Authorization` is set to `Token xxxxxxxxxxx`, following an old version of
  the [Bearer Token Usage](https://tools.ietf.org/html/rfc6750) RFC, namely
  [Token Access Authentication](https://tools.ietf.org/html/draft-hammer-http-token-auth-01).
  This is likely to change when Ruby on Rails 5 is released.
* `X-HsReport-Id` consists of a random nonce and a sequence number, separated
   by a space.
    * The random nonce that gets reset every time the program starts. It is
      provided to help the Web server detect situations where the user restarts
      hsreporter, causing it to miss logging output. The nonce may only
      contain the characters used in
      [URL-safe base64 encoding](https://tools.ietf.org/html/rfc4648#section-5)
      without padding.
    * The sequence number starts at 0 and is incremented on every HTTP request.
      It is provided to help the Web server detect situations where an
      intermediary HTTP proxy, such as a load balancer, drops hsreporter's POST
      requests, causing the server to miss logging output.
* `X-HsReport-Proto` is only sent for GET requests, and contains the protocol
  version implemented by the reporter. The server SHOULD reject clients whose
  protocol version does not match the version it understands. The server SHOULD
  include an error message that describes the problem.

When starting, hsreporter will send an HTTP `GET` request to the provided
server URL. The server must produce a JSON response containing the Hearthstone
game and network logging categories that should be uploaded. This gives the
server a way to reduce the upload bandwidth, by only asking for information
that it knows how to parse.

```json
{
  "categories": ["Power", "Zone"]
}
```

If the server can handle de-duplicating old game reports, it SHOULD ask the
reporter to upload the game log data that exists when the log is started.

```json
{
  "categories": ["Power", "Zone"],
  "existingData": true
}
```

The server should verify the supplied token, and produce an error if it is
invalid. hsreporter will immediately stop and report the error, while the user
is still paying attention to its window.

```json
{
  "error": "Invalid token"
}
```

Once initialized, hsreporter watches Hearthstone's log file, and uploads log
lines that match the requested categories. hsreporter sends `POST` requests to
the provided server URL. The log data is contained in the POST request body,
with a `Content-Type` of `application/octet-stream`. A single request may
contain multiple log lines separated by the LF (`"\n"`) character.


## Copyright and Licensing

This package is (C) Victor Costan 2015, and made available under the MIT
license, which is contained in the `LICENSE` file.
