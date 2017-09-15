package log

import (
	"fmt"
	"github.com/galaxy-solar/starstore/conf"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
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

func NewLogger(fmtType string, level string, logFile string) *logrus.Logger {
	var formatter logrus.Formatter
	if strings.ToLower(fmtType) == "json" {
		formatter = &logrus.JSONFormatter{}
	} else {
		formatter = &logrus.TextFormatter{
			FullTimestamp:    true,
		}
	}
	output, err := OutPutWriter(logFile)
	if err != nil {
		fmt.Println("Open log file error!!", err)
		os.Exit(0)
	}
	return &logrus.Logger{
		Formatter: formatter.(logrus.Formatter),
		Level:     ParseLogLevel(level),
		Out:       output,
	}
}

func ParseLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

func OutPutWriter(logFile string) (file io.Writer, err error) {
	if logFile == "" {
		return os.Stdout, nil
	}
	file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0666)
	return file, err
}
