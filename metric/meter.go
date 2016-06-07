package metric

// Conceptual meter types
// Gauge: A client side maintained counter
// Counter: A server side maintained counter
// Historam: A series of events analyzed on the server
// Timer: A series of time.Duration events analyzed on the server
// Set: Discrete strings added to a set maintained on the server
const (
	MeterGauge = iota  // A client side maintained value
	MeterCounter       // A server side maintained value
	MeterHistogram     // A general distribution of measurements events.
	MeterTimer         // ... when those measurements are milliseconds
	MeterSet           // free form string events
)

// Meters should Emit*() to a FlusherSink When asked to flush.
// The idea is that the Sink in principle needs to be go-routine safe,
// but FlusherSink.Emit*() will be called synchronized, so you can spare
// the synchronization in the Emit*() implementation it self if you provide a
// Flusher specific version of the Sink object.

// FlusherSink is a Sink for metrics data which methods are guaranteed to be called
// synchronized. This it can keep state like buffers without locking
type FlusherSink interface {
	Emit(name string, mtype int, value interface{})
	EmitNumeric64(name string, mtype int, value Numeric64)
	Flush()
}

// A factory for FlusherSinks. Metric clients are given a Sink to output to.
type Sink interface {
	FlusherSink() FlusherSink
}

// A Meter measures stuff and can be registered with a client to
// be periodically reported to the Sink
type Meter interface {
	Mtype() int
	Name() string
	Flush(FlusherSink) // Read the meter, by flushing all non-read values.
}

// An AutoFlusher can initiate a Flush throught the flusher at any time and needs
// to know the Flusher to call FlushMeter() on it
type AutoFlusher interface {
	SetFlusher(*Flusher)
}

// Flushers are created with this sink which just throws away data
// until a real sink is set.
// It's the user responsibility to not generate metrics before setting a sink if this
// is not wanted.
type nilFlusherSink struct{}

//func (n *nilFlusherSink) Emit(name string, mtype int, value int64) {
func (n *nilFlusherSink) Emit(name string, mtype int, value interface{}) {
}

func (n *nilFlusherSink) EmitNumeric64(name string, mtype int, value Numeric64) {
}

func (n *nilFlusherSink) Flush() {
}
