package roc

import (
	"errors"
)

var (
	// ErrInvalidArguments indicates that one or more arguments passed to the function
	// are invalid
	ErrInvalidArguments = errors.New("One or more arguments are invalid")

	// ErrInvalidApi should never happen and indicates that the API don't follow the
	// declared contract
	ErrInvalidApi = errors.New("Invalid return code from API")

	// ErrRuntime indicates a runtime error: memory allocation error etc
	ErrRuntime = errors.New("Runtime error")
)
