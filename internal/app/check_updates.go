package app

import (
	"strings"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/version"
)

// UpdateCheck holds the result of an update availability check.
type UpdateCheck struct {
	CLIUpdateAvailable bool
	CLICurrentVersion  string
	CLILatestVersion   string
	PackageUpdates     []PackageUpdateInfo
}

// PackageUpdateInfo describes a package that has an update available.
type PackageUpdateInfo struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
}

// CheckUpdates checks whether the CLI or any installed packages have updates available.
// It never returns an error — failures are silently ignored so the CLI is not blocked.
func CheckUpdates(projectRoot string) *UpdateCheck {
	result := &UpdateCheck{}

	if version.Version != "dev" {
		rel, err := fetchLatestRelease(defaultHTTPClient)
		if err == nil {
			latest := strings.TrimPrefix(rel.TagName, "v")
			current := strings.TrimPrefix(version.Version, "v")
			if isOutdated(current, latest) {
				result.CLIUpdateAvailable = true
				result.CLICurrentVersion = version.Version
				result.CLILatestVersion = rel.TagName
			}
		}
	}

	st, err := state.Read()
	if err != nil {
		return result
	}

	for _, inst := range st.Installations {
		latest, err := fetchLatestVersion(inst.Source, projectRoot)
		if err != nil {
			continue
		}
		if isOutdated(inst.Version, latest) {
			result.PackageUpdates = append(result.PackageUpdates, PackageUpdateInfo{
				Name:           inst.Package,
				CurrentVersion: inst.Version,
				LatestVersion:  latest,
			})
		}
	}

	return result
}
