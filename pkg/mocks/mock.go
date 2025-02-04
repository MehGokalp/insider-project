package mocks

import (
	"github.com/stretchr/testify/mock"
)

func NewLoggerMock(allowedFuncList ...string) *LoggerMock {
	logger := LoggerMock{}

	if len(allowedFuncList) == 0 {
		allowedFuncList = []string{"Infof", "Debugf"}
	}

	for _, f := range allowedFuncList {
		logger.On(
			f,
			mock.AnythingOfType("string"),
			mock.Anything,
		).Maybe()
	}

	return &logger
}

type LoggerMock struct {
	mock.Mock
}

func (logger *LoggerMock) Debugf(format string, args ...interface{}) {
	logger.Called(format, args)
}

func (logger *LoggerMock) Infof(format string, args ...interface{}) {
	logger.Called(format, args)
}

func (logger *LoggerMock) Errorf(format string, args ...interface{}) {
	logger.Called(format, args)
}
