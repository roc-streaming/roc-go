// Code generated by generate_bindings.py script from roc-streaming/bindgen
// roc-toolkit git tag: v0.3.0, commit: 57b932b8

package roc

import "time"

// Sender configuration.
// You can zero-initialize this struct to get a default config.
// See also Sender.
type SenderConfig struct {
	// The encoding used in frames passed to sender.
	//
	// Frame encoding defines sample format, channel layout, and sample rate in
	// local frames created by user and passed to sender. Should be set (zero value
	// is invalid).
	FrameEncoding MediaEncoding

	// The encoding used for packets produced by sender.
	//
	// Packet encoding defines sample format, channel layout, and sample rate in
	// network packets. If packet encoding differs from frame encoding, conversion
	// is performed automatically.
	//
	// If zero, sender selects packet encoding automatically based on
	// FrameEncoding. This automatic selection matches only encodings that have
	// exact same sample rate and channel layout, and hence don't require
	// conversions. If you need conversions, you should set packet encoding
	// explicitly.
	//
	// If you want to force specific packet encoding, and built-in set of encodings
	// is not enough, you can use ContextRegisterEncoding() to register custom
	// encoding, set PacketEncoding to registered identifier. If you use signaling
	// protocol like RTSP, it's enough to register in just on sender; otherwise,
	// you need to do the same on receiver as well.
	PacketEncoding PacketEncoding

	// The length of the packets produced by sender, in nanoseconds.
	//
	// Number of nanoseconds encoded per packet. The samples written to the sender
	// are buffered until the full packet is accumulated or the sender is flushed
	// or closed. Larger number reduces packet overhead but also increases latency.
	// If zero, default value is used.
	PacketLength time.Duration

	// Enable packet interleaving.
	//
	// If true, the sender shuffles packets before sending them. This may
	// increase robustness but also increases latency.
	PacketInterleaving bool

	// FEC encoding to use.
	//
	// If non-zero, the sender employs a FEC encoding to generate redundant packets
	// which may be used on receiver to restore lost packets. This requires both
	// sender and receiver to use two separate source and repair endpoints.
	FecEncoding FecEncoding

	// Number of source packets per FEC block.
	//
	// Used if some FEC encoding is selected.
	//
	// Sender divides stream into blocks of N source (media) packets, and adds M
	// repair (redundancy) packets to each block, where N is FecBlockSourcePackets
	// and M is FecBlockRepairPackets.
	//
	// Larger number of source packets in block increases robustness (repair
	// ratio), but also increases latency.
	//
	// If zero, default value is used.
	FecBlockSourcePackets uint32

	// Number of repair packets per FEC block.
	//
	// Used if some FEC encoding is selected.
	//
	// Sender divides stream into blocks of N source (media) packets, and adds M
	// repair (redundancy) packets to each block, where N is FecBlockSourcePackets
	// and M is FecBlockRepairPackets.
	//
	// Larger number of repair packets in block increases robustness (repair
	// ratio), but also increases traffic. Number of repair packets usually should
	// be 1/2 or 2/3 of the number of source packets.
	//
	// If zero, default value is used.
	FecBlockRepairPackets uint32

	// Clock source to use.
	//
	// Defines whether write operation will be blocking or non-blocking. If zero,
	// default value is used ( ClockSourceExternal ).
	ClockSource ClockSource

	// Resampler backend to use.
	//
	// If zero, default value is used.
	ResamplerBackend ResamplerBackend

	// Resampler profile to use.
	//
	// If non-zero, the sender employs resampler if the frame sample rate differs
	// from the packet sample rate.
	ResamplerProfile ResamplerProfile
}