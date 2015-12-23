package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
	"github.com/One-com/gonelog/syslog"
	"github.com/One-com/gonelog/term"
)


// NewStdFormatter creates a Handler which formats log lines compatible with the std library.
// Note that the std library assumes is has to add a mutex to the io.Writer for sync - unlike
// when using a pre-synced io.Writer like SyncWriter
func NewStdFormatter(w io.Writer, prefix string, flag int) *stdformatter {
	f := &stdformatter{
		out:    w,
		prefix: prefix,
		flag:   flag,
	}
	return f
}

// stdformatter is right out of the standard library
// Plus ability to log KV attributes and a level prefix
// No ordering possible
type stdformatter struct {
	flag   int    // controlling the format
	prefix string // prefix to write at beginning of each line, after any level/timestamp
	out    io.Writer

	// A buffer for when we only have one global lock around
	// both formatting and writing
	mu  sync.Mutex
	buf []byte
}

func (f *stdformatter) Flags() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.flag
}
func (f *stdformatter) Prefix() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.prefix
}

// SetFlags sets the output flags for the logger.
func (f *stdformatter) SetFlags(flag int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flag = flag
}

// SetPrefix sets the output prefix for the logger.
func (f *stdformatter) SetPrefix(prefix string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.prefix = prefix
}
func (f *stdformatter) SetOutput(w io.Writer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.out = w
}

func (l *stdformatter) formatHeader(buf *[]byte, level syslog.Priority, t time.Time, file string, line int) {

	// Add level to std log features.
	if l.flag&(Llevel) != 0 {
		// color if asked to.
		if l.flag&(Lcolor) != 0 {
			*buf = append(*buf,
				fmt.Sprintf("\x1b[%sm",
					level_colors[level])...)
		}

		*buf = append(*buf, term_lvlpfx[level]...) // level prefix

		// end any coloring
		if l.flag&(Lcolor) != 0 {
			*buf = append(*buf, "\x1b[0m"...)
		}
	}

	*buf = append(*buf, l.prefix...)

	if l.flag&LUTC != 0 {
		t = t.UTC()
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	// Allow for PID as extra feature
	if l.flag&(Lpid) != 0 {
		*buf = append(*buf, " ["...)
		itoa(buf, pid, -1)
		*buf = append(*buf, "] "...)
	}

	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *stdformatter) Log(e Event) error {
	now := e.Time() // get this early.

	var file string
	var line int
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if e.fok {
			file, line = e.FileInfo()
		} else {
			file = "???"
			line = 0
		}
	}

	s := e.Msg

	// By Locking here we can share the same buffer between all log events going
	// through this formatter.
	l.mu.Lock()
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, e.Lvl, now, file, line)
	l.buf = append(l.buf, s...)

	// Arhh.. we need to not throw info away
	if len(e.Data) > 0 {
		var xbuf bytes.Buffer
		xbuf.WriteString(" ")
		marshalKeyvals(&xbuf, e.Data...)
		l.buf = append(l.buf, xbuf.Bytes()...)
	}
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	// Now write the message to the tree of chained writers.
	// If the tree root is a EventWriter, provide the original event too.
	var err error
	if ev, ok := l.out.(EvWriter); ok {
		_, err = ev.EvWrite(e, l.buf)
	} else {
		_, err = l.out.Write(l.buf)
	}

	l.mu.Unlock()
	return err
}

// Just set the color flag if it seems like TTY
func (f *stdformatter) AutoColoring() {
	var istty bool

	f.mu.Lock()
	defer f.mu.Unlock()

	w := f.out
	if tw, ok := w.(MaybeTtyWriter); ok {
		istty = tw.IsTty()
	} else {
		istty = term.IsTty(w)
	}

	if istty {
		f.flag = f.flag | Lcolor
	} else {
		f.flag = f.flag & ^Lcolor
	}
}
