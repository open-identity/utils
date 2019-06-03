package logrusx

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	ViperKeyLogLevel  = "log.level"
	ViperKeyLogFormat = "log.format"
)

// New initializes logrus with environment variable configuration LOG_LEVEL and LOG_FORMAT.
func New() *logrus.Logger {
	l := logrus.New()
	ll, err := logrus.ParseLevel(viper.GetString(ViperKeyLogLevel))
	if err != nil {
		ll = logrus.InfoLevel
	}
	l.Level = ll

	if viper.GetString(ViperKeyLogFormat) == "json" {
		l.Formatter = new(logrus.JSONFormatter)
	}

	return l
}
