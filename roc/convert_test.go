package roc

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_go2cBool(t *testing.T) {
	tests := []struct {
		arg    bool
		result uint
	}{
		{true, 1},
		{false, 0},
	}
	for _, tt := range tests {
		t.Run(strconv.FormatBool(tt.arg), func(t *testing.T) {
			assert.Equal(t, tt.result, uint(go2cBool(tt.arg)))
		})
	}
}

func Test_go2cStr_c2goStr(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		result string
	}{
		{name: "nil"},
		{name: "str", str: "str", result: "str"},
		{name: "empty", str: "", result: ""},
		{name: "x00", str: "\x00", result: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, c2goStr(go2cStr(tt.str)))
		})
	}
}
