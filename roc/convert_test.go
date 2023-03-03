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

func Test_go2cStr(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		result []char
	}{
		{name: "str", arg: "str", result: []char{'s', 't', 'r', '\x00'}},
		{name: "str0", arg: "str\x00", result: []char{'s', 't', 'r', '\x00', '\x00'}},
		{name: "str0s", arg: "str\x00s", result: []char{'s', 't', 'r', '\x00', 's', '\x00'}},
		{name: "empty", arg: "", result: []char{'\x00'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, go2cStr(tt.arg))
		})
	}
}

func Test_c2goStr(t *testing.T) {
	tests := []struct {
		name   string
		arg    []char
		result string
	}{
		{name: "str0", arg: []char{'s', 't', 'r', '\x00'}, result: "str"},
		{name: "str00", arg: []char{'s', 't', 'r', '\x00', '\x00'}, result: "str"},
		{name: "str0str0", arg: []char{'s', 't', 'r', '\x00', 's', 't', 'r', '\x00'}, result: "str"},
		{name: "0", arg: []char{'\x00'}, result: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, c2goStr(tt.arg))
		})
	}
}
