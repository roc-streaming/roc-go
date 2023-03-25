package roc

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	v := Version()
	require.NotZero(t, v.Library.Major+v.Library.Minor+v.Library.Patch)
	require.NotZero(t, v.Bindings.Major+v.Bindings.Minor+v.Bindings.Patch)
	require.Equal(t, v.Library.Major, v.Bindings.Major)
	require.GreaterOrEqual(t, v.Bindings.Minor, v.Library.Minor)
}

func Test_bindingsVersion(t *testing.T) {
	bvs := strings.SplitN(bindingsVersion, ".", 3)
	require.Len(t, bvs, 3)
	var err error
	_, err = strconv.ParseUint(bvs[0], 10, 64)
	require.NoError(t, err)
	_, err = strconv.ParseUint(bvs[1], 10, 64)
	require.NoError(t, err)
	_, err = strconv.ParseUint(bvs[2], 10, 64)
	require.NoError(t, err)
}
