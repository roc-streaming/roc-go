package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSender_Open(t *testing.T) {
	tests := []struct {
		name    string
		config  SenderConfig
		wantErr error
	}{
		{
			name: "ok",
			config: SenderConfig{
				FrameSampleRate: 44100,
				FrameChannels:   ChannelSetStereo,
				FrameEncoding:   FrameEncodingPcmFloat,
			},
			wantErr: nil,
		},
		{
			name:    "invalid config",
			config:  SenderConfig{},
			wantErr: newNativeErr("roc_sender_open()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})

			require.NoError(t, err)
			require.NotNil(t, ctx)

			sender, err := OpenSender(ctx, tt.config)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotNil(t, sender)

				err = sender.Close()
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
				require.Nil(t, sender)
			}

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestSenderSetReuseaddr(t *testing.T) {
	ctx, err := OpenContext(ContextConfig{})
	require.NoError(t, err)

	sender, err := OpenSender(ctx, SenderConfig{
		FrameSampleRate: 44100,
		FrameChannels:   ChannelSetStereo,
		FrameEncoding:   FrameEncodingPcmFloat,
	})
	require.NoError(t, err)
	require.NotNil(t, sender)

	err = sender.SetReuseaddr(SlotDefault, InterfaceAudioSource, true)
	require.NoError(t, err)

	err = sender.Close()
	require.NoError(t, err)

	err = ctx.Close()
	require.NoError(t, err)
}
