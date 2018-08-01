package models

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

func Test_CheckQueryTag_Success(t *testing.T) {
	seat := &Seat{}

	exists := CheckQueryTag("id", seat)
	if !exists {
		t.Error("Expected to check query tag successfully")
	}
}

func Test_CheckQueryTag_Failure(t *testing.T) {
	seat := &Seat{}

	exists := CheckQueryTag("test", seat)
	if exists {
		t.Error("Expected to have error checking query tag")
	}
}

func Test_GetSearchTag_Success(t *testing.T) {
	seat := &Seat{}

	tag := GetSearchTag("id", seat)
	if tag != "id" {
		t.Error("Expected to get search tag successfully")
	}
}

func Test_GetSearchTag_Failure(t *testing.T) {
	seat := &Seat{}

	tag := GetSearchTag("test", seat)
	if tag != "" {
		t.Error("Expected to have error getting search tag")
	}
}

func Test_GetAllSearchTags_Success(t *testing.T) {
	seat := &Seat{}

	tags := GetAllSearchTags(seat)
	if len(tags) != 8 {
		t.Error("Expected to get all search tags successfully")
	}
	for _, tag := range tags {
		if tag != "id" && tag != "index" && tag != "type" && tag != "row" &&
			tag != "line" && tag != "assigned" && tag != "created_at" && tag != "updated_at" {
			t.Error("Expected to get all known search tags successfully")
		}
	}
}
