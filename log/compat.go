package log

import (
	"fmt"
	"github.com/One-com/gonelog/syslog"
	"io"
	"os"
)

// This will instantiate a logger with the same functionality (and limitations) as the std lib logger.
func New(out io.Writer, prefix string, flags int) *Logger {
	h := NewStdFormatter(out, prefix, flags)
	l := NewLogger(LvlDEFAULT, h) // not a part of the hierachy
	l.DoTime(true)
	return l
}

func (c *Logger) Output(calldepth int, s string) error {
	return c.output(calldepth, s)
}

func (c *Logger) Fatal(v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprint(v...)
	c.log(l, s)
	os.Exit(1)
}
func (c *Logger) Fatalf(format string, v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprintf(format, v...)
	c.log(l, s)
	os.Exit(1)
}
func (c *Logger) Fatalln(v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprintln(v...)
	c.log(l, s)
	os.Exit(1)
}

func (c *Logger) Panic(v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprint(v...)
	c.log(l, s)
	panic(s)
}
func (c *Logger) Panicf(format string, v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprintf(format, v...)
	c.log(l, s)
	panic(s)
}
func (c *Logger) Panicln(v ...interface{}) {
	l := syslog.LOG_ALERT
	s := fmt.Sprintln(v...)
	c.log(l, s)
	panic(s)
}

func (c *Logger) Print(v ...interface{}) {
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprint(v...)
		c.log(l, s)
	}
}
func (c *Logger) Printf(format string, v ...interface{}) {
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprintf(format, v...)
		c.log(l, s)
	}
}
func (c *Logger) Println(v ...interface{}) {
	if l, ok := c.DoingDefaultLevel(); ok {
		s := fmt.Sprintln(v...)
		c.log(l, s)
	}
}

//---

// These functions have been delegated to the swapper, since some of them might
// need to replace the handler.
// If these functions have no meaning for the actual Handler attached, then they
// result in a NOOP.

func (c *Logger) Flags() int {
	return c.h.Flags()
}

func (c *Logger) Prefix() string {
	return c.h.Prefix()
}

func (c *Logger) SetFlags(flag int) {
	// First activate needed book keeping
	if flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		c.DoTime(true)
	}
	if flag&(Llongfile|Lshortfile) != 0 {
		c.DoCodeInfo(true)
	}

	c.h.SetFlags(flag)

	// De-activate unneeded book keeping
	if flag&(Ldate|Ltime|Lmicroseconds) == 0 {
		c.DoTime(false)
	}
	if flag&(Llongfile|Lshortfile) == 0 {
		c.DoCodeInfo(false)
	}
}

func (c *Logger) SetPrefix(prefix string) {
	c.h.SetPrefix(prefix)
}

func (c *Logger) SetOutput(w io.Writer) {
	c.h.SetOutput(w)
}
