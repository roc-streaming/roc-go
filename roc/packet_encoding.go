package roc

// Packet encoding.
//
//go:generate stringer -type PacketEncoding -trimprefix PacketEncoding
type PacketEncoding int

const (
	// PCM signed 16-bit.
	// "L16" encoding from RTP A/V Profile (RFC 3551).
	// Uncompressed samples coded as interleaved 16-bit signed big-endian
	// integers in two's complement notation.
	PacketEncodingAvpL16 PacketEncoding = 2
)
