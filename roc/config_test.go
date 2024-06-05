package roc

func makeContextConfig() ContextConfig {
	return ContextConfig{
		MaxPacketSize: 2000,
		MaxFrameSize:  4000,
	}
}

func makeSenderConfig() SenderConfig {
	return SenderConfig{
		FrameEncoding: makeMediaEncoding(),
		ClockSource:   ClockSourceInternal,
		FecEncoding:   FecEncodingRs8m,
	}
}

func makeReceiverConfig() ReceiverConfig {
	return ReceiverConfig{
		FrameEncoding:    makeMediaEncoding(),
		ClockSource:      ClockSourceInternal,
		ClockSyncBackend: ClockSyncBackendDisable,
	}
}

func makeInterfaceConfig() InterfaceConfig {
	return InterfaceConfig{
		OutgoingAddress: "127.0.0.1",
		MulticastGroup:  "",
		ReuseAddress:    true,
	}
}

func makeMediaEncoding() MediaEncoding {
	return MediaEncoding{
		Rate:     44100,
		Format:   FormatPcmFloat32,
		Channels: ChannelLayoutStereo,
		Tracks:   0,
	}
}
