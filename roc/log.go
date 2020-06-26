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

type LogLevel int

const (
	LogNone  LogLevel = 0
	LogError LogLevel = 1
	LogInfo  LogLevel = 2
	LogDebug LogLevel = 3
	LogTrace LogLevel = 4
)

type LoggerFunc func(level LogLevel, component string, message string)

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

func SetLogLevel(level LogLevel) {
	C.roc_log_set_level(C.roc_log_level(level))
}

func SetLoggerFunc(logFn LoggerFunc) {
	if logFn == nil {
		logFn = makeLoggerFunc(defaultLogger{})
	}

	loggerMu.Lock()
	defer loggerMu.Unlock()

	loggerFunc = logFn
	C.roc_log_set_handler(C.roc_log_handler(C.rocGoLogHandlerProxy))
}

func SetLogger(logger Logger) {
	if logger == nil {
		logger = defaultLogger{}
	}

	SetLoggerFunc(makeLoggerFunc(logger))
}
