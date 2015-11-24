// +build go1.5

package log

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Indirection of Log() calls through a Handler which can be atomically swapped

var ErrNotLogged = errors.New("No handler found to log event")

// using atomic and mutex to support atomic reads, but also read-modify-write ops.
type swapper struct {
	mu  sync.Mutex // Locked by any who want to modify the valueStruct
	val atomic.Value
}

type valueStruct struct {
	Handler
	parent *Logger
}

// makes sure to initialize a swapper with a value
func new_swapper() (s *swapper) {
	s = new(swapper)
	s.val.Store(valueStruct{})
	return
}

// Log sends the event down the first Handler chain, it finds in the Logger tree.
// NB: This is different from pythong "logging" in that only one Handler is activated
func (h *swapper) Log(e *event) (err error) {

	// try the local handler
	v, _ := h.val.Load().(valueStruct)
	// Logger swappers *has* to have a valid valueStruct

	if v.Handler != nil {
		err = v.Handler.Log(Event{e})
	} else {
		// Have to try parents. Walk the name-tree to find the first handler
		cur := v.parent.h
		for cur != nil {
			v, _ := cur.val.Load().(valueStruct) // must be valid
			if v.Handler != nil {
				err = v.Handler.Log(Event{e})
				break
			} else {
				cur = v.parent.h
			}
		}
	}

	freePoolEvent(e)
	return err
}

func (s *swapper) SwapParent(new *Logger) (old *Logger) {
	s.mu.Lock()
	v := s.val.Load().(valueStruct)
	old = v.parent
	h := v.Handler
	s.val.Store(valueStruct{Handler: h, parent: new})
	s.mu.Unlock()
	return
}

func (s *swapper) SwapHandler(new Handler) {
	s.mu.Lock()
	p := (s.val.Load().(valueStruct)).parent
	s.val.Store(valueStruct{Handler: new, parent: p})
	s.mu.Unlock()
}

// SwapClone swaps in a handler from another swapper
func (s *swapper) swapClone(source *swapper) {
	s.mu.Lock()
	sh := source.handler()
	p := (s.val.Load().(valueStruct)).parent
	s.val.Store(valueStruct{Handler: sh, parent: p})
	s.mu.Unlock()
}

func (s *swapper) handler() Handler {
	return (s.val.Load().(valueStruct)).Handler
}

func (s *swapper) parent() *Logger {
	return (s.val.Load().(valueStruct)).parent
}
