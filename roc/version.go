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
	Native   SemanticVersion // Native library version, libroc version.
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
		versions.Native = SemanticVersion{
			Major: uint64(cVersion.major),
			Minor: uint64(cVersion.minor),
			Patch: uint64(cVersion.patch),
		}
		versions.Bindings = parseVersion(bindingsVersion)
	})

	return versions
}

func parseVersion(s string) SemanticVersion {
	vs := strings.SplitN(s, ".", 3)
	if len(vs) != 3 {
		panic("semantic version doesn't have 3 parts")
	}

	v := SemanticVersion{}
	var err error
	v.Major, err = strconv.ParseUint(vs[0], 10, 64)
	if err != nil {
		panic(err)
	}
	v.Minor, err = strconv.ParseUint(vs[1], 10, 64)
	if err != nil {
		panic(err)
	}
	v.Patch, err = strconv.ParseUint(vs[2], 10, 64)
	if err != nil {
		panic(err)
	}

	return v
}
