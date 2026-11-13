package main

import "testing"

func TestMarkCustomOrderLeavesASCIIOrderUnmarked(t *testing.T) {
	rows := customOrderTestRows("alpha.archive", "bravo.archive", "charlie.archive")
	markCustomOrder(rows)
	assertCustomRows(t, rows)
}

func TestMarkCustomOrderMarksOnlyArchiveMovedUp(t *testing.T) {
	rows := customOrderTestRows(
		"alpha.archive",
		"delta.archive",
		"bravo.archive",
		"charlie.archive",
		"echo.archive",
	)
	markCustomOrder(rows)
	assertCustomRows(t, rows, "delta.archive")
}

func TestMarkCustomOrderMarksOnlyArchiveMovedDown(t *testing.T) {
	rows := customOrderTestRows(
		"bravo.archive",
		"charlie.archive",
		"alpha.archive",
		"delta.archive",
	)
	markCustomOrder(rows)
	assertCustomRows(t, rows, "alpha.archive")
}

func TestMarkCustomOrderAllowsNewAndCustomTogether(t *testing.T) {
	rows := customOrderTestRows(
		"delta.archive",
		"alpha.archive",
		"bravo.archive",
		"charlie.archive",
	)
	rows[0].Unlisted = true
	markCustomOrder(rows)
	assertCustomRows(t, rows, "delta.archive")
	if !rows[0].Unlisted || !rows[0].CustomOrder {
		t.Fatalf("row = %+v, want both NEW and CUSTOM states", rows[0])
	}
}

func TestMarkCustomOrderIgnoresMissingRows(t *testing.T) {
	rows := customOrderTestRows("alpha.archive", "bravo.archive")
	rows = append(rows[:1], append([]DisplayRowDTO{{Name: "z-missing.archive", Missing: true}}, rows[1:]...)...)
	markCustomOrder(rows)
	assertCustomRows(t, rows)
}

func customOrderTestRows(names ...string) []DisplayRowDTO {
	rows := make([]DisplayRowDTO, len(names))
	for i, name := range names {
		rows[i] = DisplayRowDTO{Name: name, Mod: &ModDTO{Name: name}}
	}
	return rows
}

func assertCustomRows(t *testing.T, rows []DisplayRowDTO, want ...string) {
	t.Helper()
	wantSet := nameSet(want)
	for _, row := range rows {
		if row.CustomOrder != wantSet[row.Name] {
			t.Fatalf("row %s custom=%v, want %v", row.Name, row.CustomOrder, wantSet[row.Name])
		}
	}
}
