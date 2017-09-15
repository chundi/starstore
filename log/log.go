package log

import (
	"github.com/sirupsen/logrus"
	"github.com/galaxy-solar/starstore/conf"
)

func InitDefaultLogger() *logrus.Logger {
	logger := logrus.New()
	if conf.IsDevelopMode() {
		logger.Formatter = &logrus.TextFormatter{}
	} else {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.SetLevel(logrus.DebugLevel)
	return logger
}