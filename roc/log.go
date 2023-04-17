package roc

/*
#include <roc/log.h>

void rocGoLogHandlerProxy(const roc_log_message* message, void* argument);
*/
import "C"

import (
	"fmt"
	"log"
	"sync"
)

// LogLevel defines the logging verbosity.
type LogLevel int

const (
	// LogNone disables logging completely.
	LogNone LogLevel = 0

	// LogError enables only error messages.
	LogError LogLevel = 1

	// LogError enables informational messages and above.
	LogInfo LogLevel = 2

	// LogDebug enables debugging messages and above.
	LogDebug LogLevel = 3

	// LogDebug enables extra verbose logging, which may hurt performance
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

	// Message timestamp, nanoseconds since Unix epoch.
	Time uint64

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

var (
	loggerFunc LoggerFunc
	loggerMu   sync.Mutex
)

//export rocGoLogHandler
func rocGoLogHandler(cMessage *C.roc_log_message) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	if loggerFunc != nil {
		message := LogMessage{
			Level: LogLevel(cMessage.level),
			Time:  uint64(cMessage.time),
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
		loggerFunc(message)
	}
}

type defaultLogger struct{}

func (defaultLogger) Print(v ...interface{}) {
	log.Print(v...)
}

func makeLoggerFunc(logger Logger) LoggerFunc {
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

func init() {
	SetLoggerFunc(nil)
}

// SetLogLevel changes the logging level.
//
// Messages with higher verbosity than the given level will be dropped.
// Default log level is LogError.
//
// This function is thread-safe.
func SetLogLevel(level LogLevel) {
	versionCheckFn()

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
	versionCheckFn()

	if logFn == nil {
		logFn = makeLoggerFunc(defaultLogger{})
	}

	loggerMu.Lock()
	defer loggerMu.Unlock()

	loggerFunc = logFn
	C.roc_log_set_handler(C.roc_log_handler(C.rocGoLogHandlerProxy), nil)
}

// SetLogger is like SetLoggerFunc, but uses Logger interface instead of LoggerFunc.
//
// If a nil Logger is passed, default logger is used, which passes all messages to
// the standard logger using log.Print.
//
// This function is thread-safe.
func SetLogger(logger Logger) {
	versionCheckFn()

	if logger == nil {
		logger = defaultLogger{}
	}

	SetLoggerFunc(makeLoggerFunc(logger))
}
