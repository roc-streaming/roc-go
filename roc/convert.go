package roc

import (
	"C"
	"fmt"
	"time"
)

type (
	char      = C.char
	longlong  = C.longlong
	ulonglong = C.ulonglong
)

func go2cStr(str string) ([]char, error) {
	charArray := make([]char, len(str)+1)
	for ind, r := range str {
		c := (char)(r)
		if c == '\x00' {
			return nil, fmt.Errorf("unexpected zero byte in the string: %q", str)
		}
		charArray[ind] = c
	}
	charArray[len(str)] = '\x00'
	return charArray, nil
}

func c2goStr(charArray []char) string {
	byteArray := make([]byte, 0, len(charArray)-1)
	for _, c := range charArray {
		if c == '\x00' {
			break
		}
		byteArray = append(byteArray, byte(c))
	}
	return string(byteArray)
}

func go2cBool(b bool) C.uint {
	if b {
		return 1
	}
	return 0
}

func go2cSignedDuration(d time.Duration) longlong {
	return (longlong)(d)
}

func go2cUnsignedDuration(d time.Duration) (ulonglong, error) {
	if d < 0 {
		return 0, fmt.Errorf("unexpected negative duration: %v", d)
	}
	return (ulonglong)(d), nil
}
