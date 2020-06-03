package roc

/*
#include <roc/log.h>
*/
import "C"

// LogLevel as declared in roc/log.h:53
type LogLevel int32

// LogLevel enumeration from roc/log.h:53
const (
	LogNone  LogLevel = iota
	LogError LogLevel = 1
	LogInfo  LogLevel = 2
	LogDebug LogLevel = 3
	LogTrace LogLevel = 4
)

// LogHandler type as declared in roc/log.h:64
type LogHandler func(level LogLevel, component string, message string)
