package roc

import "testing"

func Test_roc_receiver_open(t *testing.T) {
	tests := []struct {
		receiverConfig ReceiverConfig
		wantErr        error
	}{
		{
			receiverConfig: ReceiverConfig{
				FrameSampleRate:         44100,
				FrameChannels:           ChannelSetStereo,
				FrameEncoding:           FrameEncodingPcmFloat,
				AutomaticTiming:         1,
				ResamplerProfile:        ResamplerDisable,
				TargetLatency:           0,
				MaxLatencyOverrun:       0,
				MaxLatencyUnderrun:      0,
				NoPlaybackTimeout:       0,
				BrokenPlaybackTimeout:   0,
				BreakageDetectionWindow: 0},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx, err := OpenContext(&ContextConfig{MaxPacketSize: 0, MaxFrameSize: 0})
		if err != nil {
			fail(nil, err, t)
		}
		receiver, err := OpenReceiver(ctx, &tt.receiverConfig)
		if err != tt.wantErr {
			fail(tt.wantErr, err, t)
		}

		// cleanup
		err = receiver.Close()
		if err != nil {
			fail(nil, err, t)
		}
		err = ctx.Close()
		if err != nil {
			fail(nil, err, t)
		}
	}
}
