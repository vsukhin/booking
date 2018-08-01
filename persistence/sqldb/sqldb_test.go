package sqldb

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/vsukhin/booking/logging"
)

func init() {
	logging.Log = NewFakeLogger()
}

// FakeLogger is fake logger
type FakeLogger struct {
	*logrus.Logger
}

// NewFakeLogger is a constructor of fake logger
func NewFakeLogger() logging.LoggerInterface {
	log := logrus.New()

	return &FakeLogger{log}
}

// Init initiates logging
func (logger *FakeLogger) Init(mode string) {
}

// WithFields logs with fields
func (logger *FakeLogger) WithFields(depthLevel int, fields logging.Fields) *logrus.Entry {
	return logrus.NewEntry(logger.Logger)
}

// Info logs info
func (logger *FakeLogger) Info(args ...interface{}) {
}

func Test_NewDB_Failure(t *testing.T) {
	db, err := NewDB("localhost:54321", nil, true, nil)
	if err == nil {
		t.Error("Expected to have db connection error")
	}
	if db != nil {
		t.Error("Expected to have nil db")
	}
}
