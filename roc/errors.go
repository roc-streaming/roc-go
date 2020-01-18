package roc

/*
#cgo LDFLAGS: -lroc
#include <roc/receiver.h>
#include <roc/address.h>
#include <roc/sender.h>
#include <roc/log.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
)

var (
	// ErrInvalidArguments indicates that one or more arguments passed to the function
	// are invalid
	ErrInvalidArguments = errors.New("One or more arguments are invalid")

	// ErrInvalidApi should never happen and indicates that the API don't follow the declared contract
	ErrInvalidApi = errors.New("Invalid return code from API")

	// ErrRuntime indicates a runtime error: memory allocation error etc
	ErrRuntime = errors.New("Runtime error")
)

// safeString ensures that the string is NULL-terminated, a NULL-terminated copy is created otherwise.
func safeString(str string) string {
	if len(str) > 0 && str[len(str)-1] != '\x00' {
		str = str + "\x00"
	} else if len(str) == 0 {
		str = "\x00"
	}
	return str
}
