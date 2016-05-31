package metric

import (
	"sync"
	"sync/atomic"
)

// An almost lock-free FIFO buffer
// Locks are only used when flushing

const bufferMaskBits = 10

const bufferSize     = uint64(1) << bufferMaskBits
const bufferMask     = bufferSize - 1
const bufferMaxCycle = uint64(1) << (64-bufferMaskBits)

const maxCycle    = 4096

const _ = uint(bufferMaxCycle-maxCycle) // assert at compile time

type event struct {
	seq uint64
	val uint64 // union of 64-bit numeric values including float64
}

// A generic stream of values which all have to be propagated to the sink.
type eventStream struct {
	name string
	mtype int // The conceptual type of the meter
	etype int // the data type of the stored event value

	widx uint64 // index of next free slot
	ridx uint64 // index of next unread slot
	slots [bufferSize]event

	flusher *Flusher

	mu   sync.Mutex
	cond *sync.Cond
}


func (c *Client) newEventStream(mtype int, name string) (e *eventStream) {
	e = &eventStream{name: name, mtype: mtype}
	switch mtype {
	case MeterGauge, MeterTimer:
		e.etype = Uint64
	case MeterCounter, MeterHistogram:
		e.etype = Int64
	default:
		e = nil
		return
	}

	// make sure first slot is not valid from the start due to zero-value
	// by invalidating it
	e.slots[0].seq = 1

	e.cond = sync.NewCond(&e.mu)
	return e
}

func (e *eventStream) SetFlusher(f *Flusher) {
	e.flusher = f
}

// Empty the buffer
func (e* eventStream) Flush(f FlusherSink) {

	var idx uint64

	// Precondition: e.ridx points to next un-eaten slot
	ridx := atomic.LoadUint64(&e.ridx)
	for {
		idx = ridx & bufferMask

		mark := atomic.LoadUint64(&(e.slots[idx].seq))
		if mark == ridx {
			// This is a valid slot
			val := e.slots[idx].val
			n := Numeric64{Type: e.etype, value: val}
			f.EmitNumeric64(e.name, e.mtype, n)
			ridx++
		} else {
			// we've reached a not yet written slot
			break
		}
	}
	// how far did we get?
	atomic.StoreUint64(&(e.ridx), ridx)
}


func (e *eventStream) Name() string {
	return e.name
}

func (e *eventStream) Mtype() int {
	return e.mtype
}

func (e *eventStream) reset() {
	// While we are here, we have no writers
	// They'll all run into maxCycle bein reached
	// - until we reset widx, then they are unleashed. So do that last

	e.flusher.FlushMeter(e)
	// buffer should now be empty

	atomic.StoreUint64(&e.ridx,0)

	// atomically let the whole thing run again
	atomic.StoreUint64(&e.widx,0)
}

func (e *eventStream) Record(val Numeric64) {

	var ridx    uint64
	var widx    uint64
	var idx     uint64
	var cycle   uint64

	// First get a slot
	for {
		widx = atomic.AddUint64(&e.widx, 1)
		widx-- // back up to get our reserved slot
		idx = widx & bufferMask
		cycle = widx >> bufferMaskBits

		if cycle >= maxCycle {
			e.mu.Lock() // stall all other appenders, but not flusher
			// let the one putting us over the top reset
			if idx == 0 && cycle == maxCycle {
				// flush and reset the whole buffer
				e.reset()
				e.cond.Broadcast()
			} else {
				e.cond.Wait()
			}
			e.mu.Unlock() // and redo slot reservation above
		} else {
			break
		}
	}

	// Then write the data
	for {
		// Where's the reader? Don't overtake it.
		ridx = atomic.LoadUint64(&e.ridx)

		if widx - ridx < bufferSize {
			// We have not catched up
			e.slots[idx].val = val.value
			// mark the slot written
			atomic.StoreUint64(&(e.slots[idx].seq),widx)
			break
		} else {
			// do some flushing
			e.flusher.FlushMeter(e)
		}
	}

}
