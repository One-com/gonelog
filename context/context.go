package context

import (
	ilog "github.com/One-com/gonelog"
	"github.com/One-com/gonelog/log"
	goctx "golang.org/x/net/context"
	"time"
)

type Context interface {
	goctx.Context
	ilog.Logger
}

type context struct {
	goctx.Context
	*log.Logger
}

var background = &context{goctx.Background(), log.Default()}

func Background() *context {
	return background
}

func (c *context) GoContext() goctx.Context {
	return c.Context
}

func WithCancel(parent goctx.Context) (Context, goctx.CancelFunc) {
	var gc goctx.Context
	var l *log.Logger
	var cf goctx.CancelFunc
	gc, cf = goctx.WithCancel(parent)
	if pctx, ok := parent.(*context); ok {
		l = pctx.Logger
	}
	return &context{gc, l}, cf
}

func WithTimeout(parent goctx.Context, timeout time.Duration) (Context, goctx.CancelFunc) {
	var gc goctx.Context
	var l *log.Logger
	var cf goctx.CancelFunc
	gc, cf = goctx.WithTimeout(parent, timeout)
	if pctx, ok := parent.(*context); ok {
		l = pctx.Logger
	}
	return &context{gc, l}, cf
}

func WithValue(parent goctx.Context, key interface{}, val interface{}) Context {
	var gc goctx.Context
	var l *log.Logger
	gc = goctx.WithValue(parent, key, val)
	if pctx, ok := parent.(*context); ok {
		l = pctx.Logger
	}
	return &context{gc, l}
}

func WithLogging(parent goctx.Context, l *log.Logger) Context {
	return &context{parent, l}
}

func WithLoggedValue(parent goctx.Context, key interface{}, val interface{}) Context {
	if p, ok := parent.(*context); ok {
		return &context{parent, p.Logger.With(key, val)}
	} else {
		return &context{parent, log.Default().With(key, val)}
	}
}
