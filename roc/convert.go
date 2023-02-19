package roc

import (
	"C"
)

func go2cStr(str string) []C.char {
	charArray := make([]C.char, len(str)+1)
	for ind, r := range str {
		charArray[ind] = (C.char)(r)
	}
	charArray[len(str)] = '\x00'
	return charArray
}

func c2goStr(charArray []C.char) string {
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
