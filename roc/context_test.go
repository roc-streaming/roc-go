package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_Open(t *testing.T) {
	tests := []struct {
		name    string
		config  ContextConfig
		wantErr error
	}{
		{
			name:    "ok",
			config:  ContextConfig{MaxPacketSize: 50, MaxFrameSize: 70},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(tt.config)

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotNil(t, ctx)

				err = ctx.Close()
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
				require.Nil(t, ctx)
			}
		})
	}
}

func TestContext_Close(t *testing.T) {
	tests := []struct {
		name         string
		hasReceivers bool
		hasSenders   bool
		wantErr      error
	}{
		{
			name:         "no_senders_or_receivers",
			hasReceivers: false,
			hasSenders:   false,
			wantErr:      nil,
		},
		{
			name:         "has_receivers",
			hasReceivers: true,
			hasSenders:   false,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
		{
			name:         "has_senders",
			hasReceivers: false,
			hasSenders:   true,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
		{
			name:         "has_senders_and_receivers",
			hasReceivers: true,
			hasSenders:   true,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})
			require.NoError(t, err)
			require.NotNil(t, ctx)

			if tt.hasReceivers {
				receiver, err := OpenReceiver(ctx, ReceiverConfig{
					FrameSampleRate:  43100,
					FrameChannels:    ChannelSetStereo,
					FrameEncoding:    FrameEncodingPcmFloat,
				})
				require.NoError(t, err)
				require.NotNil(t, receiver)

				defer func() {
					err = receiver.Close()
					if err != nil {
						t.Fail()
					}
				}()
			}
			if tt.hasSenders {
				sender, err := OpenSender(ctx, SenderConfig{
					FrameSampleRate:  44100,
					FrameChannels:    ChannelSetStereo,
					FrameEncoding:    FrameEncodingPcmFloat,
					ClockSource:      ClockInternal,
					ResamplerProfile: ResamplerProfileDisable,
					FecEncoding:      FecEncodingRs8m,
				})
				require.NoError(t, err)
				require.NotNil(t, sender)
				
				defer func() {
					err = sender.Close()
					if err != nil {
						t.Fail()
					}
				}()
			}

			err = ctx.Close()
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
