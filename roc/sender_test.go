package roc

import (
	"errors"
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

func TestSender_SetReuseaddr(t *testing.T) {
	cases := []struct {
		name               string
		slot               Slot
		iface              Interface
		senderClosedBefore bool
		enabled            bool
		wantErr            error
	}{
		{
			name:               "ok",
			slot:               SlotDefault,
			iface:              InterfaceAudioSource,
			senderClosedBefore: false,
			enabled:            true,
			wantErr:            nil,
		},
		{
			name:               "closed sender",
			slot:               SlotDefault,
			iface:              InterfaceAudioSource,
			senderClosedBefore: true,
			enabled:            true,
			wantErr:            errors.New("sender is closed"),
		},
		{
			name:               "bad iface",
			slot:               SlotDefault,
			iface:              -1,
			senderClosedBefore: false,
			enabled:            true,
			wantErr:            newNativeErr("roc_sender_set_reuseaddr()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})
			require.NoError(t, err)

			sender, err := OpenSender(ctx, SenderConfig{
				FrameSampleRate: 44100,
				FrameChannels:   ChannelSetStereo,
				FrameEncoding:   FrameEncodingPcmFloat,
			})
			require.NoError(t, err)
			require.NotNil(t, sender)

			if tt.senderClosedBefore {
				err = sender.Close()
				require.NoError(t, err)
			}

			err = sender.SetReuseaddr(tt.slot, tt.iface, tt.enabled)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if !tt.senderClosedBefore {
				err = sender.Close()
				require.NoError(t, err)
			}

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}
