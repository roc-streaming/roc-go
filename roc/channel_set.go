// DO NOT EDIT! Code generated by generate_enums script from roc-toolkit
// roc-toolkit git tag: v0.2.5-11-g14d642e9, commit: 14d642e9

package roc

// Channel set.
//
//go:generate stringer -type ChannelSet -trimprefix ChannelSet -output channel_set_string.go
type ChannelSet int

const (
	// Stereo.
	//
	// Two channels: left and right.
	ChannelSetStereo ChannelSet = 0x3
)
