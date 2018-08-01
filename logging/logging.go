package logging

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	// ModeDev is development mode
	ModeDev = "dev"
	// ModeStaging is staging mode
	ModeStaging = "staging"
	// ModeProd is production mode
	ModeProd = "prod"
)

const (
	// DepthLow is low depth
	DepthLow = 1
	// DepthModerate is moderate depth
	DepthModerate = 2
	// DepthHigh is high depth
	DepthHigh = 3
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

var (
	// Log is logger instance
	Log LoggerInterface
)

// Logger is logger
type Logger struct {
	*logrus.Logger
}

// LoggerInterface is logger interface
type LoggerInterface interface {
	Init(mode string)
	WithFields(depthLevel int, fields Fields) *logrus.Entry
	Info(args ...interface{})
}

// NewLogger is a constructor of logger
func NewLogger() LoggerInterface {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	return &Logger{log}
}

// Init initiates logging
func (logger *Logger) Init(mode string) {
	switch mode {
	case ModeDev:
		logger.Logger.SetLevel(logrus.DebugLevel)
	case ModeStaging, ModeProd:
		logger.Logger.SetLevel(logrus.InfoLevel)
	default:
		logger.Logger.SetLevel(logrus.DebugLevel)
		logger.WithFields(DepthLow, Fields{
			"mode": mode,
		}).Warn("Unexpected mode value")
	}

	logger.WithFields(DepthLow, Fields{
		"level": logrus.GetLevel().String(),
	}).Info("Service logging level")
}

// WithFields logs with fields
func (logger *Logger) WithFields(depthLevel int, fields Fields) *logrus.Entry {
	_, file, line, ok := runtime.Caller(depthLevel)
	if ok {
		fields["line"] = fmt.Sprintf("%s:%d", file, line)
	}

	host, err := os.Hostname()
	if err != nil {
		fields["host_error"] = err
	} else {
		fields["host"] = host
	}

	return logger.Logger.WithFields(logrus.Fields(fields))
}

// Info logs info
func (logger *Logger) Info(args ...interface{}) {
	logger.Logger.Info(args...)
}
