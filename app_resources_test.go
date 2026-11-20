package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"slices"
	"testing"

	"github.com/Mordeak/cp77-modorder-gui/internal/archive"
	"github.com/Mordeak/cp77-modorder-gui/internal/conflict"
)

func TestGetModResourcesGroupsCurrentOutcomes(t *testing.T) {
	const (
		uniqueHash = uint64(0x10)
		winHash    = uint64(0x20)
		lossHash   = uint64(0x30)
	)

	winner := &conflict.ModInfo{Name: "winner.archive"}
	selected := &conflict.ModInfo{
		Name: "selected.archive",
		Archive: &archive.Archive{
			Name:       "selected.archive",
			FileHashes: []uint64{uniqueHash, winHash, lossHash},
		},
		FileCount:     3,
		ConflictCount: 2,
		Wins:          1,
		Losses:        1,
	}
	loser := &conflict.ModInfo{Name: "loser.archive"}

	pathMap := map[string]string{
		resourceKey(uniqueHash): `base\unique.ent`,
		resourceKey(winHash):    `base\selected_wins.mesh`,
		resourceKey(lossHash):   `base\selected_loses.app`,
	}
	app := &App{
		result: &conflict.Result{
			Mods: []*conflict.ModInfo{winner, selected, loser},
			Conflicts: []*conflict.ConflictEntry{
				{Resource: resourceKey(winHash), Mods: []*conflict.ModInfo{selected, loser}},
				{Resource: resourceKey(lossHash), Mods: []*conflict.ModInfo{selected, winner}},
			},
		},
		pathMap:             pathMap,
		resourcePathsLoaded: map[string]bool{"selected.archive": true},
	}

	got, err := app.GetModResources("selected.archive")
	if err != nil {
		t.Fatalf("GetModResources() error = %v", err)
	}

	wantFiles := []string{`base\unique.ent`, `base\selected_wins.mesh`, `base\selected_loses.app`}
	if !slices.Equal(got.Files, wantFiles) {
		t.Fatalf("files = %v, want %v", got.Files, wantFiles)
	}
	wantConflicts := []string{`base\selected_wins.mesh`, `base\selected_loses.app`}
	if !slices.Equal(got.ConflictingFiles, wantConflicts) {
		t.Fatalf("conflicting files = %v, want %v", got.ConflictingFiles, wantConflicts)
	}
	if !slices.Equal(got.Wins, []string{`base\selected_wins.mesh`}) {
		t.Fatalf("wins = %v", got.Wins)
	}
	if !slices.Equal(got.Losses, []string{`base\selected_loses.app`}) {
		t.Fatalf("losses = %v", got.Losses)
	}

	// Moving the selected mod above the previous winner must update the lazy
	// outcome lists without rescanning or re-resolving resource paths.
	app.result.Mods = []*conflict.ModInfo{selected, winner, loser}
	got, err = app.GetModResources("selected.archive")
	if err != nil {
		t.Fatalf("GetModResources() after reorder error = %v", err)
	}
	if !slices.Equal(got.Wins, wantConflicts) {
		t.Fatalf("wins after reorder = %v, want %v", got.Wins, wantConflicts)
	}
	if len(got.Losses) != 0 {
		t.Fatalf("losses after reorder = %v, want none", got.Losses)
	}
}

func TestGetModResourcesFallsBackToHash(t *testing.T) {
	const hash = uint64(0x1234)
	selected := &conflict.ModInfo{
		Name: "selected.archive",
		Archive: &archive.Archive{
			Name:       "selected.archive",
			FileHashes: []uint64{hash},
		},
	}
	app := &App{
		result:              &conflict.Result{Mods: []*conflict.ModInfo{selected}},
		pathMap:             make(map[string]string),
		resourcePathsLoaded: map[string]bool{"selected.archive": true},
	}

	got, err := app.GetModResources("selected.archive")
	if err != nil {
		t.Fatalf("GetModResources() error = %v", err)
	}
	if want := resourceKey(hash); !slices.Equal(got.Files, []string{want}) {
		t.Fatalf("files = %v, want [%s]", got.Files, want)
	}
}

func TestWriteResourcesCSVExportsOneRowPerArchiveResource(t *testing.T) {
	const (
		uniqueHash   = uint64(0x100)
		conflictHash = uint64(0x200)
	)
	winner := &conflict.ModInfo{
		Name: "alpha.archive",
		Archive: &archive.Archive{
			Name:       "alpha.archive",
			FileHashes: []uint64{uniqueHash, conflictHash},
		},
	}
	loser := &conflict.ModInfo{
		Name: "beta.archive",
		Archive: &archive.Archive{
			Name:       "beta.archive",
			FileHashes: []uint64{conflictHash},
		},
	}
	app := &App{
		result: &conflict.Result{
			Mods: []*conflict.ModInfo{winner, loser},
			Conflicts: []*conflict.ConflictEntry{
				{Resource: resourceKey(conflictHash), Mods: []*conflict.ModInfo{loser, winner}},
			},
		},
		pathMap: map[string]string{
			resourceKey(uniqueHash):   `base\unique.ent`,
			resourceKey(conflictHash): `base\shared.mesh`,
		},
		resourcePathsLoaded: map[string]bool{
			"alpha.archive": true,
			"beta.archive":  true,
		},
	}

	var output bytes.Buffer
	if err := app.writeResourcesCSV(&output); err != nil {
		t.Fatalf("writeResourcesCSV() error = %v", err)
	}

	data := bytes.TrimPrefix(output.Bytes(), []byte{0xEF, 0xBB, 0xBF})
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		t.Fatalf("read exported CSV: %v", err)
	}
	if len(records) != 4 {
		t.Fatalf("CSV rows = %d, want 4 (header + 3 resources)", len(records))
	}

	wantHeader := []string{
		"Load Order",
		"Archive",
		"Resource Path",
		"Resource Hash",
		"Result",
		"Winning Archive",
		"Conflicting Archives",
	}
	if !slices.Equal(records[0], wantHeader) {
		t.Fatalf("CSV header = %v, want %v", records[0], wantHeader)
	}
	if got := records[1]; got[0] != "1" || got[1] != "alpha.archive" || got[2] != `base\unique.ent` || got[4] != "UNIQUE" {
		t.Fatalf("unique row = %v", got)
	}
	if got := records[2]; got[1] != "alpha.archive" || got[4] != "WIN" || got[5] != "alpha.archive" || got[6] != "beta.archive" {
		t.Fatalf("winning row = %v", got)
	}
	if got := records[3]; got[0] != "2" || got[1] != "beta.archive" || got[4] != "LOSS" || got[5] != "alpha.archive" || got[6] != "alpha.archive" {
		t.Fatalf("losing row = %v", got)
	}
}

func TestSafeCSVCellPreventsFormulaInterpretation(t *testing.T) {
	for _, value := range []string{"=formula", "+formula", "-formula", "@formula", "\tformula"} {
		if got := safeCSVCell(value); got != "'"+value {
			t.Fatalf("safeCSVCell(%q) = %q", value, got)
		}
	}
	if got := safeCSVCell(`base\safe.ent`); got != `base\safe.ent` {
		t.Fatalf("safeCSVCell changed safe value to %q", got)
	}
}

func resourceKey(hash uint64) string {
	return fmt.Sprintf("0x%016x", hash)
}
