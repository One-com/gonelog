package metric

import (
	"time"
	"sync"
)

// A flusher is either run with a fixed flushinterval with a go-routine which
// exits on stop(), or with a dynamic changeable flushinterval in a permanent go-routine.
// This is chosen by either calling run() og rundyn()
const (
	flusherTypeUndef = iota
	flusherTypeFixed
	flusherTypeDynamic
)

type Flusher struct {
	stopChan     chan struct{}
	kickChan     chan struct{}

	interval time.Duration

	mu sync.Mutex
	meters []Meter
	ftype int // only set once by the run/rundyn method

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
	f.stopChan <- struct{}{}
}

func (f *Flusher) set_interval(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.ftype == flusherTypeDynamic {
		f.interval = d
		f.kickChan <- struct{}{}
	}
}

// A go-routine wich will flush at adjustable intervals and doesn't
// exit if interval is zero.
func (f *Flusher) rundyn() {
	var interval time.Duration

	f.mu.Lock()
	if f.ftype != flusherTypeFixed {
		f.ftype = flusherTypeDynamic
	} else {
		panic("Attempt to make fixed flusher dynamic")
	}

	f.kickChan = make(chan struct{})
	f.mu.Unlock()

	var ticker *time.Ticker

	for {
	STOPPED:
		for {
			select {
			case <- f.stopChan: // wait for start signal
				break STOPPED
			case <-f.kickChan: // just accept being kicked
			}
		}
	RUNNING: // two cases - either with a flush or not
		for {
			f.mu.Lock()
			interval = f.interval
			f.mu.Unlock()
			if interval == 0 {
				select {
				case <- f.stopChan:
					break RUNNING
				case <-f.kickChan:
				}
			} else {
				ticker = time.NewTicker(interval)
			LOOP:
				for {
					select {
					case <- f.stopChan:
						ticker.Stop()
						break RUNNING
					case <-f.kickChan:
						ticker.Stop()
						break LOOP // to test to make a new ticker
					case <- ticker.C:
						f.Flush()
					}
				}
			}
		}
		f.Flush()
	}
}

func (f *Flusher) run(done *sync.WaitGroup) {
	defer done.Done()

	f.mu.Lock()
	if f.ftype != flusherTypeDynamic {
		f.ftype = flusherTypeFixed
	} else {
		panic("Attempt to make default flusher fixed")
	}
	f.mu.Unlock()

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
	m.Flush(f.sink)
	f.mu.Unlock()
}

// flush all meters
func (f *Flusher) Flush() {
	f.mu.Lock()
	for _, m := range f.meters {
		m.Flush(f.sink)
	}
	f.sink.Flush()
	f.mu.Unlock()
}

func (f *Flusher) register(m Meter) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.meters = append(f.meters, m)
	if a,ok := m.(AutoFlusher); ok {
		a.SetFlusher(f)
	}
}
