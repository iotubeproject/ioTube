// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

// Package version holds the build version, set at link time:
//
//	go build -ldflags "-X github.com/iotexproject/ioTube/witness-service/version.PackageVersion=v7.1.0 ..."
//
// The Makefile and the Dockerfiles set it (see VERSION / COMMIT / GIT_STATUS /
// BUILD_TIME there).
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

const noBuildInfo = "NoBuildInfo"

var (
	// PackageVersion is `git describe --tags --dirty` of the build, e.g. v7.1.0
	// or v7.1.0-3-gabc1234.
	PackageVersion = noBuildInfo
	// PackageCommitID is the full git commit of the build.
	PackageCommitID = noBuildInfo
	// GitStatus is "clean" or "dirty".
	GitStatus = noBuildInfo
	// BuildTime is the UTC build time.
	BuildTime = noBuildInfo
	// GoVersion is the Go toolchain the binary was built with.
	GoVersion = runtime.Version()
)

func init() {
	// A plain `go build` in a git checkout records the commit; use it when the
	// linker flags were not set.
	if PackageCommitID != noBuildInfo {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			PackageCommitID = s.Value
		case "vcs.modified":
			if GitStatus == noBuildInfo {
				GitStatus = map[string]string{"true": "dirty", "false": "clean"}[s.Value]
			}
		}
	}
}

// String describes the build in one line.
func String() string {
	return fmt.Sprintf("%s (commit %s, %s, %s, built %s)", PackageVersion, PackageCommitID, GitStatus, GoVersion, BuildTime)
}
