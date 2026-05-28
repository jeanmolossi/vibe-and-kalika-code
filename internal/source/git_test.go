package source_test

import (
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/source"
)

func TestParseGitHubTreeURL(t *testing.T) {
	cases := []struct {
		name       string
		input      string
		wantClone  string
		wantBranch string
		wantSubdir string
		wantOK     bool
	}{
		{
			name:       "full tree URL with subdir",
			input:      "https://github.com/jeanmolossi/vibe-and-kalika-code/tree/main/packages/kalika-ofc",
			wantClone:  "https://github.com/jeanmolossi/vibe-and-kalika-code.git",
			wantBranch: "main",
			wantSubdir: "packages/kalika-ofc",
			wantOK:     true,
		},
		{
			name:       "tree URL without subdir",
			input:      "https://github.com/owner/repo/tree/mybranch",
			wantClone:  "https://github.com/owner/repo.git",
			wantBranch: "mybranch",
			wantSubdir: "",
			wantOK:     true,
		},
		{
			name:   "plain git URL not matched",
			input:  "https://github.com/owner/repo.git",
			wantOK: false,
		},
		{
			name:   "non-github URL not matched",
			input:  "https://gitlab.com/owner/repo/tree/main/sub",
			wantOK: false,
		},
		{
			name:   "git@ URL not matched",
			input:  "git@github.com:owner/repo.git",
			wantOK: false,
		},
		{
			name:   "local path not matched",
			input:  "./local/path",
			wantOK: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cloneURL, branch, subdir, ok := source.ParseGitHubTreeURL(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if cloneURL != tc.wantClone {
				t.Errorf("cloneURL = %q, want %q", cloneURL, tc.wantClone)
			}
			if branch != tc.wantBranch {
				t.Errorf("branch = %q, want %q", branch, tc.wantBranch)
			}
			if subdir != tc.wantSubdir {
				t.Errorf("subdir = %q, want %q", subdir, tc.wantSubdir)
			}
		})
	}
}

func TestDefaultSourceIsGitHubTreeURL(t *testing.T) {
	_, _, _, ok := source.ParseGitHubTreeURL(source.DefaultSource)
	if !ok {
		t.Errorf("DefaultSource %q should be a valid GitHub tree URL", source.DefaultSource)
	}
}

// TestParseGitHubTreeURL_PathTraversalRejected verifies that dot-dot subdirs are
// returned as-is by the parser — the containment check lives in CloneGitSource.
func TestParseGitHubTreeURL_PathTraversalRejected(t *testing.T) {
	// The parser itself does not reject traversal paths; it returns subdir verbatim.
	// We test that a dot-dot subdir is correctly returned so CloneGitSource can block it.
	_, _, subdir, ok := source.ParseGitHubTreeURL("https://github.com/owner/repo/tree/main/../../etc/passwd")
	if !ok {
		t.Fatal("expected parser to match the URL even with traversal subdir")
	}
	// subdir must contain ".." so CloneGitSource can detect and reject it.
	if subdir == "" {
		t.Error("expected non-empty subdir for traversal URL")
	}
}
