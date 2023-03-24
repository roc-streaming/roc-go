package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	v := Version()
	require.NotEmpty(t, v)
	require.NotZero(t, v.Library.Major+v.Library.Minor+v.Library.Patch)
	require.NotZero(t, v.Binding.Major+v.Binding.Minor+v.Binding.Patch)
	require.Equal(t, v.Library.Major, v.Binding.Major)
	require.GreaterOrEqual(t, v.Binding.Minor, v.Library.Minor)
}
