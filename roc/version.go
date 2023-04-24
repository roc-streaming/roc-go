package roc

/*
#include <roc/version.h>
*/
import "C"

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	versionInfo      VersionInfo
	versionOnce      sync.Once
	versionCheckOnce int32
)

var versionCheckFn = versionCheck

const bindingsVersion = "0.2.0"

// Semantic version components.
type SemanticVersion struct {
	Major uint64 // Major version component.
	Minor uint64 // Minor version component.
	Patch uint64 // Patch version component.
}

// Return string format of semantic version
func (sv SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

// Version components.
type VersionInfo struct {
	Native   SemanticVersion // Native library version, libroc version.
	Bindings SemanticVersion // Go bindings version.
}

// Retrieve version numbers.
// This function can be used to retrieve actual run-time version of the library.
// It may be different from the compile-time version when using shared library.
func Version() VersionInfo {
	versionOnce.Do(func() {
		var cVersion C.struct_roc_version
		C.roc_version_get(&cVersion)
		versionInfo.Native = SemanticVersion{
			Major: uint64(cVersion.major),
			Minor: uint64(cVersion.minor),
			Patch: uint64(cVersion.patch),
		}
		versionInfo.Bindings = parseVersion(bindingsVersion)
	})

	return versionInfo
}

// Runs check for version compatibility
// If versions are incompatible then error describing problem is returned
// Else returns nil
func (versionInfo VersionInfo) Validate() error {
	bindingsVersion := versionInfo.Bindings
	nativeVersion := versionInfo.Native

	errMsg := fmt.Sprintf(
		"Detected incompatibility between roc bindings ( %s ) and native library ( %s ): ",
		bindingsVersion, nativeVersion,
	)

	if nativeVersion.Major != bindingsVersion.Major {
		return fmt.Errorf(errMsg + "Major versions are different")
	}

	if nativeVersion.Minor > bindingsVersion.Minor {
		return fmt.Errorf(errMsg + "Minor version of binding is less than native library")
	}

	return nil
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

// Validate version compatibility of roc bindings and native library.
// This function must be called at all library entry point at least once.
// Entry points refer to exported non-method functions of this package
func versionCheck() {
	if atomic.CompareAndSwapInt32(&versionCheckOnce, 0, 1) {
		versionInfo := Version()
		err := versionInfo.Validate()
		if err != nil {
			panic(err.Error())
		}
	}
}
