package roc

import "fmt"

type logger interface {
	Print(v ...interface{})
}

var handler LogHandler

func LogSetLevel(level LogLevel) {
	logSetLevelImpl(level)
}

func LogSetHandler(l logger) {
	handler = func(level LogLevel, component string, message string) {
		l.Print(fmt.Sprintf("%s: %s", component, message))
	}

	logSetHandlerImpl(handler)
}
