package roc

func makeContextConfig() ContextConfig {
	return ContextConfig{
		MaxPacketSize: 2000,
		MaxFrameSize:  4000,
	}
}

func makeSenderConfig() SenderConfig {
	return SenderConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		ClockSource:      ClockInternal,
		ResamplerProfile: ResamplerProfileDisable,
		FecEncoding:      FecEncodingRs8m,
	}
}

func makeReceiverConfig() ReceiverConfig {
	return ReceiverConfig{
		FrameSampleRate:  44100,
		FrameChannels:    ChannelSetStereo,
		FrameEncoding:    FrameEncodingPcmFloat,
		ClockSource:      ClockInternal,
		ResamplerProfile: ResamplerProfileDisable,
	}
}
