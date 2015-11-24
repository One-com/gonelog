package log

import (
	ilog "github.com/One-com/gonelog"
	"github.com/One-com/gonelog/syslog"
)

func // Log is the simplest Logger method
(l *Logger) Log(level syslog.Priority, msg string, kv ...interface{}) {
	if l.Does(level) {
		l.log(level, msg, kv...)
	}
}

//--- internal level loggers (to return from *ok() functions)
func (l *Logger) alert(msg string, kv ...interface{}) {
	level := syslog.LOG_ALERT
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) crit(msg string, kv ...interface{}) {
	level := syslog.LOG_CRIT
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) error(msg string, kv ...interface{}) {
	level := syslog.LOG_ERROR
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) warn(msg string, kv ...interface{}) {
	level := syslog.LOG_WARN
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) notice(msg string, kv ...interface{}) {
	level := syslog.LOG_NOTICE
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) info(msg string, kv ...interface{}) {
	level := syslog.LOG_INFO
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}
func (l *Logger) debug(msg string, kv ...interface{}) {
	level := syslog.LOG_DEBUG
	l.h.Log(l.newEvent(level, msg, normalize(kv)))
}

//---

// Methods which return a function which will do the queries logging when called.
func (l *Logger) ALERTok() (ilog.LogFunc, bool)  { return l.alert, l.Does(syslog.LOG_ALERT) }
func (l *Logger) CRITok() (ilog.LogFunc, bool)   { return l.crit, l.Does(syslog.LOG_CRIT) }
func (l *Logger) ERRORok() (ilog.LogFunc, bool)  { return l.error, l.Does(syslog.LOG_ERROR) }
func (l *Logger) WARNok() (ilog.LogFunc, bool)   { return l.warn, l.Does(syslog.LOG_WARN) }
func (l *Logger) NOTICEok() (ilog.LogFunc, bool) { return l.notice, l.Does(syslog.LOG_NOTICE) }
func (l *Logger) INFOok() (ilog.LogFunc, bool)   { return l.info, l.Does(syslog.LOG_INFO) }
func (l *Logger) DEBUGok() (ilog.LogFunc, bool)  { return l.debug, l.Does(syslog.LOG_DEBUG) }

//---

// ALERT logs a message and optinal KV values at ALERT level.
func (c *Logger) ALERT(msg string, kv ...interface{}) {
	l := syslog.LOG_ALERT
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) CRIT(msg string, kv ...interface{}) {
	l := syslog.LOG_CRIT
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) ERROR(msg string, kv ...interface{}) {
	l := syslog.LOG_ERROR
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) WARN(msg string, kv ...interface{}) {
	l := syslog.LOG_WARN
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) NOTICE(msg string, kv ...interface{}) {
	l := syslog.LOG_NOTICE
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) INFO(msg string, kv ...interface{}) {
	l := syslog.LOG_INFO
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
func (c *Logger) DEBUG(msg string, kv ...interface{}) {
	l := syslog.LOG_DEBUG
	if c.Does(l) {
		c.log(l, msg, kv...)
	}
}
