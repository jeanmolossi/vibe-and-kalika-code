package security_test

import (
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

func TestIsScriptPath(t *testing.T) {
	if !security.IsScriptPath("setup.sh") || !security.IsScriptPath("setup.ps1") {
		t.Fatal("expected script paths to be detected")
	}
	if security.IsScriptPath("README.md") {
		t.Fatal("did not expect markdown to be detected as script")
	}
}
