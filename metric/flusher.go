package metric

import (
	"time"
	"sync"
//	"log"
)

type Flusher struct {
	stopChan     chan struct{}
	interval time.Duration

	mu sync.Mutex
	meters []Meter

	sink FlusherSink
}

func newFlusher(interval time.Duration) *Flusher {
	f := &Flusher{interval: interval, sink: &nilFlusherSink{}}
	f.stopChan = make(chan struct{})
	return f
}

func (f *Flusher) setsink(sink FlusherSink) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sink = sink
}

func (f *Flusher) stop() {
	close(f.stopChan)
}

func (f *Flusher) run(done *sync.WaitGroup) {
	defer done.Done()

	if f.interval == 0 {
		return
	}
	ticker := time.NewTicker(f.interval)
LOOP:
	for {
		select {
		case <- f.stopChan:
			ticker.Stop()
			break LOOP
		case <- ticker.C:
			f.Flush()
		}
	}
	f.Flush()
}

// flush a single Meter
func (f *Flusher) FlushMeter(m Meter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	m.Flush(f.sink)
}

// flush all meters
func (f *Flusher) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, m := range f.meters {
		m.Flush(f.sink)
	}
	f.sink.Flush()
}

func (f *Flusher) register(m Meter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.meters = append(f.meters, m)
	if a,ok := m.(AutoFlusher); ok {
		a.SetFlusher(f)
	}
}
