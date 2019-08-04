package logger

import (
	"fmt"
	"log"
)

const (
	KLogLevelDebug = iota
	KLogLevelInfo
	KLogLevelWarn
	KLogLevelError
	KLogLevelFatal
)

var logLevelPrefix = []string{
	"[DBG] ",
	"[INF] ",
	"[WRN] ",
	"[ERR] ",
	"[FAL]",
}

// ILogger is an interface use for log message
type ILogger interface {
	Output(level int, calldepth int, f string) error
}

// Default logger
type defaultLogger struct {
}

func (d *defaultLogger) Output(level int, calldepth int, f string) error {
	text := logLevelPrefix[level] + "[tcpnetwork]" + f
	return log.Output(calldepth, text)
}

// Global variables
var myDefaultLogger defaultLogger
var myLogger ILogger = &myDefaultLogger

// SetLogger Set the custom logger
func SetLogger(logger ILogger) {
	myLogger = logger
}

func _log(level int, calldepth int, f string) {
	myLogger.Output(level, calldepth, f)
}

func LogDebug(f string, v ...interface{}) {
	_log(KLogLevelDebug, 2, fmt.Sprintf(f, v...))
}

func LogInfo(f string, v ...interface{}) {
	_log(KLogLevelInfo, 2, fmt.Sprintf(f, v...))
}

func LogWarn(f string, v ...interface{}) {
	_log(KLogLevelWarn, 2, fmt.Sprintf(f, v...))
}

func LogError(f string, v ...interface{}) {
	_log(KLogLevelError, 2, fmt.Sprintf(f, v...))
}

func LogFatal(f string, v ...interface{}) {
	_log(KLogLevelFatal, 2, fmt.Sprintf(f, v...))
}
