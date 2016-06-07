package metric

import (
	"sync"
	"time"
)

type Client struct {

	// Wait for all flushers to empty before exiting Stop()
	done *sync.WaitGroup

	fmu sync.Mutex
	flushers map[time.Duration]*Flusher

	default_flusher *Flusher

	// A factory for flusher local sinks with a Emit() function
	sink Sink

	running bool
}

// A single global client. Remember to set a default sink
var default_client *Client

func init() {
	default_client = NewClient(nil)
}

// Create a new metric client with a factory object for the sink.
// If sink == nil, the client will not emit metrics until a Sink is set.
func NewClient(sink Sink, opts ...MOption) (client *Client) {

	client = &Client{sink: sink}
	client.done =  new(sync.WaitGroup)
	client.flushers = make(map[time.Duration]*Flusher)

	client.default_flusher = newFlusher(0)

	go client.default_flusher.rundyn()

	client.Start()
	client.SetOptions(opts...)

	return
}

func SetDefaultOptions(opts ...MOption) {
	c := default_client
	c.SetOptions(opts...)
}

// The the Sink factory of the client
func (c *Client) SetOptions(opts ...MOption) {
	conf := make(map[string]interface{})
	for _, o := range opts {
		o(conf)
	}

	c.fmu.Lock()
	if f,ok := conf["flushInterval"]; ok {
		if d, ok := f.(time.Duration) ; ok {
			c.default_flusher.set_interval(d)
		}
	}
	c.fmu.Unlock()
}


// Set the Sink factory of the default client
func SetDefaultSink(sink Sink) {
	c := default_client
	c.SetSink(sink)
}

// The the Sink factory of the client
func (c *Client) SetSink(sink Sink) {
	c.fmu.Lock()

	c.sink = sink
	c.default_flusher.setsink(sink.FlusherSink())
	for _,f := range c.flushers {
		if sink == nil {
			fsink := &nilFlusherSink{}
			f.setsink(fsink)
		} else {
			fsink := c.sink.FlusherSink()
			f.setsink(fsink)
		}
	}
	c.fmu.Unlock()
}

// Start the default client if stopped.
func Start() {
	default_client.Start()
}

// Start a stopped client
func (c *Client) Start() {
	c.fmu.Lock()
	defer c.fmu.Unlock()

	if c.running {
		return
	}

	c.default_flusher.stop() // the default flusher flip/flops

	for _,f := range c.flushers {
		c.done.Add(1)
		f.run(c.done)
	}

	c.running = true
}

// Stops the global default metrics client
func Stop() {
	default_client.Stop()
}

// Stop a Client from flushing data.
// If any AutoFlusher meters are still in use they will still flush when overflown.
func (c *Client) Stop() {
	c.fmu.Lock()
	defer c.fmu.Unlock()

	if !c.running {
		return
	}

	c.default_flusher.stop()

	for _,f := range c.flushers {
		f.stop()
	}
	c.done.Wait()
	c.running = false
}

func (c *Client) register(m Meter, opts ...MOption) {
	c.fmu.Lock()

	var f *Flusher
	var flush time.Duration

	conf := make(map[string]interface{})
	for _, o := range opts {
		o(conf)
	}

	if fi,ok := conf["flushInterval"]; ok {
		flush = fi.(time.Duration)
		if f, ok = c.flushers[flush]; !ok {
			f = newFlusher(flush)
			if c.sink != nil {
				fsink := c.sink.FlusherSink()
				f.setsink(fsink)
			}
			c.flushers[flush] = f
			if c.running {
				c.done.Add(1)
				go f.run(c.done)
			}
		}
		f.register(m)
	} else {
		c.default_flusher.register(m)
	}
	c.fmu.Unlock()
}
