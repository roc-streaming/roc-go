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

// Hard-coded version of bindings.
// Should be updated manually each time a new release is tagged.
// This variable is modified only in tests.
var bindingsVersion = "0.2.0"

// Validate version compatibility of Go bindings and native library.
// Must be invoked at all library entry points at least once.
// Entry points refer to exported non-method functions of this package.
// This variable is modified only in tests.
var checkVersionFn = checkVersion

var (
	versionInfo      VersionInfo
	versionInfoOnce  sync.Once
	versionCheckOnce int32
)

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
		versionInfo = fetchVersion()
	})

	return versionInfo
}

// Check compatibility of versions of native library (libroc) and Go bindings.
// If versions are incompatible, then error describing the problem is
// returned, otherwise returns nil.
// When Go bindings are used first time, they automatically run this check and
// panic if it fails. You can run it manually before using bindings.
func (vi VersionInfo) Validate() error {
	errMsg := fmt.Sprintf(
		"detected incompatibility between roc bindings (%s) and native library (%s): ",
		vi.Bindings, vi.Native,
	)

	if vi.Native.Major != vi.Bindings.Major {
		return errors.New(errMsg + "major versions are different")
	}

	if vi.Native.Minor > vi.Bindings.Minor {
		return errors.New(errMsg +
			"minor version of bindings is less than minor version of native library")
	}

	return nil
}

func checkVersion() {
	if atomic.CompareAndSwapInt32(&versionCheckOnce, 0, 1) {
		vi := fetchVersion()

		if err := vi.Validate(); err != nil {
			panic(err.Error())
		}
	}
}

func fetchVersion() VersionInfo {
	var cVersion C.struct_roc_version
	C.roc_version_get(&cVersion)

	return VersionInfo{
		Bindings: parseVersion(bindingsVersion),
		Native: SemanticVersion{
			Major: uint64(cVersion.major),
			Minor: uint64(cVersion.minor),
			Patch: uint64(cVersion.patch),
		},
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
