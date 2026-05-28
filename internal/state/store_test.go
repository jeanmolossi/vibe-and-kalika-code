package state_test

import (
	"testing"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
)

func TestStateWriteAndRead(t *testing.T) {
	projectRoot := t.TempDir()
	store := &state.Store{Installations: []state.Installation{{Package: "pack", Version: "1.0.0"}}}
	if err := state.Write(projectRoot, store); err != nil {
		t.Fatal(err)
	}
	got, err := state.Read(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Installations) != 1 || got.Installations[0].Package != "pack" {
		t.Fatalf("unexpected store: %+v", got)
	}
}
