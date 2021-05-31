package logger

import (
	"fmt"
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
func NewLogger() *Logger {
	return &Logger{
		LogHistory: &LogHistory{},
	}
}

/* Insert new log in loghistory */
func (lgr Logger) Log(logLvl LogLevel, txt string, modName string, data map[string]interface{}) {
	list := &Node{
		next: lgr.LogHistory.head,
		log: Log{
			LogLevel:   logLvl,
			TimeStamp:  time.Now().Format(time.RFC850),
			ModuleName: modName,
			Text:       txt,
			Data:       data,
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
	fmt.Println(lgr.LogHistory)
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
