package security_test

import (
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/security"
)

func TestValidateRelativePathRejectsTraversal(t *testing.T) {
	if err := security.ValidateRelativePath("../evil"); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateRelativePathAllowsSafePath(t *testing.T) {
	if err := security.ValidateRelativePath("skills/test-skill"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
