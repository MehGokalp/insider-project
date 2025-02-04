package log

import "github.com/sirupsen/logrus"

func New() Logger {
	return &LogrusWrapper{
		logger: logrus.New(),
	}
}
