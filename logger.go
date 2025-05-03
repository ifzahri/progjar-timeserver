// logger.go
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	component string
	logger    *log.Logger
}

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

func NewLogger(component string) *Logger {
	return &Logger{
		component: component,
		logger:    log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Message(level string, message string, args ...interface{}) string {
	_, file, line, _ := runtime.Caller(2)
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]
	timestamp := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(message, args...)

	return fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s",
		timestamp, level, fileName, line, l.component, msg)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Println(l.Message("DEBUG", message, args...))
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Println(l.Message("INFO", message, args...))
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Println(l.Message("WARN", message, args...))
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Println(l.Message("ERROR", message, args...))
}
