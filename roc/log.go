package roc

/*
#include <roc/log.h>

unsigned long long rocGoThreadID();
void rocGoLogHandlerProxy(const roc_log_message* message, void* argument);
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// LogLevel defines the logging verbosity.
//
//go:generate stringer -type LogLevel -trimprefix Log
type LogLevel int

const (
	// LogNone disables logging completely.
	LogNone LogLevel = 0

	// LogError enables only error messages.
	LogError LogLevel = 1

	// LogInfo enables informational messages and above.
	LogInfo LogLevel = 2

	// LogDebug enables debugging messages and above.
	LogDebug LogLevel = 3

	// LogTrace enables extra verbose logging, which may hurt performance
	// and should not be used in production.
	LogTrace LogLevel = 4
)

// LogMessage defines message written to log.
type LogMessage struct {
	// Message log level.
	Level LogLevel

	// Name of the module that originated the message.
	Module string

	// Name of the source code file.
	// May be empty.
	File string

	// Line number in the source code file.
	Line int

	// Message timestamp, unix time at which message was logged
	Time time.Time

	// Platform-specific process ID.
	Pid uint64

	// Platform-specific thread ID.
	Tid uint64

	// Message text.
	Text string
}

// LogFunc is a handler for log messages.
// It is called for every message, if the corresponding log level is enabled.
// Its calls are serialized, so it doesn't need to be thread-safe.
type LoggerFunc func(LogMessage)

// Logger interface is an alternative way to handle log messages.
// It is used like LoggerFunc, but receives a single string with formatted message.
// This interface is compatible with log.Logger from standard library.
type Logger interface {
	Print(v ...interface{})
}

// SetLogLevel changes the logging level.
//
// Messages with higher verbosity than the given level will be dropped.
// Default log level is LogError.
//
// This function is thread-safe.
func SetLogLevel(level LogLevel) {
	checkVersionFn()

	atomic.StoreInt32(&loggerLevel, int32(level))
	C.roc_log_set_level(C.roc_log_level(level))
}

// SetLoggerFunc sets the handler for log messages.
//
// Starting from this call, all log messages produced by the library, will be passed
// to the given function. It may be called from different threads, but the calls will
// be always serialized, so it doesn't need to be thread-safe.
//
// If a nil function is passed, default logger is used, which passes all messages to
// the standard logger using log.Print.
//
// This function is thread-safe.
func SetLoggerFunc(logFn LoggerFunc) {
	checkVersionFn()

	if logFn == nil {
		logFn = logger2func(standardLogger{})
	}

	loggerFunc.Store(logFn)
}

// SetLogger is like SetLoggerFunc, but uses Logger interface instead of LoggerFunc.
//
// If a nil Logger is passed, default logger is used, which passes all messages to
// the standard logger using log.Print.
//
// This function is thread-safe.
func SetLogger(logger Logger) {
	checkVersionFn()

	if logger == nil {
		logger = standardLogger{}
	}

	loggerFunc.Store(logger2func(logger))
}

func logger2func(logger Logger) LoggerFunc {
	return func(message LogMessage) {
		level := ""
		switch message.Level {
		case LogError:
			level = "err"
		case LogInfo:
			level = "inf"
		case LogDebug:
			level = "dbg"
		case LogTrace:
			level = "trc"
		}
		logger.Print(fmt.Sprintf("[%s] %s: %s", level, message.Module, message.Text))
	}
}

type standardLogger struct{}

func (standardLogger) Print(v ...interface{}) {
	log.Print(v...)
}

var (
	loggerLevel int32
	loggerFunc  atomic.Value
	loggerChan  = make(chan LogMessage, 1024)
)

// Write structured message to log.
// Invoked from C library when it needs to log something.
//
//export rocGoLogHandler
func rocGoLogHandler(cMessage *C.roc_log_message) {
	message := LogMessage{
		Level: LogLevel(cMessage.level),
		Time:  time.Unix(int64(cMessage.time), 0),
		Pid:   uint64(cMessage.pid),
		Tid:   uint64(cMessage.tid),
	}
	if cMessage.module != nil {
		message.Module = C.GoString(cMessage.module)
	}
	if cMessage.file != nil {
		message.File = C.GoString(cMessage.file)
		message.Line = int(cMessage.line)
	}
	if cMessage.text != nil {
		message.Text = C.GoString(cMessage.text)
	}

	loggerChan <- message
}

// Write formatted message to log.
// Invoked from Go code when it needs to log something.
func logWrite(level LogLevel, text string, params ...interface{}) {
	if level > logLevel() {
		return
	}

	file, line := logLocation(2)

	message := LogMessage{
		Level:  level,
		Time:   time.Now(),
		Pid:    uint64(os.Getpid()),
		Tid:    uint64(C.rocGoThreadID()),
		Module: "roc_go",
		File:   file,
		Line:   line,
		Text:   fmt.Sprintf(text, params...),
	}

	loggerChan <- message
}

func logLevel() LogLevel {
	return LogLevel(atomic.LoadInt32(&loggerLevel))
}

func logLocation(stack int) (string, int) {
	_, file, line, ok := runtime.Caller(stack)
	if !ok {
		return "", -1
	}

	parts := strings.FieldsFunc(file, func(c rune) bool {
		return c == '/' || c == '\\'
	})
	if len(parts) == 0 {
		return "", -1
	}

	for i := len(parts) - 1; i != 0; i-- {
		if parts[i] == "roc" {
			parts = parts[i:]
			break
		}
	}

	return filepath.Join(parts...), line
}

func logRoutine() {
	for message := range loggerChan {
		fn := loggerFunc.Load().(LoggerFunc)
		if fn != nil {
			fn(message)
		}
	}
}

func init() {
	SetLogLevel(LogError)
	SetLoggerFunc(nil)

	// rocGoLogHandlerProxy calls rocGoLogHandler,
	// rocGoLogHandler writes messages to channel
	C.roc_log_set_handler(C.roc_log_handler(C.rocGoLogHandlerProxy), nil)

	// logRoutine reads messages from channel and passes them to
	// Logger or LoggerFunc
	go logRoutine()
}
