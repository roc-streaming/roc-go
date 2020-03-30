package roc

import "testing"

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
