package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	old := [...]string{PackageVersion, PackageCommitID, GitStatus, BuildTime}
	t.Cleanup(func() { PackageVersion, PackageCommitID, GitStatus, BuildTime = old[0], old[1], old[2], old[3] })
	PackageVersion, PackageCommitID, GitStatus, BuildTime = "v7.1.0", "abc", "clean", "2026-11-08T00:00:00Z"
	got := String()
	for _, want := range []string{"v7.1.0", "commit abc", "clean", GoVersion, "built 2026-11-08T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Fatalf("String() = %q, missing %q", got, want)
		}
	}
}
