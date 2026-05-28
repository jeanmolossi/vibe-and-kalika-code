package report_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/manifest"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/platform"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/report"
)

func TestReportGeneratedWithCorrectContent(t *testing.T) {
	projectRoot := t.TempDir()
	path, err := report.WriteInstallReport(report.InstallReportInput{
		ProjectRoot: projectRoot,
		Manifest:    &manifest.Manifest{Name: "pack", Version: "1.0.0"},
		Source:      "./pack",
		Platforms:   []platform.Platform{platform.PlatformCopilotCLI},
		Operations:  []platform.PlannedOperation{{Type: platform.OperationCreate, TargetPath: filepath.Join(projectRoot, "file.md")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "Package: pack") || !strings.Contains(content, "[create]") {
		t.Fatalf("unexpected report content: %s", content)
	}
}
