package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Mordeak/cp77-modorder-gui/internal/config"
	"github.com/Mordeak/cp77-modorder-gui/internal/conflict"
)

func TestWriteModlistClearsNewStateImmediately(t *testing.T) {
	tempDir := t.TempDir()
	a := &App{
		cfg: &config.Config{},
		result: &conflict.Result{Mods: []*conflict.ModInfo{
			{Name: "alpha.archive"},
			{Name: "charlie.archive"},
			{Name: "bravo.archive"},
		}},
		modDir:              tempDir,
		modlistOrder:        []string{"alpha.archive", "charlie.archive"},
		modlistSet:          map[string]bool{"alpha.archive": true, "charlie.archive": true},
		initialModlistOrder: []string{"alpha.archive", "charlie.archive"},
	}

	before := a.buildScanResult()
	if !before.Rows[2].Unlisted || before.Rows[2].Name != "bravo.archive" {
		t.Fatalf("row before Apply = %+v, want bravo.archive marked new", before.Rows[2])
	}

	after, err := a.WriteModlist()
	if err != nil {
		t.Fatalf("WriteModlist() error = %v", err)
	}
	if after.Rows[2].Unlisted || after.Rows[2].Name != "bravo.archive" {
		t.Fatalf("row after Apply = %+v, want bravo.archive no longer marked new", after.Rows[2])
	}
	if !a.modlistSet["bravo.archive"] {
		t.Fatal("applied archive should be part of the in-memory modlist membership set")
	}

	want := []string{"alpha.archive", "charlie.archive", "bravo.archive"}
	if !slices.Equal(a.initialModlistOrder, want) {
		t.Fatalf("applied snapshot = %v, want %v", a.initialModlistOrder, want)
	}
	written, err := os.ReadFile(filepath.Join(tempDir, "modlist.txt"))
	if err != nil {
		t.Fatalf("read modlist.txt: %v", err)
	}
	if got := strings.Fields(string(written)); !slices.Equal(got, want) {
		t.Fatalf("written order = %v, want %v", got, want)
	}
}
