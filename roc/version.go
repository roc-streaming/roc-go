package roc

/*
#include <roc/version.h>
*/
import "C"

import (
	"strconv"
	"strings"
	"sync"
)

var (
	versions    Versions
	versionOnce sync.Once
)

const bindingsVersion = "0.2.0"

// Version components.
type Versions struct {
	Library  SemanticVersion // Native library version, libroc version.
	Bindings SemanticVersion // Go bindings version.
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
		var cVersion C.struct_roc_version
		C.roc_version_get(&cVersion)
		versions.Library = SemanticVersion{
			Major: uint64(cVersion.major),
			Minor: uint64(cVersion.minor),
			Patch: uint64(cVersion.patch),
		}

		bvs := strings.SplitN(bindingsVersion, ".", 3)
		if len(bvs) != 3 {
			return
		}
		versions.Bindings.Major, _ = strconv.ParseUint(bvs[0], 10, 64)
		versions.Bindings.Minor, _ = strconv.ParseUint(bvs[1], 10, 64)
		versions.Bindings.Patch, _ = strconv.ParseUint(bvs[2], 10, 64)
	})

	return versions
}
