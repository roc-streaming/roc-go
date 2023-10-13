package roc

// Network protocol.
// Defines URI scheme of Endpoint.
//
//go:generate stringer -type Protocol -trimprefix Proto
type Protocol int

const (
	// RTSP 1.0 (RFC 2326) or RTSP 2.0 (RFC 7826).
	//
	// Interfaces:
	//  - InterfaceConsolidated
	//
	// Transports:
	//   - for signaling: TCP
	//   - for media: RTP and RTCP over UDP or TCP
	//
	ProtoRtsp Protocol = 10

	// RTP over UDP (RFC 3550).
	//
	// Interfaces:
	//  - InterfaceAudioSource
	//
	// Transports:
	//  - UDP
	//
	// Audio encodings:
	//   - PacketEncodingAvpL16
	//
	// FEC encodings:
	//   - none
	//
	ProtoRtp Protocol = 20

	// RTP source packet (RFC 3550) + FECFRAME Reed-Solomon footer (RFC 6865) with m=8.
	//
	// Interfaces:
	//  - InterfaceAudioSource
	//
	// Transports:
	//  - UDP
	//
	// Audio encodings:
	//  - similar to ProtoRtp
	//
	// FEC encodings:
	//  - FecEncodingRs8m
	//
	ProtoRtpRs8mSource Protocol = 30

	// FEC repair packet + FECFRAME Reed-Solomon header (RFC 6865) with m=8.
	//
	// Interfaces:
	//  - InterfaceAudioRepair
	//
	// Transports:
	//  - UDP
	//
	// FEC encodings:
	//  - FecEncodingRs8m
	//
	ProtoRs8mRepair Protocol = 31

	// RTP source packet (RFC 3550) + FECFRAME LDPC-Staircase footer (RFC 6816).
	//
	// Interfaces:
	//  - InterfaceAudioSource
	//
	// Transports:
	//  - UDP
	//
	// Audio encodings:
	//  - similar to ProtoRtp
	//
	// FEC encodings:
	//  - FecEncodingLdpcStaircase
	//
	ProtoRtpLdpcSource Protocol = 32

	// FEC repair packet + FECFRAME LDPC-Staircase header (RFC 6816).
	//
	// Interfaces:
	//  - InterfaceAudioRepair
	//
	// Transports:
	//  - UDP
	//
	// FEC encodings:
	//  - FecEncodingLdpcStaircase
	//
	ProtoLdpcRepair Protocol = 33

	// RTCP over UDP (RFC 3550).
	//
	// Interfaces:
	//  - InterfaceAudioControl
	//
	// Transports:
	//  - UDP
	//
	ProtoRtcp Protocol = 70
)
