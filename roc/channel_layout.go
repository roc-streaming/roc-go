// Code generated by generate_bindings.py script from roc-streaming/bindgen
// roc-toolkit git tag: v0.3.0, commit: 57b932b8

package roc

// Channel layout.
//
// Defines number of channels and meaning of each channel.
//
//go:generate stringer -type ChannelLayout -trimprefix ChannelLayout -output channel_layout_string.go
type ChannelLayout int

const (
	// Multi-track audio.
	//
	// In multitrack layout, stream contains multiple channels which represent
	// independent "tracks" without any special meaning (unlike stereo or surround)
	// and hence without any special processing or mapping.
	//
	// The number of channels is arbitrary and is defined by Tracks field of
	// MediaEncoding struct.
	ChannelLayoutMultitrack ChannelLayout = 1

	// Mono.
	//
	// One channel with monophonic sound.
	ChannelLayoutMono ChannelLayout = 2

	// Stereo.
	//
	// Two channels: left, right.
	ChannelLayoutStereo ChannelLayout = 3
)
