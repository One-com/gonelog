package metric

import (
	"time"
)

// A histogram is a series of int64 events all sent to the server
type Histogram eventStream

func NewHistogram(name string, opts ...MOption) *Histogram {
	return default_client.NewHistogram(name, opts...)
}

func (c *Client) NewHistogram(name string, opts ...MOption) *Histogram {
	dequeuef := func(f FlusherSink, val uint64) {
		n := Numeric64{Type: Int64, value: val}
		f.EmitNumeric64(name, MeterHistogram, n)
	}
	t := c.newEventStream(name, MeterHistogram, dequeuef)
	c.register(t, opts...)
	return (*Histogram)(t)
}

// Update the Histogram by registering a new event
func (h *Histogram) Update(d int64) {
	e := (*eventStream)(h)
	e.Enqueue(uint64(d))
}

// Timer is like Histogram, but the event is a time.Duration.
// values are remembered as milliseconds
type Timer eventStream

func NewTimer(name string, opts ...MOption) *Timer {
	return default_client.NewTimer(name, opts...)
}

func (c *Client) NewTimer(name string, opts ...MOption) *Timer {
 	dequeuef := func(f FlusherSink, val uint64) {
		n := Numeric64{Type: Uint64, value: val}
		f.EmitNumeric64(name, MeterTimer, n)
	}
	t := c.newEventStream(name, MeterTimer, dequeuef)
	c.register(t, opts...)
	return (*Timer)(t)
}

// Register a new duration event.
func (t *Timer) Update(d time.Duration) {
	e := (*eventStream)(t)
	e.Enqueue(uint64(d.Nanoseconds()/int64(1000000)))
}
