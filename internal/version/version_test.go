package version

import (
	"runtime/debug"
	"testing"
)

func TestFromBuildInfo(t *testing.T) {
	bi := &debug.BuildInfo{
		Main: debug.Module{Version: "v1.2.3"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abcdef1234567890"},
			{Key: "vcs.time", Value: "2026-11-23T19:00:00Z"},
			{Key: "vcs.modified", Value: "true"},
		},
	}

	t.Run("injected version takes precedence over build info main version", func(t *testing.T) {
		got := fromBuildInfo("v1.0.0-beta.2", bi)
		if got.Version != "1.0.0-beta.2" {
			t.Fatalf("expected Version 1.0.0-beta.2, got %q", got.Version)
		}
		if got.Commit != "abcdef1234567890" {
			t.Fatalf("expected Commit abcdef1234567890, got %q", got.Commit)
		}
		if got.Date != "2026-11-23T19:00:00Z" {
			t.Fatalf("expected Date 2026-11-23T19:00:00Z, got %q", got.Date)
		}
		if !got.Dirty {
			t.Fatalf("expected Dirty true, got false")
		}
	})

	t.Run("falls back to build info main version when dev", func(t *testing.T) {
		got := fromBuildInfo("dev", bi)
		if got.Version != "1.2.3" {
			t.Fatalf("expected Version 1.2.3, got %q", got.Version)
		}
	})

	t.Run("keeps dev when build info main version is devel", func(t *testing.T) {
		develBI := &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
		}
		got := fromBuildInfo("dev", develBI)
		if got.Version != "dev" {
			t.Fatalf("expected Version dev, got %q", got.Version)
		}
	})

	t.Run("handles nil build info", func(t *testing.T) {
		got := fromBuildInfo("", nil)
		if got.Version != "dev" {
			t.Fatalf("expected Version dev, got %q", got.Version)
		}
	})
}

func TestNumericFileVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "1.0.0-beta.2", want: "1.0.0.0"},
		{in: "v1.2.3", want: "1.2.3.0"},
		{in: "1.2.3.4", want: "1.2.3.4"},
		{in: "v2.5.1+build.42", want: "2.5.1.0"},
		{in: "1.0", want: "1.0.0.0"},
		{in: "dev", want: "0.0.0.0"},
		{in: "", want: "0.0.0.0"},
		{in: "1.a.0", want: "0.0.0.0"},
	}

	for _, tc := range tests {
		if got := NumericFileVersion(tc.in); got != tc.want {
			t.Errorf("NumericFileVersion(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
