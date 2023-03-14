package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	v := Version()
	require.NotEmpty(t, v)
	require.NotZero(t, v.Library.Major+v.Library.Minor+v.Library.Patch)
}
