package healthx

import (
	log "github.com/InVisionApp/go-logger"
	"github.com/sirupsen/logrus"
)

type shim struct {
	logrus.FieldLogger
}

// NewLogrus can be used to override the default logger.
// Optionally pass in an existing logrus logger or pass in
// `nil` to use the default logger.
func NewShim(logger logrus.FieldLogger) log.Logger {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	return &shim{FieldLogger: logger}
}

// WithFields will return a new logger based on the original logger with
// the additional supplied fields. Wrapper for logrus Entry.WithFields()
func (s *shim) WithFields(fields log.Fields) log.Logger {
	cp := &shim{
		s.FieldLogger.WithFields(logrus.Fields(fields)),
	}
	return cp
}
