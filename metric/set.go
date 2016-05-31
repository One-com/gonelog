package metric

import (
	"sync"
)


type Set struct {
	name string
	flusher *Flusher

	mu sync.Mutex
	set map[string]struct{}
}


func NewSet(name string, opts ...MOption) *Set {
	return default_client.NewSet(name, opts...)
}

func (c *Client) NewSet(name string, opts ...MOption) *Set {
	s := &Set{name: name}
	s.set = make(map[string]struct{})
	c.register(s, opts...)
	return s
}


func (s *Set) Flush(f FlusherSink) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k,_ := range s.set {
		f.Emit(s.name, MeterSet, k)
	}
	s.set = make(map[string]struct{})
}

func (s *Set) Name() string {
	return s.name
}

func (s *Set) Mtype() int {
	return MeterSet
}

func (s *Set) Update(val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[val] = struct{}{}
}

func (s *Set) SetFlusher(f *Flusher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flusher = f
}
