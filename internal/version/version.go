// Package version holds build-time version information injected via ldflags.
package version

// Version, Commit and Date are populated at build time via ldflags:
//
//	-ldflags "-X github.com/jeanmolossi/vibe-and-kalika-code/internal/version.Version=v1.2.3"
var (
	// Version is the semantic version of the CLI binary.
	Version = "dev"
	// Commit is the git commit SHA at build time.
	Commit = "none"
	// Date is the build timestamp.
	Date = "unknown"
)
