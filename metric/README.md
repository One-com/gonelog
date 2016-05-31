# gonelog metrics

Fast Golang metrics library

Package gonelog/metric is an expandable library for metrics. Initally only sending data to statsd

The design goals:

* Basic statsd metric types (gauge, counter, timer, set)
* Client side buffered
* Fast

Timer and Histogram is basically the same except for the argument type.

Counter is reset to zero on each flush. Gauges are not.

Still somewhat experimental.

## Example

```go
import "github.com/One-com/gonelog/metric"

var flushPeriod = 2*time.Second

func main() {

	sink, err := metric.NewStatsdSink("1.2.3.4:8125", "prefix", 1432)
	if err != nil {
		log.Fatal(err)
	}
	flushPeriod := metric.FlushInterval(flushPeriod)
	c := metric.NewClient(sink,flushPeriod)
	
	gauge   := c.NewGauge("gauge",flushPeriod)
	timer   := c.NewTimer("timer")
	histo   := c.NewHistogram("histo",flushPeriod)
	counter := c.NewCounter("counter",flushPeriod)
	set     := c.NewSet("set", flushPeriod)

	g := 100
	for g != 0 {
		counter.Inc(1)
		gauge.Update(uint64(g))
		timer.Update(time.Duration(g)*time.Millisecond)
		histo.Update(int64(g))
		set.Update(strconv.FormatInt(int64(g), 10))
		
		time.Sleep(time.Second)
		g--
	}
	c.Stop()
	
}
```
