// Code generated by "stringer -type ClockSource -trimprefix ClockSource -output clock_source_string.go"; DO NOT EDIT.

package roc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ClockSourceDefault-0]
	_ = x[ClockSourceExternal-1]
	_ = x[ClockSourceInternal-2]
}

const _ClockSource_name = "DefaultExternalInternal"

var _ClockSource_index = [...]uint8{0, 7, 15, 23}

func (i ClockSource) String() string {
	if i < 0 || i >= ClockSource(len(_ClockSource_index)-1) {
		return "ClockSource(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ClockSource_name[_ClockSource_index[i]:_ClockSource_index[i+1]]
}
