package roc

import (
	"errors"
)

var (
	// ErrInvalidArgs indicates that one or more function arguments are invalid
	ErrInvalidArgs = errors.New("One or more arguments are invalid")
)
