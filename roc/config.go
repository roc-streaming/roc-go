package roc

/*
#include "roc/config.h"
*/
import "C"

// Network slot.
//
// A peer (sender or receiver) may have multiple slots, which may be independently
// bound or connected. You can use multiple slots on sender to connect it to multiple
// receiver addresses, and you can use multiple slots on receiver to bind it to
// multiple receiver address.
//
// Slots are numbered from zero and are created implicitly. Just specify slot index
// when binding or connecting endpoint, and slot will be automatically created if it
// was not created yet.
//
// In simple cases, just use SlotDefault.
//
// Each slot has its own set of interfaces, dedicated to different kinds of endpoints.
// See Interface for details.
type Slot int

const (
	// Alias for the slot with index zero.
	SlotDefault Slot = 0
)

// Network interface.
//
// Interface is a way to access the peer (sender or receiver) via network.
//
// Each peer slot has multiple interfaces, one of each type. The user interconnects
// peers by binding one of the first peer's interfaces to an URI and then connecting the
// corresponding second peer's interface to that URI.
//
// A URI is represented by Endpoint object.
//
// The interface defines the type of the communication with the remote peer and the
// set of protocols (URI schemes) that can be used with this particular interface.
//
// InterfaceConsolidated is an interface for high-level protocols which
// automatically manage all necessary communication: transport streams, control messages,
// parameter negotiation, etc. When a consolidated connection is established, peers may
// automatically setup lower-level interfaces like InterfaceAudioSource,
// InterfaceAudioRepair, and InterfaceAudioControl.
//
// InterfaceConsolidated is mutually exclusive with lower-level interfaces.
// In most cases, the user needs only InterfaceConsolidated. However, the
// lower-level interfaces may be useful if an external signaling mechanism is used or for
// compatibility with third-party software.
//
// InterfaceAudioSource and InterfaceAudioRepair are lower-level
// unidirectional transport-only interfaces. The first is used to transmit audio stream,
// and the second is used to transmit redundant repair stream, if FEC is enabled.
//
// InterfaceAudioControl is a lower-level interface for control streams.
// If you use InterfaceAudioSource and InterfaceAudioRepair, you
// usually also need to use InterfaceAudioControl to enable carrying additional
// non-transport information.
//
//go:generate stringer -type Interface -trimprefix Interface
type Interface int

const (
	// Interface that consolidates all types of streams (source, repair, control).
	//
	// Allowed operations:
	//  - bind    (sender, receiver)
	//  - connect (sender, receiver)
	//
	// Allowed protocols:
	//  - ProtoRtsp
	//
	InterfaceConsolidated Interface = 1

	// Interface for audio stream source data.
	//
	// Allowed operations:
	//  - bind    (receiver)
	//  - connect (sender)
	//
	// Allowed protocols:
	//  -  ProtoRtp
	//  -  ProtoRtpRs8mSource
	//  -  ProtoRtpLdpcSource
	//
	InterfaceAudioSource Interface = 11

	// Interface for audio stream repair data.
	//
	// Allowed operations:
	//  - bind    (receiver)
	//  - connect (sender)
	//
	// Allowed protocols:
	//  -  ProtoRs8mRepair
	//  -  ProtoLdpcRepair
	//
	InterfaceAudioRepair Interface = 12

	// Interface for audio control messages.
	//
	// Allowed operations:
	//  - bind    (sender, receiver)
	//  - connect (sender, receiver)
	//
	// Allowed protocols:
	//  -  ProtoRrcp
	//
	InterfaceAudioControl Interface = 13
)

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

// Forward Error Correction encoding.
//
//go:generate stringer -type FecEncoding -trimprefix FecEncoding
type FecEncoding int

const (
	// No FEC encoding.
	// Compatible with ProtoRtp protocol.
	FecEncodingDisable FecEncoding = -1

	// Default FEC encoding.
	// Current default is FecEncodingRs8m.
	FecEncodingDefault FecEncoding = 0

	// Reed-Solomon FEC encoding (RFC 6865) with m=8.
	// Good for small block sizes (below 256 packets).
	// Compatible with ProtoRtpRs8mSource and ProtoRs8mRepair
	// protocols for source and repair endpoints.
	FecEncodingRs8m FecEncoding = 1

	// LDPC-Staircase FEC encoding (RFC 6816).
	// Good for large block sizes (above 1024 packets).
	// Compatible with ProtoRtpLdpcSource and ProtoLdpcRepair
	// protocols for source and repair endpoints.
	FecEncodingLdpcStaircase FecEncoding = 2
)

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

// Channel set.
//
//go:generate stringer -type ChannelSet -trimprefix ChannelSet
type ChannelSet int

const (
	// Stereo.
	// Two channels: left and right.
	ChannelSetStereo ChannelSet = 0x3
)

// Resampler backend.
// Affects speed and quality.
// Some backends may be disabled at build time.
//
//go:generate stringer -type ResamplerBackend -trimprefix ResamplerBackend
type ResamplerBackend int

const (
	// Default backend.
	// Depends on what was enabled at build time.
	ResamplerBackendDefault ResamplerBackend = 0

	// Slow built-in resampler.
	// Always available.
	ResamplerBackendBuiltin ResamplerBackend = 1

	// Fast good-quality resampler from SpeexDSP.
	// May be disabled at build time.
	ResamplerBackendSpeex ResamplerBackend = 2
)

// Resampler profile.
// Affects speed and quality.
// Each resampler backend treats profile in its own way.
//
//go:generate stringer -type ResamplerProfile -trimprefix ResamplerProfile
type ResamplerProfile int

const (
	// Do not perform resampling.
	// Clock drift compensation will be disabled in this case.
	// If in doubt, do not disable resampling.
	ResamplerProfileDisable ResamplerProfile = -1

	// Default profile.
	// Current default is ResamplerProfileMedium.
	ResamplerProfileDefault ResamplerProfile = 0

	// High quality, low speed.
	ResamplerProfileHigh ResamplerProfile = 1

	// Medium quality, medium speed.
	ResamplerProfileMedium ResamplerProfile = 2

	// Low quality, high speed.
	ResamplerProfileLow ResamplerProfile = 3
)

// Clock source for sender or receiver.
//
//go:generate stringer -type ClockSource -trimprefix Clock
type ClockSource int

const (
	// Sender or receiver is clocked by external user-defined clock.
	// Write and read operations are non-blocking. The user is responsible
	// to call them in time, according to the external clock.
	ClockExternal ClockSource = 0

	// Sender or receiver is clocked by an internal clock.
	// Write and read operations are blocking. They automatically wait until it's time
	// to process the next bunch of samples according to the configured sample rate.
	ClockInternal ClockSource = 1
)

// Context configuration.
// You can zero-initialize this struct to get a default config.
// See also Context.
type ContextConfig struct {
	// Maximum size in bytes of a network packet.
	// Defines the amount of bytes allocated per network packet.
	// Sender and receiver won't handle packets larger than this.
	// If zero, default value is used.
	MaxPacketSize uint32

	// Maximum size in bytes of an audio frame.
	// Defines the amount of bytes allocated per intermediate internal frame in the
	// pipeline. Does not limit the size of the frames provided by user.
	// If zero, default value is used.
	MaxFrameSize uint32
}

// Sender configuration.
// You can zero-initialize this struct to get a default config.
// See also Sender.
type SenderConfig struct {
	// The rate of the samples in the frames passed to sender.
	// Number of samples per channel per second.
	// If FrameSampleRate and PacketSampleRate are different,
	// resampler should be enabled.
	// Should be set.
	FrameSampleRate uint32

	// The channel set in the frames passed to sender.
	// Should be set.
	FrameChannels ChannelSet

	// The sample encoding in the frames passed to sender.
	// Should be set.
	FrameEncoding FrameEncoding

	// The rate of the samples in the packets generated by sender.
	// Number of samples per channel per second.
	// If zero, default value is used.
	PacketSampleRate uint32

	// The channel set in the packets generated by sender.
	// If zero, default value is used.
	PacketChannels ChannelSet

	// The sample encoding in the packets generated by sender.
	// If zero, default value is used.
	PacketEncoding PacketEncoding

	// The length of the packets produced by sender, in nanoseconds.
	// Number of nanoseconds encoded per packet.
	// The samples written to the sender are buffered until the full packet is
	// accumulated or the sender is flushed or closed. Larger number reduces
	// packet overhead but also increases latency.
	// If zero, default value is used.
	PacketLength uint64

	// Enable packet interleaving.
	// If true, the sender shuffles packets before sending them. This
	// may increase robustness but also increases latency.
	PacketInterleaving bool

	// Clock source to use.
	// Defines whether write operation will be blocking or non-blocking.
	// If zero, default value is used.
	ClockSource ClockSource

	// Resampler backend to use.
	ResamplerBackend ResamplerBackend

	// Resampler profile to use.
	// If non-zero, the sender employs resampler if the frame sample rate differs
	// from the packet sample rate.
	ResamplerProfile ResamplerProfile

	// FEC encoding to use.
	// If non-zero, the sender employs a FEC codec to generate redundant packets
	// which may be used on receiver to restore lost packets. This requires both
	// sender and receiver to use two separate source and repair ports.
	FecEncoding FecEncoding

	// Number of source packets per FEC block.
	// Used if some FEC code is selected.
	// Larger number increases robustness but also increases latency.
	// If zero, default value is used.
	FecBlockSourcePackets uint32

	// Number of repair packets per FEC block.
	// Used if some FEC code is selected.
	// Larger number increases robustness but also increases traffic.
	// If zero, default value is used.
	FecBlockRepairPackets uint32
}

// Receiver configuration.
// You can zero-initialize this struct to get a default config.
// See also Receiver.
type ReceiverConfig struct {
	// The rate of the samples in the frames returned to the user.
	// Number of samples per channel per second.
	// Should be set.
	FrameSampleRate uint32

	// The channel set in the frames returned to the user.
	// Should be set.
	FrameChannels ChannelSet

	// The sample encoding in the frames returned to the user.
	// Should be set.
	FrameEncoding FrameEncoding

	// Clock source to use.
	// Defines whether read operation will be blocking or non-blocking.
	// If zero, default value is used.
	ClockSource ClockSource

	// Resampler backend to use.
	ResamplerBackend ResamplerBackend

	// Resampler profile to use.
	// If non-zero, the receiver employs resampler for two purposes:
	//  - adjust the sender clock to the receiver clock, which may differ a bit
	//  - convert the packet sample rate to the frame sample rate if they are different
	ResamplerProfile ResamplerProfile

	// Target latency, in nanoseconds.
	// The session will not start playing until it accumulates the requested latency.
	// Then, if resampler is enabled, the session will adjust its clock to keep actual
	// latency as close as close as possible to the target latency.
	// If zero, default value is used.
	TargetLatency uint64

	// Maximum delta between current and target latency, in nanoseconds.
	// If current latency becomes larger than the target latency plus this value, the
	// session is terminated.
	// If zero, default value is used.
	MaxLatencyOverrun uint64

	// Maximum delta between target and current latency, in nanoseconds.
	// If current latency becomes smaller than the target latency minus this value, the
	// session is terminated.
	// May be larger than the target latency because current latency may be negative,
	// which means that the playback run ahead of the last packet received from network.
	// If zero, default value is used.
	MaxLatencyUnderrun uint64

	// Timeout for the lack of playback, in nanoseconds.
	// If there is no playback during this period, the session is terminated.
	// This mechanism allows to detect dead, hanging, or broken clients
	// generating invalid packets.
	// If zero, default value is used. If negative, the timeout is disabled.
	NoPlaybackTimeout int64

	// Timeout for broken playback, in nanoseconds.
	// If there the playback is considered broken during this period, the session
	// is terminated. The playback is broken if there is a breakage detected at every
	// BreakageDetectionWindow during BrokenPlaybackTimeout.
	// This mechanism allows to detect vicious circles like when all client packets
	// are a bit late and receiver constantly drops them producing unpleasant noise.
	// If zero, default value is used. If negative, the timeout is disabled.
	BrokenPlaybackTimeout int64

	// Breakage detection window, in nanoseconds.
	// If zero, default value is used.
	// See BrokenPlaybackTimeout.
	BreakageDetectionWindow uint64
}
