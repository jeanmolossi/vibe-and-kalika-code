package backup_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/backup"
)

func TestBackupCreatedAndReportWritten(t *testing.T) {
	projectRoot := t.TempDir()
	source := filepath.Join(projectRoot, "file.txt")
	if err := os.WriteFile(source, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := backup.Create(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := backup.CopyInto(res, projectRoot, source); err != nil {
		t.Fatal(err)
	}
	if err := backup.WriteReport(res); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "file.txt")); err != nil {
		t.Fatalf("backup file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(res.Dir, "backup-report.txt")); err != nil {
		t.Fatalf("backup report missing: %v", err)
	}
}
