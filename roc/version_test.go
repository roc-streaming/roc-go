package roc

import (
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
