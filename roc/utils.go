package roc

import "C"

func toCStr(str string) []C.char {
	charArray := make([]C.char, len(str)+1)
	for ind, r := range str {
		charArray[ind] = (C.char)(r)
	}
	charArray[len(str)] = '\x00'
	return charArray
}
