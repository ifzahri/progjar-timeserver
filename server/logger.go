package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	component string
	logger    *log.Logger
	level     LogLevel
}

func NewLogger(component string) *Logger {
	return &Logger{
		component: component,
		logger:    log.New(os.Stdout, "", 0),
		level:     LevelInfo,
	}
}

func NewLoggerOutput(component string, output io.Writer) *Logger {
	return &Logger{
		component: component,
		logger:    log.New(output, "", 0),
		level:     LevelInfo,
	}
}

func (l *Logger) Message(level LogLevel, message string, args ...interface{}) string {
	_, file, line, _ := runtime.Caller(3)
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]
	timestamp := time.Now().Format("15:04:05.000")
	msg := fmt.Sprintf(message, args...)

	return fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s",
		timestamp, level.String(), fileName, line, l.component, msg)
}

func (l *Logger) CheckLevel(level LogLevel) bool {
	return level >= l.level
}

func (l *Logger) Log(level LogLevel, message string, args ...interface{}) {
	if l.CheckLevel(level) {
		l.logger.Println(l.Message(level, message, args...))
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.Log(LevelDebug, message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.Log(LevelInfo, message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.Log(LevelWarn, message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.Log(LevelError, message, args...)
}
