package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReceiver_Open(t *testing.T) {
	tests := []struct {
		name    string
		config  ReceiverConfig
		wantErr error
	}{
		{
			name: "ok",
			config: ReceiverConfig{
				FrameSampleRate: 44100,
				FrameChannels:   ChannelSetStereo,
				FrameEncoding:   FrameEncodingPcmFloat,
			},
			wantErr: nil,
		},
		{
			name:    "invalid config",
			config:  ReceiverConfig{},
			wantErr: newNativeErr("roc_receiver_open()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})

			require.NoError(t, err)
			require.NotNil(t, ctx)

			receiver, err := OpenReceiver(ctx, tt.config)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotNil(t, receiver)

				err = receiver.Close()
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
				require.Nil(t, receiver)
			}

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}
