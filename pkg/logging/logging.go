package logging

import "github.com/sirupsen/logrus"

// GetLogger will get the configured logger
func GetLogger(loglevel string) *logrus.Entry {
	logger := logrus.NewEntry(logrus.New())
	lvl, err := logrus.ParseLevel(loglevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.Logger.SetLevel(lvl)
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
