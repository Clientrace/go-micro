package logger

import (
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
	LogHistory LogHistory
}

/* Insert Item into linked list */
func (L *LogHistory) insert(logLvl LogLevel, txt string, modName string, data map[string]interface{}) {
	list := &Node{
		next: L.head,
		log: Log{
			LogLevel:   logLvl,
			TimeStamp:  time.Now().Format(time.RFC850),
			ModuleName: modName,
			Text:       txt,
			Data:       data,
		},
	}
	if L.head != nil {
		L.head.prev = list
	}

	L.head = list

	l := L.head
	for l.next != nil {
		l = l.next
	}
	L.tail = l
}

/* Create Logger Instance */
func NewLogger() Logger {
	lh := LogHistory{}
	return Logger{
		LogHistory: lh,
	}
}

func (lgr Logger) Log(logLvl LogLevel, txt string, modName string, data map[string]interface{}) {
	lgr.LogHistory.insert(logLvl, txt, modName, data)
}
