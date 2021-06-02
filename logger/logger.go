package logger

import (
	"fmt"
	"reflect"
	"time"
)

type LogLevel int

const (
	INFO  LogLevel = iota
	DEBUG          = iota
	ERROR          = iota
	WARN           = iota
	FATAL          = iota
)

// LogLevelMap is a maaping for LogLevel enum
var logLevelMap = map[int]string{
	int(INFO):  "INFO",
	int(DEBUG): "DEBUG",
	int(ERROR): "ERROR",
	int(WARN):  "WARN",
	int(FATAL): "FATAL",
}

// Log is the type of object that can be logged in the LogHistory.
type Log struct {
	LogLevel   LogLevel
	ModuleName string
	TimeStamp  string
	Text       string
	Data       map[string]interface{}
}

// Node is a node for implementing linkedlist in LogHistory.
type Node struct {
	prev *Node
	next *Node
	log  Log
}

// LogHistory is the object for recording all logs in linkedlist fashion.
type LogHistory struct {
	head *Node
	tail *Node
}

// Logger is an struct for logging.
type Logger struct {
	LogHistory *LogHistory
}

// NewLogger will create new Logger instance.
func NewLogger() Logger {
	return Logger{
		LogHistory: &LogHistory{},
	}
}

// structToMap converts struct to map[string]interface{}.
func structToMap(in interface{}, tag string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	mapVal := reflect.ValueOf(in)
	if mapVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid input struct")
	}
	if mapVal.Kind() == reflect.Ptr {
		mapVal = mapVal.Elem()
	}

	for i := 0; i < mapVal.Type().NumField(); i++ {
		field := mapVal.Type().Field(i)
		if tagVal := field.Tag.Get(tag); tagVal != "" {
			ret[tagVal] = mapVal.Field(i).Interface()
		}
	}
	return ret, nil
}

// Log will insert new log (data) into log history in a linked-list fashion.
// It will panic if given a isStruct = true value and the data value fed isn't a struct.
// The data parameter expects a map[string]interface{} unless stated via isStruct.
func (lgr Logger) Log(logLvl LogLevel, txt string, modName string, data interface{}, dataTag string, isStruct bool) {
	var dataMap map[string]interface{}
	if isStruct {
		var err interface{}
		dataMap, err = structToMap(data, "json")
		if err != nil {
			panic("Invalid Log Data. Data should be of type Struct")
		}
	} else {
		dataMap = data.(map[string]interface{})
	}
	list := &Node{
		next: lgr.LogHistory.head,
		log: Log{
			LogLevel:   logLvl,
			TimeStamp:  time.Now().Format(time.RFC850),
			ModuleName: modName,
			Text:       txt,
			Data:       dataMap,
		},
	}
	if lgr.LogHistory.head != nil {
		lgr.LogHistory.head.prev = list
	}
	lgr.LogHistory.head = list
	l := lgr.LogHistory.head
	for l.next != nil {
		l = l.next
	}
	lgr.LogHistory.tail = l
}

// DisplayLogs will display all saved logs using fmt.Printf
func (lgr Logger) DisplayLogs() {
	list := lgr.LogHistory.head
	for list != nil {
		fmt.Printf(
			"%v [%v]<%v> %v\n",
			list.log.TimeStamp,
			logLevelMap[int(list.log.LogLevel)],
			list.log.ModuleName,
			list.log.Text,
		)
		list = list.next
	}
}
