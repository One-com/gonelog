package metric

import (
	"time"
)

type Histogram eventStream

func NewHistogram(name string, opts ...MOption) *Histogram {
	return default_client.NewHistogram(name, opts...)
}

func (c *Client) NewHistogram(name string, opts ...MOption) *Histogram {
	t := c.newEventStream(MeterHistogram, name)
	c.register(t, opts...)
	return (*Histogram)(t)
}

func (t *Histogram) Update(d int64) {
	e := (*eventStream)(t)
	n := Numeric64{Type: Int64, value: uint64(d)}
	e.Record(n)
}


type Timer eventStream

func NewTimer(name string, opts ...MOption) *Timer {
	return default_client.NewTimer(name, opts...)
}

func (c *Client) NewTimer(name string, opts ...MOption) *Timer {
	t := c.newEventStream(MeterTimer, name)
	c.register(t, opts...)
	return (*Timer)(t)
}

func (t *Timer) Update(d time.Duration) {
	e := (*eventStream)(t)
	n := Numeric64{Type:Uint64, value: uint64(d.Nanoseconds()/int64(1000000))}
	e.Record(n)
}
