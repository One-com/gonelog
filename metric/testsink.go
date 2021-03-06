package metric

import (
	"strconv"
	"fmt"
)


type testSink struct {
	max    int
	prefix string
}

// a non-thread safe surrogate of the sink
type testFlusherSink struct {
	prefix string // "prefix."
	max int
	buf []byte
}

func NewTestSink(prefix string, maxDatagramsize int) (sink *testSink, err error) {
	sink = &testSink{prefix: prefix, max : maxDatagramsize}
	return
}

func (s *testSink) FlusherSink() FlusherSink {
	fs := &testFlusherSink{max: s.max, prefix: s.prefix + "."}
	fs.buf = make([]byte,0,1500)
	return fs
}

func (s *testFlusherSink) Emit(name string, mtype int, value interface{}) {
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

func (s *testFlusherSink) EmitNumeric64(name string, mtype int, value Numeric64) {
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

func (s *testFlusherSink) flushIfBufferFull(lastSafeLen int) {
	if len(s.buf) > s.max {
		s.flush(lastSafeLen)
	}
}

func (s *testFlusherSink) Flush() {
	s.flush(0)
}

func (s *testFlusherSink) flush(n int) {
	if len(s.buf) == 0 {
		return
	}
	if n == 0 {
		n = len(s.buf)
	}

	fmt.Println(string(s.buf[:n-1]))

	if n < len(s.buf) {
		copy(s.buf, s.buf[n:])
	}
	s.buf = s.buf[:len(s.buf)-n]
}

func (s *testFlusherSink) appendType(t int) {
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

func (s *testFlusherSink) appendNumeric64(v Numeric64) {
	switch v.Type {
	case Uint64:
		s.buf = strconv.AppendUint(s.buf, v.Uint64(), 10)		
	case Int64:
		s.buf = strconv.AppendInt(s.buf, v.Int64(), 10)
	case Float64:
		s.buf = strconv.AppendFloat(s.buf, v.Float64(), 'f', -1, 64)
	}
}

func (s *testFlusherSink) appendNumber(v interface{}) {
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
	
