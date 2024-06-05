package roc

import "time"

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
	// latency as close as possible to the target latency.
	// If zero, default value is used. Should not be negative, otherwise an error is returned.
	TargetLatency time.Duration

	// Maximum delta between current and target latency, in nanoseconds.
	// If current latency becomes larger than the target latency plus this value, the
	// session is terminated.
	// If zero, default value is used. Should not be negative, otherwise an error is returned.
	MaxLatencyOverrun time.Duration

	// Maximum delta between target and current latency, in nanoseconds.
	// If current latency becomes smaller than the target latency minus this value, the
	// session is terminated.
	// May be larger than the target latency because current latency may be negative,
	// which means that the playback run ahead of the last packet received from network.
	// If zero, default value is used.  Should not be negative, otherwise an error is returned.
	MaxLatencyUnderrun time.Duration

	// Timeout for the lack of playback, in nanoseconds.
	// If there is no playback during this period, the session is terminated.
	// This mechanism allows to detect dead, hanging, or broken clients
	// generating invalid packets.
	// If zero, default value is used. If negative, the timeout is disabled.
	NoPlaybackTimeout time.Duration

	// Timeout for broken playback, in nanoseconds.
	// If there the playback is considered broken during this period, the session
	// is terminated. The playback is broken if there is a breakage detected at every
	// BreakageDetectionWindow during BrokenPlaybackTimeout.
	// This mechanism allows to detect vicious circles like when all client packets
	// are a bit late and receiver constantly drops them producing unpleasant noise.
	// If zero, default value is used. If negative, the timeout is disabled.
	BrokenPlaybackTimeout time.Duration

	// Breakage detection window, in nanoseconds.
	// If zero, default value is used. Should not be negative, otherwise an error is returned.
	// See BrokenPlaybackTimeout.
	BreakageDetectionWindow time.Duration
}
