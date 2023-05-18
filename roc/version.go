package roc

/*
#include <roc/version.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const bindingsVersion = "0.2.0"

var (
	versionInfo      VersionInfo
	versionInfoOnce  sync.Once
	versionCheckOnce int32
)

var versionCheckFn = versionCheck

// Semantic version components.
type SemanticVersion struct {
	Major uint64
	Minor uint64
	Patch uint64
}

// Return string format of semantic version.
func (sv SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
}

// Version components.
type VersionInfo struct {
	Native   SemanticVersion // Version of native library (libroc).
	Bindings SemanticVersion // Version of Go bindings.
}

// Retrieve version numbers.
// This function can be used to retrieve actual run-time version of the library.
// It may be different from the compile-time version when using shared library.
func Version() VersionInfo {
	versionInfoOnce.Do(func() {
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

// Check compatibility of versions of native library (libroc) and Go bindings.
// If versions are incompatible, then error describing the problem is
// returned, otherwise returns nil.
// When Go bindings are used first time, they automatically run this check and
// panic if it fails. You can run it manually before using bindings.
func (vi VersionInfo) Validate() error {
	bindingsVersion := vi.Bindings
	nativeVersion := vi.Native

	errMsg := fmt.Sprintf(
		"detected incompatibility between roc bindings (%s) and native library (%s): ",
		bindingsVersion, nativeVersion,
	)

	if nativeVersion.Major != bindingsVersion.Major {
		return errors.New(errMsg + "major versions are different")
	}

	if nativeVersion.Minor > bindingsVersion.Minor {
		return errors.New(errMsg +
			"minor version of bindings is less than minor version of native library")
	}

	return nil
}

// Validate version compatibility of roc bindings and native library.
// This function must be called at all library entry point at least once.
// Entry points refer to exported non-method functions of this package.
func versionCheck() {
	var versionInfo VersionInfo
	logWrite(LogDebug, "entering versionCheck()")
	defer logWrite(LogDebug, "leaving versionCheck(): version=%+v", versionInfo)

	if atomic.CompareAndSwapInt32(&versionCheckOnce, 0, 1) {
		versionInfo = Version()
		err := versionInfo.Validate()
		if err != nil {
			panic(err.Error())
		}
	}
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
