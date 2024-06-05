package roc

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
