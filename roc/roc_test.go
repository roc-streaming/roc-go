package roc

import "testing"

func fail(expected interface{}, got interface{}, t *testing.T) {
	t.Errorf("Mismatch, expected: %v, got: %v", expected, got)
	t.Fail()
}

func Test_roc_address_init(t *testing.T) {
	tests := []struct {
		f    Family
		ip   string
		port int
	}{
		{
			f:    AfIpv4,
			ip:   "192.168.0.1",
			port: 4567,
		},
		{
			f:    AfIpv6,
			ip:   "2001:db8:85a3:1:2:8a2e:37:7334",
			port: 9858,
		},
	}

	for _, tt := range tests {
		a, err := NewAddress(tt.f, tt.ip, tt.port)

		if err != nil {
			fail(nil, err, t)
		}

		if a == nil {
			fail(nil, "Address is nil", t)
		}

		fam, err := a.Family()
		if err != nil {
			fail(nil, err, t)
		}
		if fam != tt.f {
			fail(tt.f, fam, t)
		}

		ip, err := a.IP()
		if err != nil {
			fail(nil, err, t)
		}
		if ip != tt.ip {
			fail(tt.ip, ip, t)
		}

		port, err := a.Port()
		if err != nil {
			fail(nil, err, t)
		}
		if port != tt.port {
			fail(tt.port, port, t)
		}
	}
}

func Test_roc_context_open(t *testing.T) {
	tests := []struct {
		config  ContextConfig
		wantErr error
	}{
		{
			config:  ContextConfig{MaxPacketSize: 50, MaxFrameSize: 70},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx, err := OpenContext(&tt.config)

		if err != tt.wantErr {
			fail(tt.wantErr, err, t)
		}

		err = ctx.Close()
		if err != nil {
			fail(nil, err, t)
		}
	}
}

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
