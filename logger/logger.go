package logger

import (
  "github.com/sirupsen/logrus"
)

var (
	log           *logrus.Logger
	AppLog        *logrus.Entry
	HttpLog       *logrus.Entry
)

func init() {
	log = logrus.New()
	AppLog = log.WithFields(logrus.Fields{"component": "middlewareApp", "category": "App"})
}