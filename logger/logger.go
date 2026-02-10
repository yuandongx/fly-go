// Package log provides the logging functionality for the application.
package log

import (
	"fmt"
	"log"
)

const (
	DEBUG = 1 << iota
	INFO
	ERROR
	PANIC
)

type Zip struct {
	Key   string
	Value interface{}
}
type ILogger struct {
	Level  int
	Logger *log.Logger
}

func Zap(key string, value interface{}) Zip {
	return Zip{key, value}
}

func (z Zip) String() string {
	return fmt.Sprintf("%s=%v", z.Key, z.Value)
}

func DefaultLogger() *ILogger {
	logger := &ILogger{}
	l := log.Default()
	l.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	logger.Level = DEBUG
	logger.Logger = l
	return logger
}
func (l *ILogger) Println(msg string, items ...Zip) {
	tmp := ""
	for _, item := range items {
		tmp += fmt.Sprintf("%s %s, ", item.Key, item.Value)
	}
	tmp += msg
	l.Logger.Println(tmp)
}
func (l *ILogger) Debug(msg string, items ...Zip) {
	if l.Level >= DEBUG {
		l.Println(msg, items...)
	}
}
func (l *ILogger) Info(msg string, items ...Zip) {
	if l.Level >= INFO {
		l.Println(msg, items...)
	}
}
func (l *ILogger) Error(msg string, items ...Zip) {
	if l.Level >= ERROR {
		l.Println(msg, items...)
	}
}
func (l *ILogger) Panic(msg string, items ...Zip) {
	if l.Level >= PANIC {
		l.Println(msg, items...)
	}
}
