// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"runtime/debug"
	"sync"
)

func Version() string {
	versionOnce.Do(func() {
		version, _ = versionFn()
		if version == "" {
			version = "(unknown version)"
		}
	})
	return version
}

var (
	version     string
	versionOnce sync.Once
)

const pkgPath = "github.com/antoniszymanski/invidious-go"

func versionFn() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, module := range info.Deps {
		if pkgPath == module.Path {
			return module.Version, true
		}
	}
	if pkgPath == info.Main.Path {
		return info.Main.Version, true
	}
	return "", false
}
