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
