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
