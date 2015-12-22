package log

import (
	ilog "github.com/One-com/gonelog"
	"io"
)

// First some interfaces to test on what can be done with Handlers

// CloneableHandler is a Handler which can clone itself for modification with
// the purpose of being swapped in to replace the current handler - thus utilizing the
// atomiciy of the swapper to change Handler behaviour during flight.
// Handlers not being Cloneable must be manually manipulated by the application
// and replaced by Logger.SetHandler()
// Making a Handler Cloneable makes it possible for the framework to support the
// standard top-level operations on it like StdFormatter and AutoColorer
// The framework promises not to modify the Handler after it's in use.
// In other words, once the first Log() method is called the Handler is not modifed.
type CloneableHandler interface {
	Handler
	Clone() CloneableHandler
}

// USAutoColorer is the ability of a formatter to do AutoColoring detection
// before being swapped in to start Logging.
type USAutoColorer interface {
	UnsyncedAutoColoring()
}

// USFlagger is the ability of a formatter to set the formatter flags before
// being swapped in to start logging
type USFlagger interface {
	UnsyncedSetFlags(flag int)
}

// USPrefixer is the ability of a formatter to set the prefix before
// being swapped in to start logging
type USPrefixer interface {
	UnsyncedSetPrefix(prefix string)
}

// USOutputter is the ability of a formatter to set the output Writer before
// being swapped in to start logging
type USOutputter interface {
	UnsyncedSetOutput(w io.Writer)
}

type WriterGetter interface {
	// Not in the standard library, but needed here to swap formatting handlers
	GetWriter() io.Writer
}

// AutoColoring swaps in a equivalent handler doing AutoColoring if possible
func (h *swapper) AutoColoring() {
	old := h.handler()
	if clo, ok := old.(CloneableHandler); ok {
		if _, ok := old.(USAutoColorer); ok {
			new := clo.Clone()
			new.(USAutoColorer).UnsyncedAutoColoring()
			h.SwapHandler(new)
		}
	}
}

/*****************************************************************************/
// Functions for manipulating the stored handler in std lib compatible ways
// These functions are a no-op for handlers not supporting the concepts
// though the swapper goes out of its way to let as many handlers as possible support
// these operations by implementing the below interfaces

// Flags return the Handler flags. Since Handlers are not modfied after being swapped in
// (unless they are StdMutables) this is safe for all.
func (h *swapper) Flags() int {
	if handler, ok := h.handler().(ilog.StdFormatter); ok {
		return handler.Flags()
	}
	return 0
}

// Prefix - same as for flags
func (h *swapper) Prefix() string {
	if handler, ok := h.handler().(ilog.StdFormatter); ok {
		return handler.Prefix()
	}
	return ""
}

func (h *swapper) SetFlags(flag int) {
	old := h.handler()
	if handler, ok := old.(ilog.StdMutableFormatter); ok {
		handler.SetFlags(flag)
		return
	}
	// we have to atomically replace the handler with one with the new flag,
	// since locking can only be assumed to be done in 2 places:
	// swapper and stdformatter.out (when it's a syncwriter or equivalent),
	// nothing protects the flags of the formatter except replacing it entirely
	// Note that this is not a "compare-and-swap". A bad application might
	// end up swapping out another handler than the one it got the original
	// flags from. That's your own fault.
	// This operation only protects against outputting log-lines which
	// are not well defined for "some" handler.
	if clo, ok := old.(CloneableHandler); ok {
		if _, ok := old.(USFlagger); ok {
			new := clo.Clone()
			new.(USFlagger).UnsyncedSetFlags(flag)
			h.SwapHandler(new)
		}
	}
}

func (h *swapper) SetPrefix(prefix string) {
	old := h.handler()
	if handler, ok := old.(ilog.StdMutableFormatter); ok {
		handler.SetPrefix(prefix)
		return
	}
	if clo, ok := old.(CloneableHandler); ok {
		if _, ok := old.(USPrefixer); ok {
			new := clo.Clone()
			new.(USPrefixer).UnsyncedSetPrefix(prefix)
			h.SwapHandler(new)
		}
	}
}

func (h *swapper) SetOutput(w io.Writer) {
	old := h.handler()
	if handler, ok := old.(ilog.StdMutableFormatter); ok {
		handler.SetOutput(w)
		return
	}
	// changing output for some Handlers is actually possible without a swap,
	// courtesy of the syncwriter
	if f, ok := old.(WriterGetter); ok {
		if s, ok := f.GetWriter().(*syncWriter); ok {
			s.SetOutput(w)
		} else { // then we have to swap
			if clo, ok := old.(CloneableHandler); ok {
				if _, ok := old.(USOutputter); ok {
					new := clo.Clone()
					new.(USOutputter).UnsyncedSetOutput(w)
					h.SwapHandler(new)
				}
			}
		}
	}
}
