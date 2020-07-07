package roc

/*
#include <roc/log.h>

void rocGoLogHandlerProxy(roc_log_level level, char* component, char* message);
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
	// LogNone disables logging completely
	LogNone LogLevel = 0

	// LogError enables only error messages
	LogError LogLevel = 1

	// LogError enables informational messages and above
	LogInfo LogLevel = 2

	// LogDebug enables debugging messages and above
	LogDebug LogLevel = 3

	// LogDebug enables extra verbose logging, which may hurt performance
	// and should not be used in production
	LogTrace LogLevel = 4
)

// LogFunc is a handler for log messages.
//
// It is called for every message, if the corresponding log level is enabled.
//
// Its calls are serialized, so it doesn't need to be thread-safe.
type LoggerFunc func(level LogLevel, component string, message string)

// Logger interface is an alternative way to handle log messages.
//
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
func rocGoLogHandler(level C.roc_log_level, component *C.char, message *C.char) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	if loggerFunc != nil {
		loggerFunc(LogLevel(level), C.GoString(component), C.GoString(message))
	}
}

type defaultLogger struct{}

func (defaultLogger) Print(v ...interface{}) {
	log.Print(v...)
}

func makeLoggerFunc(logger Logger) LoggerFunc {
	return func(level LogLevel, component string, message string) {
		levStr := ""
		switch level {
		case LogError:
			levStr = "err"
		case LogInfo:
			levStr = "inf"
		case LogDebug:
			levStr = "dbg"
		case LogTrace:
			levStr = "trc"
		}
		logger.Print(fmt.Sprintf("[%s] %s: %s", levStr, component, message))
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
	if logFn == nil {
		logFn = makeLoggerFunc(defaultLogger{})
	}

	loggerMu.Lock()
	defer loggerMu.Unlock()

	loggerFunc = logFn
	C.roc_log_set_handler(C.roc_log_handler(C.rocGoLogHandlerProxy))
}

// SetLogger is like SetLoggerFunc, but uses Logger interface instead of LoggerFunc.
//
// If a nil Logger is passed, default logger is used, which passes all messages to
// the standard logger using log.Print.
//
// This function is thread-safe.
func SetLogger(logger Logger) {
	if logger == nil {
		logger = defaultLogger{}
	}

	SetLoggerFunc(makeLoggerFunc(logger))
}
