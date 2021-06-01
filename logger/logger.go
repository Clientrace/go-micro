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

var logLevelMap = map[int]string{
	int(INFO):  "INFO",
	int(DEBUG): "DEBUG",
	int(ERROR): "ERROR",
	int(WARN):  "WARN",
	int(FATAL): "FATAL",
}

/* Log Object */
type Log struct {
	LogLevel   LogLevel
	ModuleName string
	TimeStamp  string
	Text       string
	Data       map[string]interface{}
}

/* Linked List Node */
type Node struct {
	prev *Node
	next *Node
	log  Log
}

/* Linked List */
type LogHistory struct {
	head *Node
	tail *Node
}

type Logger struct {
	LogHistory *LogHistory
}

/* Create Logger Instance */
func NewLogger() Logger {
	return Logger{
		LogHistory: &LogHistory{},
	}
}

// structToMap - Convert struct to map[string]interface{}
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

/* Insert new log in loghistory */
func (lgr Logger) Log(logLvl LogLevel, txt string, modName string, data interface{}, dataTag string) {
	dataMap, err := structToMap(data, "json")
	if err != nil {
		panic("Invalid Log Data. Data should be of type Struct")
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
