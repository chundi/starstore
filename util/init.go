package util

import (
	"github.com/sirupsen/logrus"
	"github.com/galaxy-solar/starstore/log"
)

var Logger *logrus.Logger

func init() {
	Logger = log.InitDefaultLogger()
}
