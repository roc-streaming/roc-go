package roc

// Frame encoding.
//
//go:generate stringer -type FrameEncoding -trimprefix FrameEncoding
type FrameEncoding int

const (
	// PCM floats.
	// Uncompressed samples coded as floats in range [-1; 1].
	// Channels are interleaved, e.g. two channels are encoded as "L R L R ...".
	FrameEncodingPcmFloat FrameEncoding = 1
)
