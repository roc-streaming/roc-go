package roc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSender_Open(t *testing.T) {
	tests := []struct {
		name        string
		contextFunc func() *Context
		config      SenderConfig
		wantErr     error
	}{
		{
			name: "ok",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			config:  makeSenderConfig(),
			wantErr: nil,
		},

		{
			name: "invalid config.PacketLength",
			contextFunc: func() *Context {
				ctx, err := OpenContext(ContextConfig{})
				require.NoError(t, err)
				return ctx
			},
			config: invalidconfigPacketLength(),
			wantErr: fmt.Errorf("invalid config.PacketLength: %w",
				fmt.Errorf("unexpected negative duration: -1ns")),
		},
		{
			name: "invalid config",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			config:  SenderConfig{},
			wantErr: newNativeErr("roc_sender_open()", -1),
		},
		{
			name: "nil context",
			contextFunc: func() *Context {
				return nil
			},
			config:  makeSenderConfig(),
			wantErr: errors.New("context is nil"),
		},
		{
			name: "closed context",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)

				err = ctx.Close()
				require.NoError(t, err)
				return ctx
			},
			config:  makeSenderConfig(),
			wantErr: errors.New("context is closed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextFunc()

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

			if ctx != nil {
				err = ctx.Close()
				require.NoError(t, err)
			}
		})
	}
}

func TestSender_SetOutgoingAddress(t *testing.T) {
	cases := []struct {
		name    string
		slot    Slot
		iface   Interface
		ip      string
		wantErr error
	}{
		{
			name:    "ok",
			slot:    SlotDefault,
			iface:   InterfaceAudioSource,
			ip:      "127.0.0.1",
			wantErr: nil,
		},
		{
			name:    "bad iface",
			slot:    SlotDefault,
			iface:   -1,
			ip:      "127.0.0.1",
			wantErr: newNativeErr("roc_sender_set_outgoing_address()", -1),
		},
		{
			name:  "invalid ip",
			slot:  SlotDefault,
			iface: InterfaceAudioSource,
			ip:    "127.0.0.1\x00",
			wantErr: fmt.Errorf("invalid ip: %w",
				fmt.Errorf("unexpected zero byte in the string: \"127.0.0.1\\x00\"")),
		},
		{
			name:    "out of range ip",
			slot:    SlotDefault,
			iface:   InterfaceAudioSource,
			ip:      "256.256.256.256",
			wantErr: newNativeErr("roc_sender_set_outgoing_address()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			sender, err := OpenSender(ctx, makeSenderConfig())
			require.NoError(t, err)
			require.NotNil(t, sender)

			err = sender.SetOutgoingAddress(tt.slot, tt.iface, tt.ip)
			require.Equal(t, tt.wantErr, err)

			err = sender.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestSender_SetReuseaddr(t *testing.T) {
	cases := []struct {
		name    string
		slot    Slot
		iface   Interface
		enabled bool
		wantErr error
	}{
		{
			name:    "ok",
			slot:    SlotDefault,
			iface:   InterfaceAudioSource,
			enabled: true,
			wantErr: nil,
		},
		{
			name:    "bad iface",
			slot:    SlotDefault,
			iface:   -1,
			enabled: true,
			wantErr: newNativeErr("roc_sender_set_reuseaddr()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			sender, err := OpenSender(ctx, makeSenderConfig())
			require.NoError(t, err)
			require.NotNil(t, sender)

			err = sender.SetReuseaddr(tt.slot, tt.iface, tt.enabled)
			require.Equal(t, tt.wantErr, err)

			err = sender.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestSender_Connect(t *testing.T) {
	baseEndpoint, err := ParseEndpoint("rtp+rs8m://127.0.0.1:0")
	require.NoError(t, err)
	require.NotNil(t, baseEndpoint)

	cases := []struct {
		name     string
		slot     Slot
		iface    Interface
		endpoint *Endpoint
		wantErr  error
	}{
		{
			name:     "ok",
			slot:     SlotDefault,
			iface:    InterfaceAudioSource,
			endpoint: baseEndpoint,
			wantErr:  nil,
		},
		{
			name:    "nil endpoint",
			slot:    SlotDefault,
			iface:   InterfaceAudioSource,
			wantErr: errors.New("endpoint is nil"),
		},
		{
			name:     "bad endpoint",
			slot:     SlotDefault,
			iface:    InterfaceAudioSource,
			endpoint: &Endpoint{Host: "127.0.0.1", Port: 0, Protocol: ProtoRs8mRepair},
			wantErr:  newNativeErr("roc_sender_connect()", -1),
		},
		{
			name:     "bad protocol",
			slot:     SlotDefault,
			iface:    InterfaceAudioSource,
			endpoint: &Endpoint{Host: "127.0.0.1", Port: 0, Protocol: 1},
			wantErr:  newNativeErr("roc_endpoint_set_protocol()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			sender, err := OpenSender(ctx, makeSenderConfig())
			require.NoError(t, err)
			require.NotNil(t, sender)

			err = sender.Connect(tt.slot, tt.iface, tt.endpoint)
			require.Equal(t, tt.wantErr, err)

			err = sender.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestSender_WriteFloats(t *testing.T) {
	baseFrameCnt := 2
	baseFrame := make([]float32, baseFrameCnt)
	for i := 0; i < baseFrameCnt; i++ {
		baseFrame[i] = float32(i + 1)
	}

	cases := []struct {
		name    string
		frame   []float32
		wantErr error
	}{
		{
			name:    "ok",
			frame:   baseFrame,
			wantErr: nil,
		},
		{
			name:    "nil frame",
			wantErr: errors.New("frame is nil"),
		},
		{
			name:    "empty frame",
			frame:   []float32{},
			wantErr: nil,
		},
		{
			name:    "bad frame",
			frame:   []float32{1.0},
			wantErr: newNativeErr("roc_sender_write()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			sender, err := OpenSender(ctx, makeSenderConfig())
			require.NoError(t, err)
			require.NotNil(t, sender)

			err = sender.WriteFloats(tt.frame)
			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			err = sender.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestSender_Close(t *testing.T) {
	cases := []struct {
		name      string
		operation func(sender *Sender) error
	}{
		{
			name: "SetReuseaddr after close",
			operation: func(sender *Sender) error {
				return sender.SetReuseaddr(SlotDefault, InterfaceAudioSource, true)
			},
		},
		{
			name: "SetOutgoingAddress after close",
			operation: func(sender *Sender) error {
				return sender.SetOutgoingAddress(SlotDefault, InterfaceAudioSource, "127.0.0.1")
			},
		},
		{
			name: "Connect after close",
			operation: func(sender *Sender) error {
				return sender.Connect(SlotDefault, InterfaceAudioSource, nil)
			},
		},
		{
			name: "WriteFloats after close",
			operation: func(sender *Sender) error {
				recFloats := make([]float32, 2)
				return sender.WriteFloats(recFloats)
			},
		},
	}
	for _, tt := range cases {
		ctx, err := OpenContext(makeContextConfig())
		require.NoError(t, err)
		require.NotNil(t, ctx)

		sender, err := OpenSender(ctx, makeSenderConfig())
		require.NoError(t, err)
		require.NotNil(t, sender)

		err = sender.Close()
		require.NoError(t, err)

		require.Equal(t, errors.New("sender is closed"), tt.operation(sender))

		err = ctx.Close()
		require.NoError(t, err)
	}
}
