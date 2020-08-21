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

func boolToUint(b bool) C.uint {
	if b {
		return 1
	}
	return 0
}
