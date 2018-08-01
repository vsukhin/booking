package logging

import (
	"testing"
)

func Test_Logger_Init_Dev_Success(t *testing.T) {
	logger := NewLogger()

	logger.Init(ModeDev)
}

func Test_Logger_Init_Staging_Success(t *testing.T) {
	logger := NewLogger()

	logger.Init(ModeStaging)
}

func Test_Logger_Init_Prod_Success(t *testing.T) {
	logger := NewLogger()

	logger.Init(ModeProd)
}

func Test_Logger_Init_Unknown_Success(t *testing.T) {
	logger := NewLogger()

	logger.Init("Unknown")
}

func Test_Logger_WithFields_Success(t *testing.T) {
	logger := NewLogger()

	var fields = Fields{}

	entry := logger.WithFields(DepthModerate, fields)
	if entry == nil {
		t.Error("Expected to log fields successfully")
	}
	if fields["line"] == "" {
		t.Error("Expected not empty line")
	}
	if fields["host"] == "" {
		t.Error("Expected not empty host")
	}
}

func Test_Logger_Info_Success(t *testing.T) {
	logger := NewLogger()

	logger.Info("test")
}
