package roc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_Info(t *testing.T) {
	v := Version()

	assert.NotZero(t, v.Native.Major+v.Native.Minor+v.Native.Patch)
	assert.NotZero(t, v.Bindings.Major+v.Bindings.Minor+v.Bindings.Patch)
	assert.Equal(t, v.Native.Major, v.Bindings.Major)
	assert.GreaterOrEqual(t, v.Bindings.Minor, v.Native.Minor)
}

func TestVersion_Check(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		wantPanic bool
	}{
		{
			name:      "ok",
			version:   bindingsVersion,
			wantPanic: false,
		},
		{
			name:      "validation error",
			version:   "999.999.999",
			wantPanic: true,
		},
		{
			name:      "parsing error",
			version:   "xxx",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalBindingsVersion := bindingsVersion
			defer func() { bindingsVersion = originalBindingsVersion }()

			bindingsVersion = tt.version
			versionCheckOnce = 0

			if tt.wantPanic {
				require.Panics(t, func() { checkVersionFn() })
			} else {
				require.NotPanics(t, func() { checkVersionFn() })
			}
		})
	}
}

func TestVersion_Entrypoints(t *testing.T) {
	tests := []struct {
		name       string
		entrypoint func()
	}{
		{
			name: "OpenContext",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				require.NoError(t, err)
				defer ctx.Close()
			},
		},
		{
			name: "OpenSender",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				require.NoError(t, err)
				defer ctx.Close()

				sender, err := OpenSender(ctx, makeSenderConfig())
				require.NoError(t, err)
				defer sender.Close()
			},
		},
		{
			name: "OpenReceiver",
			entrypoint: func() {
				ctx, err := OpenContext(ContextConfig{})
				require.NoError(t, err)
				defer ctx.Close()

				recv, err := OpenReceiver(ctx, makeReceiverConfig())
				require.NoError(t, err)
				defer recv.Close()
			},
		},
		{
			name: "ParseEndpoint",
			entrypoint: func() {
				_, err := ParseEndpoint("rtp://192.168.0.1:1234")
				require.NoError(t, err)
			},
		},
		{
			name: "SetLogLevel",
			entrypoint: func() {
				SetLogLevel(defaultLogLevel)
			},
		},
		{
			name: "SetLogger",
			entrypoint: func() {
				SetLogger(nil)
			},
		},
		{
			name: "SetLoggerFunc",
			entrypoint: func() {
				SetLoggerFunc(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCheckVersionFn := checkVersionFn
			defer func() { checkVersionFn = originalCheckVersionFn }()

			versionCheckCalled := false
			checkVersionFn = func() { versionCheckCalled = true }

			tt.entrypoint()
			require.True(t, versionCheckCalled)
		})
	}
}

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name     string
		versions VersionInfo
		wantErr  error
	}{
		{
			name: "compatible: equal versions",
			versions: VersionInfo{
				Bindings: SemanticVersion{1, 1, 1},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: nil,
		},
		{
			name: "incompatible: binding major version less than native",
			versions: VersionInfo{
				Bindings: SemanticVersion{0, 1, 1},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: errors.New(
				"detected incompatibility between roc bindings (0.1.1) and native library (1.1.1):" +
					" major versions are different",
			),
		},
		{
			name: "incompatible: binding major version greater than native",
			versions: VersionInfo{
				Bindings: SemanticVersion{2, 1, 1},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: errors.New(
				"detected incompatibility between roc bindings (2.1.1) and native library (1.1.1):" +
					" major versions are different",
			),
		},
		{
			name: "compatible: binding minor version greater than native",
			versions: VersionInfo{
				Bindings: SemanticVersion{1, 2, 1},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: nil,
		},
		{
			name: "incompatible: binding minor version less than native",
			versions: VersionInfo{
				Bindings: SemanticVersion{1, 0, 1},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: errors.New(
				"detected incompatibility between roc bindings (1.0.1) and native library (1.1.1):" +
					" minor version of bindings is less than minor version of native library",
			),
		},
		{
			name: "compatible: binding patch less than native",
			versions: VersionInfo{
				Bindings: SemanticVersion{1, 1, 0},
				Native:   SemanticVersion{1, 1, 1},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.versions.Validate()
			require.Equal(t, err, tt.wantErr)
		})
	}
}

func TestVersion_Parse(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		want      SemanticVersion
		wantPanic bool
	}{
		{
			name: "valid",
			arg:  "3.2.1",
			want: SemanticVersion{3, 2, 1},
		},
		{
			name:      "empty",
			arg:       "",
			wantPanic: true,
		},
		{
			name:      "garbage",
			arg:       "xxx",
			wantPanic: true,
		},
		{
			name:      "less parts",
			arg:       "3.2",
			wantPanic: true,
		},
		{
			name:      "more parts",
			arg:       "3.2.1.0",
			wantPanic: true,
		},
		{
			name:      "major invalid",
			arg:       "a.2.1",
			wantPanic: true,
		},
		{
			name:      "minor invalid",
			arg:       "3.a.1",
			wantPanic: true,
		},
		{
			name:      "patch invalid",
			arg:       "3.2.a",
			wantPanic: true,
		},
		{
			name: "hard-coded bindings version",
			arg:  bindingsVersion,
			want: parseVersion(bindingsVersion),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { parseVersion(tt.arg) })
			} else {
				require.Equal(t, tt.want, parseVersion(tt.arg))
			}
		})
	}
}
