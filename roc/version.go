package roc

/*
#include <roc/version.h>
*/
import "C"

import (
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

var (
	versions    *Versions
	versionOnce sync.Once
)

// Version components.
type Versions struct {
	Library SemanticVersion // Native library version, libroc version.
	Binding SemanticVersion // Go binding version.
}

// Semantic version components.
type SemanticVersion struct {
	Major uint64 // Major version component.
	Minor uint64 // Minor version component.
	Patch uint64 // Patch version component.
}

// Retrieve version numbers.
// This function can be used to retrieve actual run-time version of the library.
// It may be different from the compile-time version when using shared library.
func Version() Versions {
	versionOnce.Do(func() {
		versions = &Versions{}

		var cVersion C.struct_roc_version
		C.roc_version_get(&cVersion)
		versions.Library = SemanticVersion{
			Major: uint64(cVersion.major),
			Minor: uint64(cVersion.minor),
			Patch: uint64(cVersion.patch),
		}

		if info, ok := debug.ReadBuildInfo(); ok {
			v := strings.TrimPrefix(info.Main.Version, "v")
			vs := strings.SplitN(v, ".", 3)
			if len(vs) != 3 {
				return
			}
			versions.Binding.Major, _ = strconv.ParseUint(vs[0], 10, 64)
			versions.Binding.Minor, _ = strconv.ParseUint(vs[1], 10, 64)
			versions.Binding.Patch, _ = strconv.ParseUint(vs[2], 10, 64)
		}
	})

	return *versions
}
