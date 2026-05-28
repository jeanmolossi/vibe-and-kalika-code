package app_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

func writeState(t *testing.T, root string, st *state.Store) {
	t.Helper()
	if err := state.Write(root, st); err != nil {
		t.Fatalf("write state: %v", err)
	}
}

func TestUninstall(t *testing.T) {
	t.Run("created_files_deleted", func(t *testing.T) {
		root := t.TempDir()
		fpath := filepath.Join(root, "somefile.md")
		if err := os.WriteFile(fpath, []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{Package: "pkg-a", CreatedFiles: []string{fpath}},
			},
		})
		_, code, err := app.Uninstall(app.UninstallOptions{Package: "pkg-a", ProjectRoot: root})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != app.ExitSuccess {
			t.Fatalf("expected exit success, got %d", code)
		}
		if _, statErr := os.Stat(fpath); !os.IsNotExist(statErr) {
			t.Fatal("created file should have been deleted")
		}
	})

	t.Run("backup_restored", func(t *testing.T) {
		root := t.TempDir()
		backupDir := t.TempDir()
		originalFile := filepath.Join(backupDir, "original.md")
		if err := os.WriteFile(originalFile, []byte("original content"), 0o644); err != nil {
			t.Fatal(err)
		}
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{Package: "pkg-a", BackupPath: backupDir, Files: []string{"original.md"}},
			},
		})
		res, _, err := app.Uninstall(app.UninstallOptions{Package: "pkg-a", ProjectRoot: root})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.FilesRestored) == 0 {
			t.Fatal("expected files to be reported as restored")
		}
	})

	t.Run("managed_block_removed", func(t *testing.T) {
		root := t.TempDir()
		agentsPath := filepath.Join(root, "AGENTS.md")
		content := "<!-- BEGIN VKC AGENT: agent-a -->\n## Agent: agent-a\n\nhello\n<!-- END VKC AGENT: agent-a -->\n"
		if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{Package: "pkg-a", AgentBlocks: []state.AgentBlock{{Path: agentsPath, AgentName: "agent-a"}}},
			},
		})
		_, _, err := app.Uninstall(app.UninstallOptions{Package: "pkg-a", ProjectRoot: root})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, _ := os.ReadFile(agentsPath)
		if strings.Contains(string(got), "BEGIN VKC AGENT: agent-a") {
			t.Fatalf("managed block should be removed, got:\n%s", got)
		}
	})

	t.Run("agent_name_used_not_package_name", func(t *testing.T) {
		root := t.TempDir()
		agentsPath := filepath.Join(root, "AGENTS.md")
		// agent name differs from package name
		content := "# Header\n\n<!-- BEGIN VKC AGENT: kalika-reviewer -->\n## Agent: kalika-reviewer\n\ncontent\n<!-- END VKC AGENT: kalika-reviewer -->\n"
		if err := os.WriteFile(agentsPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{
					Package:     "basic-kalika-pack",
					AgentBlocks: []state.AgentBlock{{Path: agentsPath, AgentName: "kalika-reviewer"}},
				},
			},
		})
		_, _, err := app.Uninstall(app.UninstallOptions{Package: "basic-kalika-pack", ProjectRoot: root})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, _ := os.ReadFile(agentsPath)
		if strings.Contains(string(got), "BEGIN VKC AGENT: kalika-reviewer") {
			t.Fatalf("agent block should be removed, got:\n%s", got)
		}
		if !strings.Contains(string(got), "# Header") {
			t.Fatalf("surrounding content should be preserved, got:\n%s", got)
		}
	})

	t.Run("state_record_removed", func(t *testing.T) {
		root := t.TempDir()
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{Package: "pkg-a"},
				{Package: "pkg-b"},
			},
		})
		_, _, err := app.Uninstall(app.UninstallOptions{Package: "pkg-a", ProjectRoot: root})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		st, _ := state.Read(root)
		if len(st.Installations) != 1 || st.Installations[0].Package != "pkg-b" {
			t.Fatalf("expected only pkg-b to remain, got: %+v", st.Installations)
		}
	})

	t.Run("package_not_found", func(t *testing.T) {
		root := t.TempDir()
		writeState(t, root, &state.Store{})
		_, _, err := app.Uninstall(app.UninstallOptions{Package: "nonexistent", ProjectRoot: root})
		if err == nil {
			t.Fatal("expected error for missing package")
		}
		if !strings.Contains(err.Error(), "nonexistent") {
			t.Fatalf("error should mention package name, got: %v", err)
		}
	})

	t.Run("already_absent_created_file", func(t *testing.T) {
		root := t.TempDir()
		missingFile := filepath.Join(root, "does-not-exist.md")
		writeState(t, root, &state.Store{
			Installations: []state.Installation{
				{Package: "pkg-a", CreatedFiles: []string{missingFile}},
			},
		})
		_, _, err := app.Uninstall(app.UninstallOptions{Package: "pkg-a", ProjectRoot: root})
		if err != nil {
			t.Fatalf("expected no error for absent file (idempotent), got: %v", err)
		}
	})
}
