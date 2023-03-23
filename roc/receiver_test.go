package roc

import (
	"errors"
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

func TestReceiver_SetReuseaddr(t *testing.T) {
	cases := []struct {
		name                 string
		slot                 Slot
		iface                Interface
		receiverClosedBefore bool
		enabled              bool
		wantErr              error
	}{
		{
			name:                 "ok",
			slot:                 SlotDefault,
			iface:                InterfaceAudioSource,
			receiverClosedBefore: false,
			enabled:              true,
			wantErr:              nil,
		},
		{
			name:                 "closed receiver",
			slot:                 SlotDefault,
			iface:                InterfaceAudioSource,
			receiverClosedBefore: true,
			enabled:              true,
			wantErr:              errors.New("receiver is closed"),
		},
		{
			name:                 "bad iface",
			slot:                 SlotDefault,
			iface:                -1,
			receiverClosedBefore: false,
			enabled:              true,
			wantErr:              newNativeErr("roc_receiver_set_reuseaddr()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(ContextConfig{})
			require.NoError(t, err)

			receiver, err := OpenReceiver(ctx, ReceiverConfig{
				FrameSampleRate: 44100,
				FrameChannels:   ChannelSetStereo,
				FrameEncoding:   FrameEncodingPcmFloat,
			})
			require.NoError(t, err)
			require.NotNil(t, receiver)

			if tt.receiverClosedBefore {
				err = receiver.Close()
				require.NoError(t, err)
			}

			err = receiver.SetReuseaddr(tt.slot, tt.iface, tt.enabled)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if !tt.receiverClosedBefore {
				err = receiver.Close()
				require.NoError(t, err)
			}

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}
