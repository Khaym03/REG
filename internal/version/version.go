package version

import (
	"runtime/debug"
	"strings"
)

// version is set at build time via:
// -ldflags "-X github.com/Khaym03/REG/internal/version.version=1.0.0"
var version = "dev"

// Info holds build and VCS metadata for the application.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Dirty   bool   `json:"dirty"`
}

// Get returns the current application version and VCS metadata embedded by the Go toolchain.
func Get() Info {
	return fromBuildInfo(version, readBuildInfo())
}

func readBuildInfo() *debug.BuildInfo {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	return bi
}

func fromBuildInfo(injected string, bi *debug.BuildInfo) Info {
	info := Info{
		Version: strings.TrimPrefix(strings.TrimSpace(injected), "v"),
	}
	if info.Version == "" {
		info.Version = "dev"
	}

	if bi == nil {
		return info
	}

	if info.Version == "dev" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		info.Version = strings.TrimPrefix(bi.Main.Version, "v")
	}

	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.Commit = setting.Value
		case "vcs.time":
			info.Date = setting.Value
		case "vcs.modified":
			info.Dirty = setting.Value == "true"
		}
	}

	return info
}

// NumericFileVersion converts a SemVer string (e.g. "1.0.0-beta.2", "v1.2.3")
// into a Windows-compatible 4-part numeric version ("X.X.X.X").
// Falls back to "0.0.0.0" if the input cannot be parsed.
func NumericFileVersion(semver string) string {
	s := strings.TrimSpace(semver)
	s = strings.TrimPrefix(s, "v")

	if idx := strings.IndexAny(s, "-+"); idx != -1 {
		s = s[:idx]
	}

	if s == "" {
		return "0.0.0.0"
	}

	parts := strings.Split(s, ".")
	if len(parts) > 4 {
		parts = parts[:4]
	}

	for _, p := range parts {
		if p == "" {
			return "0.0.0.0"
		}
		for i := 0; i < len(p); i++ {
			if p[i] < '0' || p[i] > '9' {
				return "0.0.0.0"
			}
		}
	}

	for len(parts) < 4 {
		parts = append(parts, "0")
	}

	return strings.Join(parts, ".")
}
