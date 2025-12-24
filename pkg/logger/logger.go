package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(serviceName string) {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	Log.SetOutput(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
		Log.WithField("log_level", logLevel).Warn("Invalid log level, using default 'info'")
	}

	Log.SetLevel(level)

	Log.WithField("service", serviceName)

	Log.WithFields(logrus.Fields{
		"level":   level.String(),
		"service": serviceName,
	}).Info("Logger initialized")
}
