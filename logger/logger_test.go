package logger

import (
	"fmt"
	"testing"
)

type TestObj struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
	Param3 string
}

func TestLogging(t *testing.T) {
	logger := NewLogger()
	logger.LogTxt(INFO, "Test info log")
	logger.LogTxt(DEBUG, "Test debug log")
	logger.LogTxt(ERROR, "Test error log")
	logger.LogTxt(WARN, "Test warning log")
	logger.LogTxt(FATAL, "Test fatal log")

	logger.LogObj(
		INFO,
		"Test log info with app data map[string]interface{}",
		map[string]interface{}{"test": "test"},
		"",
		false,
	)
	logger.LogObj(
		WARN,
		"Test log warning with app data TestObj struct",
		TestObj{
			Param1: "test1",
			Param2: "test2",
			Param3: "test3",
		},
		"json",
		true,
	)

	logger.DisplayLogsForward()
	logger.DisplayLogsBackward()
	fmt.Println("run without errors")
}

func TestStructToMapMethod(t *testing.T) {
	covertedMap, err := structToMap(TestObj{
		Param1: "testvalue1",
		Param2: "testvalue2",
		Param3: "testvalue3",
	}, "json")
	if err != nil || covertedMap["param1"] != "testvalue1" || covertedMap["param2"] != "testvalue2" {
		t.Errorf("struct to map failed")
	}
}

func TestStructToMapInvalidMethod(t *testing.T) {
	_, err := structToMap("test invalid type", "json")
	if err == nil {
		t.Errorf("invalid input struct not raised")
	}
}

func TestLoggerInvalidStruct(t *testing.T) {
	logger := NewLogger()
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("invalid struct checking failed")
		}
	}()
	logger.LogObj(INFO, "test", &TestObj{}, "test", true)
}
