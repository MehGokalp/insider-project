package log

import "github.com/sirupsen/logrus"

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type LogrusWrapper struct {
	logger logrus.FieldLogger
}

func (l *LogrusWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *LogrusWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *LogrusWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
