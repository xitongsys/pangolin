package logging

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	Log.SetFormatter(formatter)
	Log.SetReportCaller(false)
	Log.SetLevel(logrus.InfoLevel)
}

func SetLevel(level string) {
	if level == "debug" {
		Log.SetLevel(logrus.DebugLevel)
	} else if level == "warn" {
		Log.SetLevel(logrus.WarnLevel)
	} else if level == "info" {
		Log.SetLevel(logrus.InfoLevel)
	}
}
