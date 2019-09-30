package micro

import (
	"fmt"
	"github.com/ProtocolONE/go-core/logger"
	logIface "github.com/go-log/log"
)

type logAdapter struct {
	redirectLevel logger.Level
	logFunc       func(level logger.Level, format string, o ...logger.Option)
}

// Log
func (l *logAdapter) Log(v ...interface{}) {
	l.logFunc(l.redirectLevel, "%v", logger.Args(fmt.Sprint(v...)))
}

// Logf
func (l *logAdapter) Logf(format string, v ...interface{}) {
	l.logFunc(l.redirectLevel, format, logger.Args(v...))
}

// NewLoggerAdapter
func NewLoggerAdapter(log logger.Logger, lvl logger.Level) logIface.Logger {
	return &logAdapter{
		redirectLevel: lvl,
		logFunc:       log.Log,
	}
}
