package roc

// Channel set.
//
//go:generate stringer -type ChannelSet -trimprefix ChannelSet
type ChannelSet int

const (
	// Stereo.
	// Two channels: left and right.
	ChannelSetStereo ChannelSet = 0x3
)
