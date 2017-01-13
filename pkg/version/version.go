package version

import (
	"fmt"
	"runtime"
)

// Info contains versioning information.
type Info struct {
	GitVersion   string
	GitCommit    string
	GitTreeState string
	GoVersion    string
	Compiler     string
	Platform     string
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return fmt.Sprintf("%s+%s-%s", info.GitVersion, info.GitCommit, info.GitTreeState)
}
