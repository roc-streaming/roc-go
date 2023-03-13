package roc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion_Get(t *testing.T) {
	v := GetVersion()
	require.NotEmpty(t, v)
	require.NotEmpty(t, v.Library)
}
