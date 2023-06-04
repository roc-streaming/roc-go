package roc

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvert_go2cBool(t *testing.T) {
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

func TestConvert_go2cStr(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    []char
		wantErr error
	}{
		{
			name:    "str",
			arg:     "str",
			want:    []char{'s', 't', 'r', '\x00'},
			wantErr: nil,
		},
		{
			name:    "str0",
			arg:     "str\x00",
			want:    nil,
			wantErr: errors.New("unexpected zero byte in the string: \"str\\x00\""),
		},
		{
			name:    "str0s",
			arg:     "str\x00s",
			want:    nil,
			wantErr: errors.New("unexpected zero byte in the string: \"str\\x00s\""),
		},
		{
			name:    "empty",
			arg:     "",
			want:    []char{'\x00'},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := go2cStr(tt.arg)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestConvert_c2goStr(t *testing.T) {
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

func TestConvert_go2cUnsignedDuration(t *testing.T) {
	tests := []struct {
		name    string
		arg     time.Duration
		want    ulonglong
		wantErr error
	}{
		{
			name:    "gtezero",
			arg:     1,
			want:    (ulonglong)(1),
			wantErr: nil,
		},
		{
			name:    "ltzero",
			arg:     -1,
			want:    0,
			wantErr: fmt.Errorf("unexpected negative duration: %v", (time.Duration)(-1)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := go2cUnsignedDuration(tt.arg)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
