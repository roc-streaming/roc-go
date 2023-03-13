package roc

/*
#include <roc/version.h>
*/
import "C"

import (
	"runtime/debug"
	"sync"
)

var (
	version     *Version
	versionOnce sync.Once
)

// Version components.
type Version struct {
	Library *SemanticVersion // Native library version, librock version.
	Binding string           // Binding version.
}

// Semantic version components.
type SemanticVersion struct {
	Major uint // Major version component.
	Minor uint // Minor version component.
	Patch uint // Patch version component.
}

// Retrieve version numbers.
// This function can be used to retrieve actual run-time version of the library.
// It may be different from the compile-time version when using shared library.
func GetVersion() *Version {
	versionOnce.Do(func() {
		version = &Version{}

		var cVersion C.struct_roc_version
		C.roc_version_get(&cVersion)
		version.Library = &SemanticVersion{
			Major: uint(cVersion.major),
			Minor: uint(cVersion.minor),
			Patch: uint(cVersion.patch),
		}

		if info, ok := debug.ReadBuildInfo(); ok {
			version.Binding = info.Main.Version
		}
	})

	return version
}
