package roc

import (
	"C"
	"fmt"
)

type nativeErr struct {
	op   string
	code int
}

func newNativeErr(op string, code C.int) nativeErr {
	return nativeErr{
		op:   op,
		code: int(code),
	}
}

func (e nativeErr) Error() string {
	return fmt.Sprintf("%s failed with code %d", e.op, e.code)
}
