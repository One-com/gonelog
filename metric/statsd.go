package metric

import (
	"net"
	"strconv"
	"time"
	"fmt"
)

type config struct {
	Addr          string
	MaxPacketSize int
}


type Option func(*config)


type statsdSink struct {
	conn   net.Conn
	max    int
	prefix string
}

// a non-thread safe surrogate of the sink
type statsdFlusherSink struct {
	prefix string // "prefix."
	conn net.Conn
	max int
	buf []byte
}

// A Sink sending data to a UDP statsd server
func NewStatsdSink(addr string, prefix string, maxDatagramsize int) (sink *statsdSink, err error) {

	conn, err := net.DialTimeout("udp", addr, time.Second)
	if err != nil {
		return
	}
	sink = &statsdSink{conn:conn, prefix: prefix, max : maxDatagramsize}
	return
}

func (s *statsdSink) FlusherSink() FlusherSink {
	fs := &statsdFlusherSink{conn:s.conn, max: s.max, prefix: s.prefix + "."}
	fs.buf = make([]byte,0,1500)
	return fs
}

func (s *statsdFlusherSink) Emit(name string, mtype int, value interface{}) {
	curbuflen := len(s.buf)
	s.buf = append(s.buf,s.prefix...)
	s.buf = append(s.buf,name...)
	s.buf = append(s.buf,':')
	switch v := value.(type) {
	case string:
		s.buf = append(s.buf,v...)
	case fmt.Stringer:
		s.buf = append(s.buf,v.String()...)
	default:
		panic("Not stringable")
	}
	s.buf = append(s.buf,'|')
	s.appendType(mtype)
	// sampe rate not supported
	s.buf = append(s.buf,'\n')
	s.flushIfBufferFull(curbuflen)
}

func (s *statsdFlusherSink) EmitNumeric64(name string, mtype int, value Numeric64) {
	curbuflen := len(s.buf)
	s.buf = append(s.buf,s.prefix...)
	s.buf = append(s.buf,name...)
	s.buf = append(s.buf,':')
	s.appendNumeric64(value)
	s.buf = append(s.buf,'|')
	s.appendType(mtype)
	// sampe rate not supported
	s.buf = append(s.buf,'\n')
	s.flushIfBufferFull(curbuflen)
}

func (s *statsdFlusherSink) flushIfBufferFull(lastSafeLen int) {
	if len(s.buf) > s.max {
		s.flush(lastSafeLen)
	}
}

func (s *statsdFlusherSink) Flush() {
	s.flush(0)
}

func (s *statsdFlusherSink) flush(n int) {
	if len(s.buf) == 0 {
		return
	}
	if n == 0 {
		n = len(s.buf)
	}

	// Trim the last \n, StatsD does not like it.
//	fmt.Println(string(s.buf[:n-1]))
	s.conn.Write(s.buf[:n-1])
//	c.handleError(err)
	if n < len(s.buf) {
		copy(s.buf, s.buf[n:])
	}
	s.buf = s.buf[:len(s.buf)-n]
}

func (s *statsdFlusherSink) appendType(t int) {
	switch t {
	case MeterGauge:
		s.buf = append(s.buf,'g')
	case MeterCounter:
		s.buf = append(s.buf,'c')
	case MeterTimer, MeterHistogram:  // until we are sure the statsd server supports otherwise
		s.buf = append(s.buf,"ms"...)
	case MeterSet:
		s.buf = append(s.buf,'s')

	}
}

func (s *statsdFlusherSink) appendNumeric64(v Numeric64) {
	switch v.Type {
	case Uint64:
		s.buf = strconv.AppendUint(s.buf, v.Uint64(), 10)
	case Int64:
		s.buf = strconv.AppendInt(s.buf, v.Int64(), 10)
	case Float64:
		s.buf = strconv.AppendFloat(s.buf, v.Float64(), 'f', -1, 64)
	}
}

func (s *statsdFlusherSink) appendNumber(v interface{}) {
	switch n := v.(type) {
	case int:
		s.buf = strconv.AppendInt(s.buf, int64(n), 10)
	case uint:
		s.buf = strconv.AppendUint(s.buf, uint64(n), 10)
	case int64:
		s.buf = strconv.AppendInt(s.buf, n, 10)
	case uint64:
		s.buf = strconv.AppendUint(s.buf, n, 10)
	case int32:
		s.buf = strconv.AppendInt(s.buf, int64(n), 10)
	case uint32:
		s.buf = strconv.AppendUint(s.buf, uint64(n), 10)
	case int16:
		s.buf = strconv.AppendInt(s.buf, int64(n), 10)
	case uint16:
		s.buf = strconv.AppendUint(s.buf, uint64(n), 10)
	case int8:
		s.buf = strconv.AppendInt(s.buf, int64(n), 10)
	case uint8:
		s.buf = strconv.AppendUint(s.buf, uint64(n), 10)
	case float64:
		s.buf = strconv.AppendFloat(s.buf, n, 'f', -1, 64)
	case float32:
		s.buf = strconv.AppendFloat(s.buf, float64(n), 'f', -1, 32)
	}
}

/* Some of the above code has been borrowed from github.com/alexcesaro/statsd

... which carries the license:

The MIT License (MIT)

Copyright (c) 2015 Alexandre Cesaro

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
