package logger

import (
	"log"

	"github.com/sirupsen/logrus"
)

type Logger struct{}

func New() *Logger {
	logLevel, err := logrus.ParseLevel("debug")
	if err != nil {
		log.Fatalf("failed to parse the level: %v", err)
	}

	logrus.SetLevel(logLevel)

	return &Logger{}
}

func (l Logger) Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func (l Logger) Error(msg ...interface{}) {
	logrus.Error(msg...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
