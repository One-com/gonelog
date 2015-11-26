# gonelog
Golang logging library [![GoDoc](https://godoc.org/github.com/one-com/gonelog/log?status.svg)](https://godoc.org/github.com/one-com/gonelog/log)

Package gonelog/log is a drop-in replacement for the standard Go logging library "log" which is fully source code compatible support all the standard library API while at the same time offering advanced logging features through an extended API.

The design goals of gonelog was:

    * Standard library source level compatibility with mostly preserved behaviour.
    * Leveled logging with syslog levels.
    * Structured key/value logging
    * Hierarchical contextable logging to have k/v data in context logged automatically.
    * Low resource usage to allow more (debug) log-statements even if they don't result in output.
    * Light syntax to encourage logging on INFO/DEBUG level. (and low cost of doing so)
    * Explore compatiblity with http://golang.org/x/net/context to make the context object a Logger.
    * Flexibility in how log events are output.
    * A fast simple lightweight default in systemd newdaemon style only outputting <level>message
      to standard output.

See the examples in api_test.go

