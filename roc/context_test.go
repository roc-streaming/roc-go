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
			name:    "default config",
			config:  ContextConfig{},
			wantErr: nil,
		},
		{
			name:    "custom config",
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
			name:         "no senders or receivers",
			hasReceivers: false,
			hasSenders:   false,
			wantErr:      nil,
		},
		{
			name:         "has receivers",
			hasReceivers: true,
			hasSenders:   false,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
		{
			name:         "has senders",
			hasReceivers: false,
			hasSenders:   true,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
		{
			name:         "has senders and receivers",
			hasReceivers: true,
			hasSenders:   true,
			wantErr:      newNativeErr("roc_context_close()", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				receiver *Receiver
				sender   *Sender
			)

			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)
			require.NotNil(t, ctx)

			if tt.hasReceivers {
				receiver, err = OpenReceiver(ctx, makeReceiverConfig())
				require.NoError(t, err)
				require.NotNil(t, receiver)
			}
			if tt.hasSenders {
				sender, err = OpenSender(ctx, makeSenderConfig())
				require.NoError(t, err)
				require.NotNil(t, sender)
			}

			err = ctx.Close()
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.hasReceivers || tt.hasSenders {
				if tt.hasReceivers {
					err = receiver.Close()
					require.NoError(t, err)
				}
				if tt.hasSenders {
					err = sender.Close()
					require.NoError(t, err)
				}

				err = ctx.Close()
				require.NoError(t, err)
			}
		})
	}
}
