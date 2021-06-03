package logger

import (
	"testing"
)

func TestLogging(t *testing.T) {
	logger := NewLogger()
	logger.LogTxt(INFO, "Test log")
	logger.LogTxt(INFO, "Test log")
	logger.DisplayLogsForward()
	logger.DisplayLogsBackward()
}
