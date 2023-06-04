package roc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReceiver_Open(t *testing.T) {
	tests := []struct {
		name        string
		contextFunc func() *Context
		configFunc  func() ReceiverConfig
		wantErr     error
	}{
		{
			name: "ok",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: makeReceiverConfig,
			wantErr:    nil,
		},
		{
			name: "nil context",
			contextFunc: func() *Context {
				return nil
			},
			configFunc: makeReceiverConfig,
			wantErr:    errors.New("context is nil"),
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
			configFunc: makeReceiverConfig,
			wantErr:    errors.New("context is closed"),
		},
		{
			name: "invalid config.FrameSampleRate",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.FrameSampleRate = 0
				return rc
			},
			wantErr: newNativeErr("roc_receiver_open()", -1),
		},
		{
			name: "invalid config.FrameChannels",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.FrameChannels = 0
				return rc
			},
			wantErr: newNativeErr("roc_receiver_open()", -1),
		},
		{
			name: "invalid config.FrameEncoding",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.FrameEncoding = 0
				return rc
			},
			wantErr: newNativeErr("roc_receiver_open()", -1),
		},
		{
			name: "invalid config.TargetLatency",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.TargetLatency = -1
				return rc
			},
			wantErr: fmt.Errorf("invalid config.TargetLatency: %w",
				fmt.Errorf("unexpected negative duration: -1ns")),
		},
		{
			name: "invalid config.MaxLatencyOverrun",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.MaxLatencyOverrun = -1
				return rc
			},
			wantErr: fmt.Errorf("invalid config.MaxLatencyOverrun: %w",
				fmt.Errorf("unexpected negative duration: -1ns")),
		},
		{
			name: "invalid config.MaxLatencyUnderrun",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.MaxLatencyUnderrun = -1
				return rc
			},
			wantErr: fmt.Errorf("invalid config.MaxLatencyUnderrun: %w",
				fmt.Errorf("unexpected negative duration: -1ns")),
		},
		{
			name: "invalid config.BreakageDetectionWindow",
			contextFunc: func() *Context {
				ctx, err := OpenContext(makeContextConfig())
				require.NoError(t, err)
				return ctx
			},
			configFunc: func() ReceiverConfig {
				rc := makeReceiverConfig()
				rc.BreakageDetectionWindow = -1
				return rc
			},
			wantErr: fmt.Errorf("invalid config.BreakageDetectionWindow: %w",
				fmt.Errorf("unexpected negative duration: -1ns")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextFunc()

			receiver, err := OpenReceiver(ctx, tt.configFunc())

			if tt.wantErr == nil {
				require.NoError(t, err)
				require.NotNil(t, receiver)

				err = receiver.Close()
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.wantErr, err)
				require.Nil(t, receiver)
			}

			if ctx != nil {
				err = ctx.Close()
				require.NoError(t, err)
			}
		})
	}
}

func TestReceiver_SetMulticastGroup(t *testing.T) {
	cases := []struct {
		name    string
		slot    Slot
		iface   Interface
		ip      string
		wantErr error
	}{
		{
			name:  "ok",
			slot:  SlotDefault,
			iface: InterfaceAudioSource,
			ip:    "127.0.0.1",
		},
		{
			name:    "bad iface",
			slot:    SlotDefault,
			iface:   -1,
			ip:      "127.0.0.1",
			wantErr: newNativeErr("roc_receiver_set_multicast_group()", -1),
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
			wantErr: newNativeErr("roc_receiver_set_multicast_group()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			receiver, err := OpenReceiver(ctx, makeReceiverConfig())
			require.NoError(t, err)
			require.NotNil(t, receiver)

			err = receiver.SetMulticastGroup(tt.slot, tt.iface, tt.ip)
			require.Equal(t, tt.wantErr, err)

			err = receiver.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}

}

func TestReceiver_SetReuseaddr(t *testing.T) {
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
			wantErr: newNativeErr("roc_receiver_set_reuseaddr()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			receiver, err := OpenReceiver(ctx, makeReceiverConfig())
			require.NoError(t, err)
			require.NotNil(t, receiver)

			err = receiver.SetReuseaddr(tt.slot, tt.iface, tt.enabled)
			require.Equal(t, tt.wantErr, err)

			err = receiver.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestReceiver_Bind(t *testing.T) {
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
			wantErr:  newNativeErr("roc_receiver_bind()", -1),
		},
		{
			name:     "bad protocol",
			slot:     SlotDefault,
			iface:    InterfaceAudioSource,
			endpoint: &Endpoint{Host: "127.0.0.1", Port: 0, Protocol: 1},
			wantErr:  newNativeErr("roc_endpoint_set_protocol()", -1),
		},
		{
			name:     "bad iface",
			slot:     SlotDefault,
			iface:    -1,
			endpoint: baseEndpoint,
			wantErr:  newNativeErr("roc_receiver_bind()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			receiver, err := OpenReceiver(ctx, makeReceiverConfig())
			require.NoError(t, err)
			require.NotNil(t, receiver)

			err = receiver.Bind(tt.slot, tt.iface, tt.endpoint)
			require.Equal(t, tt.wantErr, err)

			err = receiver.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestReceiver_ReadFloats(t *testing.T) {
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
			wantErr: newNativeErr("roc_receiver_read()", -1),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := OpenContext(makeContextConfig())
			require.NoError(t, err)

			receiver, err := OpenReceiver(ctx, makeReceiverConfig())
			require.NoError(t, err)
			require.NotNil(t, receiver)

			err = receiver.ReadFloats(tt.frame)
			require.Equal(t, tt.wantErr, err)

			err = receiver.Close()
			require.NoError(t, err)

			err = ctx.Close()
			require.NoError(t, err)
		})
	}
}

func TestReceiver_Close(t *testing.T) {
	cases := []struct {
		name      string
		operation func(receiver *Receiver) error
	}{
		{
			name: "SetReuseaddr after close",
			operation: func(receiver *Receiver) error {
				return receiver.SetReuseaddr(SlotDefault, InterfaceAudioSource, true)
			},
		},
		{
			name: "SetMulticastGroup after close",
			operation: func(receiver *Receiver) error {
				return receiver.SetMulticastGroup(SlotDefault, InterfaceAudioSource, "127.0.0.1")
			},
		},
		{
			name: "Bind after close",
			operation: func(receiver *Receiver) error {
				return receiver.Bind(SlotDefault, InterfaceAudioSource, nil)
			},
		},
		{
			name: "ReadFloats after close",
			operation: func(receiver *Receiver) error {
				recFloats := make([]float32, 2)
				return receiver.ReadFloats(recFloats)
			},
		},
	}
	for _, tt := range cases {
		ctx, err := OpenContext(makeContextConfig())
		require.NoError(t, err)
		require.NotNil(t, ctx)

		receiver, err := OpenReceiver(ctx, makeReceiverConfig())
		require.NoError(t, err)
		require.NotNil(t, receiver)

		err = receiver.Close()
		require.NoError(t, err)

		require.Equal(t, errors.New("receiver is closed"), tt.operation(receiver))

		err = ctx.Close()
		require.NoError(t, err)
	}
}
