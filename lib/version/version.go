// Package version exposes the terrafmt version and git commit, set at build time.
package version

import "runtime/debug"

var Version = "dev"

var GitCommit string

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
			Version = info.Main.Version
		}
	}
}
