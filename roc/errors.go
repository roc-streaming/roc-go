package roc

import (
	"errors"
)

var (
	// ErrInvalidArguments indicates that one or more arguments
	// passed to the function are invalid
	ErrInvalidArguments = errors.New("One or more arguments are invalid")
)
