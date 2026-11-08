package main

import (
	"path/filepath"
	"slices"
	"testing"
)

func TestReconcileScanOrderReinsertsVortexAppendedArchives(t *testing.T) {
	knownOrder := []string{
		"alpha.archive",
		"charlie.archive",
		"echo.archive",
	}
	a, _ := newReorderTestApp(t,
		knownOrder,
		[]string{
			"alpha.archive",
			"charlie.archive",
			"echo.archive",
			"bravo.archive",
			"delta.archive",
		},
	)
	a.cfg.KnownArchives = append([]string(nil), knownOrder...)
	a.cfg.ModlistSnapshot = append([]string(nil), knownOrder...)

	// Vortex preserves the known order and appends newly deployed archives.
	diskOrder := []string{
		"alpha.archive",
		"charlie.archive",
		"echo.archive",
		"bravo.archive",
		"delta.archive",
	}
	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	want := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
		"delta.archive",
		"echo.archive",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("reconciled order = %v, want %v", got, want)
	}
	if listed["bravo.archive"] || listed["delta.archive"] {
		t.Fatal("Vortex-appended archives should remain marked NEW before Apply")
	}
	for _, name := range knownOrder {
		if !listed[name] {
			t.Fatalf("known archive %s should remain listed", name)
		}
	}
}

func TestReconcileScanOrderAdoptsExternalAlphabeticalReset(t *testing.T) {
	oldOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"bravo.archive",
	}
	a, _ := newReorderTestApp(t, oldOrder, oldOrder)
	a.cfg.KnownArchives = append([]string(nil), oldOrder...)
	a.cfg.ModlistSnapshot = append([]string(nil), oldOrder...)

	diskOrder := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
	}
	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	if !slices.Equal(got, diskOrder) {
		t.Fatalf("reconciled order = %v, want external disk order %v", got, diskOrder)
	}
	for _, name := range diskOrder {
		if !listed[name] {
			t.Fatalf("known archive %s unexpectedly marked NEW", name)
		}
	}
}

func TestReconcileScanOrderReinsertsKnownArchiveReenabledByVortex(t *testing.T) {
	knownArchives := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
	}
	previousDiskOrder := []string{
		"alpha.archive",
		"charlie.archive",
	}
	a, _ := newReorderTestApp(t, previousDiskOrder, knownArchives)
	a.cfg.KnownArchives = append([]string(nil), knownArchives...)
	a.cfg.ModlistSnapshot = append([]string(nil), previousDiskOrder...)

	// Vortex appends the re-enabled archive even though MO saw it in the past.
	diskOrder := []string{
		"alpha.archive",
		"charlie.archive",
		"bravo.archive",
	}
	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	want := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("reconciled order = %v, want %v", got, want)
	}
	if !listed["bravo.archive"] {
		t.Fatal("a previously known re-enabled archive should not be marked NEW")
	}
}

func TestReconcileScanOrderRestoresUnappliedOrderWhenDiskIsUnchanged(t *testing.T) {
	diskOrder := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
	}
	proposedOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"bravo.archive",
	}
	a, mods := newReorderTestApp(t, diskOrder, proposedOrder)
	a.cfg.KnownArchives = append([]string(nil), diskOrder...)
	a.cfg.ModlistSnapshot = append([]string(nil), diskOrder...)
	a.cfg.Priorities = make(map[string]int, len(proposedOrder))
	for i, name := range proposedOrder {
		a.cfg.Priorities[name] = i + 1
		mods[name].Priority = i + 1
	}
	a.result.ApplyPriorities()

	got, _ := a.reconcileScanOrder(a.modDir, diskOrder, false)
	if !slices.Equal(got, proposedOrder) {
		t.Fatalf("reconciled order = %v, want saved proposed order %v", got, proposedOrder)
	}
}

func TestReconcileScanOrderMigratesFullLegacyPriorities(t *testing.T) {
	knownOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"echo.archive",
	}
	a, _ := newReorderTestApp(t,
		knownOrder,
		[]string{
			"charlie.archive",
			"alpha.archive",
			"echo.archive",
			"bravo.archive",
		},
	)
	// No KnownArchives or ModlistSnapshot: this is a config from v0.4.1.
	diskOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"echo.archive",
		"bravo.archive",
	}
	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	want := []string{
		"charlie.archive",
		"alpha.archive",
		"bravo.archive",
		"echo.archive",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("migrated order = %v, want %v", got, want)
	}
	if listed["bravo.archive"] {
		t.Fatal("archive absent from full legacy priorities should migrate as NEW")
	}
}

func TestReconcileScanOrderMigratesVortexResetWithAppendedSuffix(t *testing.T) {
	oldKnownOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"zulu.archive",
		"bravo.archive",
		"delta.archive",
	}
	a, _ := newReorderTestApp(t, oldKnownOrder, oldKnownOrder)
	// Vortex first reset the known archives alphabetically, then appended a
	// re-enabled archive. No explicit snapshot exists in this legacy config yet.
	diskOrder := []string{
		"alpha.archive",
		"bravo.archive",
		"delta.archive",
		"zulu.archive",
		"charlie.archive",
	}
	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	want := []string{
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
		"delta.archive",
		"zulu.archive",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("migrated order = %v, want %v", got, want)
	}
	if !listed["charlie.archive"] {
		t.Fatal("a known archive appended after reset should not be marked NEW")
	}
}

func TestReconcileScanOrderTrustsDiskForSparseLegacyPriorities(t *testing.T) {
	diskOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"bravo.archive",
		"delta.archive",
	}
	a, mods := newReorderTestApp(t, diskOrder, diskOrder)
	a.cfg.Priorities = map[string]int{"charlie.archive": 1}
	for name, mod := range mods {
		if name == "charlie.archive" {
			mod.Priority = 1
		} else {
			mod.Priority = 0
		}
	}
	a.result.ApplyPriorities()

	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	if !slices.Equal(got, diskOrder) {
		t.Fatalf("migrated order = %v, want disk order %v", got, diskOrder)
	}
	for _, name := range diskOrder {
		if !listed[name] {
			t.Fatalf("sparse legacy migration incorrectly marked %s NEW", name)
		}
	}
}

func TestReconcileScanOrderDoesNotReuseStateFromAnotherFolder(t *testing.T) {
	diskOrder := []string{
		"charlie.archive",
		"alpha.archive",
		"bravo.archive",
	}
	a, _ := newReorderTestApp(t, diskOrder, diskOrder)
	a.cfg.ArchiveStateDir = filepath.Join(t.TempDir(), "other-mod-folder")
	a.cfg.KnownArchives = []string{"alpha.archive"}
	a.cfg.ModlistSnapshot = []string{"alpha.archive"}

	got, listed := a.reconcileScanOrder(a.modDir, diskOrder, false)
	if !slices.Equal(got, diskOrder) {
		t.Fatalf("order = %v, want current folder disk order %v", got, diskOrder)
	}
	for _, name := range diskOrder {
		if !listed[name] {
			t.Fatalf("state from another folder incorrectly marked %s NEW", name)
		}
	}
}
