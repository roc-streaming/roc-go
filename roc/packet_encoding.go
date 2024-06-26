// Code generated by generate_bindings.py script from roc-streaming/bindgen
// roc-toolkit git tag: v0.4.0, commit: 62401be9

package roc

// Packet encoding.
//
// Each packet encoding defines sample format, channel layout, and rate. Each
// packet encoding is compatible with specific protocols.
//
//go:generate stringer -type PacketEncoding -trimprefix PacketEncoding -output packet_encoding_string.go
type PacketEncoding int

const (
	// PCM signed 16-bit, 1 channel, 44100 rate.
	//
	// Represents 1-channel L16 stereo encoding from RTP A/V Profile (RFC 3551).
	// Uses uncompressed samples coded as interleaved 16-bit signed big-endian
	// integers in two's complement notation.
	//
	// Supported by protocols:
	//  - ProtoRtp
	//  - ProtoRtpRs8mSource
	//  - ProtoRtpLdpcSource
	PacketEncodingAvpL16Mono PacketEncoding = 11

	// PCM signed 16-bit, 2 channels, 44100 rate.
	//
	// Represents 2-channel L16 stereo encoding from RTP A/V Profile (RFC 3551).
	// Uses uncompressed samples coded as interleaved 16-bit signed big-endian
	// integers in two's complement notation.
	//
	// Supported by protocols:
	//  - ProtoRtp
	//  - ProtoRtpRs8mSource
	//  - ProtoRtpLdpcSource
	PacketEncodingAvpL16Stereo PacketEncoding = 10
)
