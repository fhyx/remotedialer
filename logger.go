package remotedialer

import (
	"fmt"
	syslog "log"
)

// dftLogger default instance
var dftLogger Logger

func init() {
	syslog.SetFlags(syslog.Ltime | syslog.Lshortfile)
	dftLogger = &logImpl{}
}

// SetLogger ...
func SetLogger(logger Logger) {
	if logger != nil {
		dftLogger = logger
	}
}

// GetLogger ...
func GetLogger() Logger {
	return dftLogger
}

// Logger like zap.Sugar
type Logger interface {

	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})
	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})
	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})
	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})
}

type logImpl struct{}

func (z *logImpl) Debugf(template string, args ...interface{}) {
	// syslog.Output(2, fmt.Sprintf(template, args...))
}

func (z *logImpl) Infof(template string, args ...interface{}) {
	syslog.Output(2, fmt.Sprintf(template, args...))
}

func (z *logImpl) Warnf(template string, args ...interface{}) {
	syslog.Output(2, fmt.Sprintf(template, args...))
}

func (z *logImpl) Errorf(template string, args ...interface{}) {
	syslog.Output(2, fmt.Sprintf(template, args...))
}
