package metric_test

import (
	"github.com/One-com/gonelog/metric"
	"log"
	"time"
	"strconv"
)

func ExampleNewClient() {

	var _flushPeriod = 4*time.Second
	
	sink, err := metric.NewTestSink("prefix", 1432)
	if err != nil {
		log.Fatal(err)
	}
	flushPeriod := metric.FlushInterval(_flushPeriod)

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

