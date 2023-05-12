package roc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion_Get(t *testing.T) {
	v := Version()
	require.NotZero(t, v.Native.Major+v.Native.Minor+v.Native.Patch)
	require.NotZero(t, v.Bindings.Major+v.Bindings.Minor+v.Bindings.Patch)
	require.Equal(t, v.Native.Major, v.Bindings.Major)
	require.GreaterOrEqual(t, v.Bindings.Minor, v.Native.Minor)
}

func TestVersion_Parse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		want      SemanticVersion
		wantPanic bool
	}{
		{
			name: "valid",
			args: args{
				s: "3.2.1",
			},
			want: SemanticVersion{
				Major: 3,
				Minor: 2,
				Patch: 1,
			},
		},
		{
			name: "doesn't have 3 parts",
			args: args{
				s: "3.2",
			},
			wantPanic: true,
		},
		{
			name: "major invalid",
			args: args{
				s: "a.2.1",
			},
			wantPanic: true,
		},
		{
			name: "minor invalid",
			args: args{
				s: "3.a.1",
			},
			wantPanic: true,
		},
		{
			name: "patch invalid",
			args: args{
				s: "3.2.a",
			},
			wantPanic: true,
		},
		{
			name: "hard-coded bindings version",
			args: args{
				s: bindingsVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { parseVersion(tt.args.s) })
			} else if tt.want != (SemanticVersion{}) {
				require.Equal(t, tt.want, parseVersion(tt.args.s))
			}
		})
	}
}

func TestVersion_Validite(t *testing.T) {
	tests := []struct {
		name        string
		versionInfo VersionInfo
		wantErr     error
	}{
		{
			name: "compatible: equal versions",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{1, 1, 1},
			},
			wantErr: nil,
		},
		{
			name: "incompatible: binding's major version less than native",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{0, 1, 1},
			},
			wantErr: errors.New(
				"Detected incompatibility between roc bindings ( 0.1.1 ) and native library ( 1.1.1 ):" +
					" Major versions are different",
			),
		},
		{
			name: "incompatible: binding's major version greater than native",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{2, 1, 1},
			},
			wantErr: errors.New(
				"Detected incompatibility between roc bindings ( 0.1.1 ) and native library ( 1.1.1 ):" +
					" Major versions are different",
			),
		},
		{
			name: "compatible: binding's minor version greater than native",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{1, 2, 1},
			},
			wantErr: nil,
		},
		{
			name: "incompatible: binding's minor version less than native",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{1, 0, 1},
			},
			wantErr: errors.New(
				"Detected incompatibility between roc bindings ( 1.0.1 ) and native library ( 1.1.1 ):" +
					" Minor version of bindings is less than native library",
			),
		},
		{
			name: "compatible: binding's patch less than native",
			versionInfo: VersionInfo{
				Native:   SemanticVersion{1, 1, 1},
				Bindings: SemanticVersion{1, 1, 0},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.versionInfo.Validate()
			if tt.wantErr != nil {
				require.Error(t, err, tt.wantErr)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestVersion_Check(t *testing.T) {
	require.NotPanics(t, func() { versionCheck() })
}

func TestVersion_Entrypoints(t *testing.T) {
	var versionCheckCalled = false

	versionCheckFn = func() {
		versionCheckCalled = true
	}

	defer func() {
		versionCheckFn = versionCheck
	}()

	tests := []struct {
		name       string
		entrypoint func()
	}{
		{
			name: "OpenContext",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				if err != nil {
					panic(err)
				}
				defer ctx.Close()
			},
		},
		{
			name: "OpenSender",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				if err != nil {
					panic(err)
				}
				sender, err := OpenSender(ctx, makeSenderConfig())
				if err != nil {
					panic(err)
				}
				defer ctx.Close()
				defer sender.Close()
			},
		},
		{
			name: "OpenReceiver",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				if err != nil {
					panic(err)
				}
				recv, err := OpenReceiver(ctx, makeReceiverConfig())
				if err != nil {
					panic(err)
				}
				defer ctx.Close()
				defer recv.Close()
			},
		},
		{
			name: "ParseEndpoint",
			entrypoint: func() {
				_, err := ParseEndpoint("rtp://192.168.0.1:1234")
				if err != nil {
					panic(err)
				}
			},
		},
		{
			name:       "SetLogLevel",
			entrypoint: func() { SetLogLevel(LogError) },
		},
		{
			name:       "SetLogger",
			entrypoint: func() { SetLogger(nil) },
		},
		{
			name:       "SetLoggerFunc",
			entrypoint: func() { SetLoggerFunc(nil) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.entrypoint()
			require.Equal(t, versionCheckCalled, true)
			versionCheckCalled = false
		})
	}
}
