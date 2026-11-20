package main

import (
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

func resourceKey(hash uint64) string {
	return fmt.Sprintf("0x%016x", hash)
}
