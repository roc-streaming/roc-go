package roc

import "fmt"

type Logger interface {
	Print(v ...interface{})
}

func LogSetLevel(level LogLevel) {
	logSetLevel(level)
}

func LogSetHandler(logger Logger) {
	logSetHandler(func(level LogLevel, component string, message string) {
		logger.Print(fmt.Sprintf("%s: %s", component, message))
	})
}
