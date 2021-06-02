package logger

import (
	"testing"
)

func TestLogging(t *testing.T) {
	logger := NewLogger()
	logger.LogTxt(INFO, "Test log", "TestLogging")
	logger.LogTxt(INFO, "Test log", "TestLogging")
	logger.DisplayLogs()

	// lh := LogHistory{}
	// lh.insert(INFO, "TEST", "service", nil)
	// lh.insert(INFO, "TEST2", "service", nil)
	// lh.insert(INFO, "TEST3", "service", nil)
	// fmt.Println(lh.head)
}
