package syslog

type Priority int

const (
	LOG_EMERG Priority = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

// aliases

const (
	LOG_ERROR Priority = LOG_ERR
	LOG_WARN  Priority = LOG_WARNING
)
