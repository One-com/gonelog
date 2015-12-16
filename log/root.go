package log

import (
	"fmt"
	ilog "github.com/One-com/gonelog"
	"github.com/One-com/gonelog/syslog"
	"io"
	"os"
)

// All the toplevel package functionality

// The default log context
var defaultLogger *Logger

func Default() *Logger {
	return defaultLogger
}

func init() {
	// Default Logger is an ordinary stdlib like logger, to be compatible
	defaultLogger = New(os.Stderr, "", LstdFlags)
}

// Sets the default logger to the minimal mode, where it doesn't log timestamps
// But only emits systemd/syslog-compatible "<level>message" lines.
func Minimal() {
	minHandler := NewMinFormatter(SyncWriter(os.Stdout))
	defaultLogger.SetHandler(minHandler)
	// turn of doing timestamps *after* not using them
	defaultLogger.DoTime(false)
}

// Create a child K/V logger of the default logger
func With(kv ...interface{}) *Logger {
	return defaultLogger.With(kv...)
}

// Clone the default logger
func Clone(kv ...interface{}) *Logger {
	return defaultLogger.Clone(kv...)
}

// AutoColoring turns on coloring if the output Writer is connected to a TTY
func AutoColoring() {
	defaultLogger.AutoColoring()
}

//--- level logger stuff

// ALERT makes the default Logger create a log event at ALERT level.
func ALERT(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func CRIT(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_CRIT
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func ERROR(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ERROR
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func WARN(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_WARN
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func NOTICE(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_NOTICE
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func INFO(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_INFO
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func DEBUG(msg string, kv ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_DEBUG
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}

// Methods which return a function which will do the queries logging when called.
func ALERTok() (ilog.LogFunc, bool)  { l := defaultLogger; return l.alert, l.Does(syslog.LOG_ALERT) }
func CRITok() (ilog.LogFunc, bool)   { l := defaultLogger; return l.crit, l.Does(syslog.LOG_CRIT) }
func ERRORok() (ilog.LogFunc, bool)  { l := defaultLogger; return l.error, l.Does(syslog.LOG_ERROR) }
func WARNok() (ilog.LogFunc, bool)   { l := defaultLogger; return l.warn, l.Does(syslog.LOG_WARN) }
func NOTICEok() (ilog.LogFunc, bool) { l := defaultLogger; return l.notice, l.Does(syslog.LOG_NOTICE) }
func INFOok() (ilog.LogFunc, bool)   { l := defaultLogger; return l.info, l.Does(syslog.LOG_INFO) }
func DEBUGok() (ilog.LogFunc, bool)  { l := defaultLogger; return l.debug, l.Does(syslog.LOG_DEBUG) }

//---

func IncLevel() bool {
	return defaultLogger.IncLevel()
}

func DecLevel() bool {
	return defaultLogger.DecLevel()
}

// SetLevel set the Logger log level.
// returns success
func SetLevel(level syslog.Priority) bool {
	return defaultLogger.SetLevel(level)
}

func SetDefaultLevel(level syslog.Priority, respect bool) bool {
	return defaultLogger.SetDefaultLevel(level, respect)
}

// Level returns the default Loggers log level.
func Level() syslog.Priority {
	return defaultLogger.Level()
}

//--- std logger stuff

func Flags() int {
	return defaultLogger.Flags()
}
func Prefix() string {
	return defaultLogger.Prefix()
}

func SetFlags(flag int) {
	defaultLogger.SetFlags(flag)
}
func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}
func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func Fatal(v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprint(v...)
		c.log(l, s)
	}
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprintf(format, v...)
		c.log(l, s)
	}
	os.Exit(1)

}
func Fatalln(v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprintln(v...)
		c.log(l, s)
	}
	os.Exit(1)
}

func Panic(v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprint(v...)
		c.log(l, s)
		panic(s)
	}
}
func Panicf(format string, v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprintf(format, v...)
		c.log(l, s)
		panic(s)
	}
}
func Panicln(v ...interface{}) {
	c := defaultLogger
	l := syslog.LOG_ALERT
	if c.Does(l) {
		s := fmt.Sprintln(v...)
		c.log(l, s)
		panic(s)
	}
}

func Print(v ...interface{}) {
	c := defaultLogger
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprint(v...)
		c.log(l, s)
	}
}
func Printf(format string, v ...interface{}) {
	c := defaultLogger
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprintf(format, v...)
		c.log(l, s)
	}
}
func Println(v ...interface{}) {
	c := defaultLogger
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprintln(v...)
		c.log(l, s)
	}
}

// Log is the simplest Logger method. Provide the log level your self.
func Log(level syslog.Priority, msg string, kv ...interface{}) {
	c := defaultLogger
	if c.Does(level) {
		c.log(level, msg, kv...)
	}
}
