package logger

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// setup log file for rotate and max size
func SetupLogger() {
	logger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   "rest4smpp.log",
		MaxSize:    5, // MB
		MaxBackups: 10,
		Compress:   true, // disabled by default
	}

	logFormatter := &logrus.TextFormatter{
		TimestampFormat:           time.RFC3339,
		FullTimestamp:             true,
		EnvironmentOverrideColors: true,
		PadLevelText:              true,
	}

	logrus.SetFormatter(logFormatter)
	logrus.SetOutput(io.Writer(logger))
}
