// The main package of gonelog is gonelog/log. This top-level package only contains interfaces used to organize methods.
package gonelog

import (
	"io"
	"github.com/One-com/gonelog/syslog"
)

// Logger is the the main interface of gonelog
// A "Logger" makes available methods compatible with the stdlib logger and
// an extended API for leveled logging.
// Logger is implemented by *log.Logger
type Logger interface {

	// Will generate a log event with this level if the Logger log level is
	// high enough.
	// The event will have the given log message and key/value structured data.
	Log(level syslog.Priority, message string, kv ...interface{}) error

	// further interfaces
	StdLogger
	LevelLogger
}

// StdLogger is the interface used by the standard lib *log.Logger
// This is the API for actually logging stuff.
type StdLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// LogFunc is the type of the function returned by *ok() methods, which will
// log at the level queried about if called.
type LogFunc func(msg string, kv ...interface{})

// LevelLogger is the extended leveled log API
type LevelLogger interface {

	// Returns true if the Logger will generate events at this log level.
	Does(level syslog.Priority) bool

	// Shorthand functions to Log() generating events on syslog levels.
	ALERT(msg string, kv ...interface{})
	CRIT(msg string, kv ...interface{})
	ERROR(msg string, kv ...interface{})
	WARN(msg string, kv ...interface{})
	NOTICE(msg string, kv ...interface{})
	INFO(msg string, kv ...interface{})
	DEBUG(msg string, kv ...interface{})

	// Returns a function which will log at the given level if called and a boolean
	// indicating whether the log level is enabled.
	// To be used like:
	// if f,ok := log.ERRORok(); ok {f("message")}
	ALERTok() (LogFunc, bool)
	CRITok() (LogFunc, bool)
	ERRORok() (LogFunc, bool)
	WARNok() (LogFunc, bool)
	NOTICEok() (LogFunc, bool)
	INFOok() (LogFunc, bool)
	DEBUGok() (LogFunc, bool)
}

// LevelLoggerFull includes more methods than those needed for actual logging
type LevelLoggerFull interface {
	LevelLogger

	Level() syslog.Priority

	SetLevel(level syslog.Priority) bool
	IncLevel() bool
	DecLevel() bool

	With(kv ...interface{}) *Logger
}

//---

// StdFormatter allows quering a formatting handler for flags and prefix compatible with the
// stdlib log library
type StdFormatter interface {
	Flags() int
	Prefix() string
}

// StdMutableFormatter is the interface for a Logger which directly
// can change the stdlib flags, prefix and output io.Writer attributes in a synchronized manner.
// Since gonelog Handlers are immutable, it's not used for Formatters.
type StdMutableFormatter interface {
	StdFormatter
	SetFlags(flag int)
	SetPrefix(prefix string)
	SetOutput(w io.Writer)
}

// StdLoggerFull is mostly for documentation purposes. This is the full set of methods
// supported by the standard logger. You would only use the extra methods when you know
// exactly which kind of logger you are dealing with anyway.
type StdLoggerFull interface {
	StdLogger
	StdMutableFormatter
	Output(calldepth int, s string) error
}
