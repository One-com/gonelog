package metric

import (
	"sync/atomic"
	"math"
)

// A client maintained gauge which is only sampled regulary without information loss
// wrt. the absolute value.
// Can be used as a client side maintained counter too.

type GaugeUint64 struct {
	name string
	val  uint64
}

type GaugeInt64 struct {
	name string
	val  int64
}

// A go-metric compatible float64 gauge which stores its value as a uint64
// to implement Flush() fast
type GaugeFloat64 struct {
	name string
	val  uint64
}

/*==================================================================*/
// Gauge is alias for GaugeUint64
func NewGauge(name string, opts ...MOption) *GaugeUint64 {
	return default_client.NewGauge(name, opts...)
}

func (c *Client) NewGauge(name string, opts ...MOption) *GaugeUint64 {
	return c.NewGaugeUint64(name, opts...)
}
/*==================================================================*/

func (c *Client) NewGaugeUint64(name string, opts ...MOption) *GaugeUint64 {
	g := &GaugeUint64{name: name}
	c.register(g, opts...)
	return g
}

// Flush to implement Meter interface
func (g *GaugeUint64) Flush(f FlusherSink) {
	val := atomic.LoadUint64(&g.val)
	n := Numeric64{Type:Uint64, value: val}
	f.EmitNumeric64(g.name, MeterGauge, n)
}

// Name to implement Meter interface
func (g *GaugeUint64) Name() string {
	return g.name
}

// Mtype to implement Meter interface
func (g *GaugeUint64) Mtype() int {
	return MeterGauge
}

func (g *GaugeUint64) Update(val uint64) {
	atomic.StoreUint64(&g.val, val)
}

func (g *GaugeUint64) Value() uint64 {
	return atomic.LoadUint64(&g.val)
}

func (g *GaugeUint64) Inc(i uint64) {
	atomic.AddUint64(&g.val, i)
}

func (g *GaugeUint64) Dec(i int64) {
	atomic.AddUint64(&g.val, ^uint64(i-1))
}

/*==================================================================*/

// An Int64 Gauge. Can be used as a go-metric client side gauge or counter
func (c *Client) NewGaugeInt64(name string, opts ...MOption) *GaugeInt64 {
	g := &GaugeInt64{name: name}
	c.register(g, opts...)
	return g
}

func (g *GaugeInt64) Flush(f FlusherSink) {
	val := atomic.LoadInt64(&g.val)
	n := Numeric64{Type:Int64, value: uint64(val)}
	f.EmitNumeric64(g.name, MeterGauge, n)
}

func (g *GaugeInt64) Name() string {
	return g.name
}

func (g *GaugeInt64) Mtype() int {
	return MeterGauge
}

func (g *GaugeInt64) Update(val int64) {
	atomic.StoreInt64(&g.val, val)
}

func (g *GaugeInt64) Value() int64 {
	return atomic.LoadInt64(&g.val)
}

// Clear sets the counter to zero.
func (g *GaugeInt64) Clear() {
	atomic.StoreInt64(&g.val, 0)
}

// Count returns the current count.
func (g *GaugeInt64) Count() int64 {
	return g.Value()
}

// Dec decrements the counter by the given amount.
func (g *GaugeInt64) Dec(i int64) {
	atomic.AddInt64(&g.val, -i)
}

// Inc increments the counter by the given amount.
func (g *GaugeInt64) Inc(i int64) {
	atomic.AddInt64(&g.val, i)
}

/*==================================================================*/
// An Float64 Gauge.
func (c *Client) NewGaugeFloat64(name string, opts ...MOption) *GaugeFloat64 {
	g := &GaugeFloat64{name: name}
	c.register(g, opts...)
	return g
}

func (g *GaugeFloat64) Name() string {
	return g.name
}

func (g *GaugeFloat64) Mtype() int {
	return MeterGauge
}

// Update updates the gauge's value.
func (g *GaugeFloat64) Update(v float64) {
	atomic.StoreUint64(&g.val, math.Float64bits(v))

}

// Value returns the gauge's current value.
func (g *GaugeFloat64) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.val))
}

func (g *GaugeFloat64) Flush(f FlusherSink) {
	val := atomic.LoadUint64(&g.val)
	n := Numeric64{Type:Float64, value: val}
	f.EmitNumeric64(g.name, MeterGauge, n)
}
