// Package conflict detects hash collisions between archives and manages load order.
package conflict

import (
	"fmt"
	"sort"

	"github.com/Mordeak/cp77-modorder-gui/internal/archive"
)

// ConflictPair records a single resource hash shared between two mods.
type ConflictPair struct {
	Opponent *ModInfo
	Resource string // "0x<hex>" when no path mapping is known
}

// ModInfo enriches an archive with conflict analysis and user priority.
type ModInfo struct {
	Archive       *archive.Archive
	Name          string // shorthand: Archive.Name
	FileCount     int
	Priority      int // 1-99; 0 = unset
	ConflictCount int
	Wins          int
	Losses        int
	ConflictsWith []*ConflictPair
}

// ConflictEntry is one resource hash contested by two or more mods.
type ConflictEntry struct {
	Resource string     // "0x<hex>"
	Mods     []*ModInfo // all mods containing this hash, in load order
}

// Result is the output of Detect — the complete conflict analysis.
type Result struct {
	Mods      []*ModInfo      // in current load order
	Conflicts []*ConflictEntry // all contested resources
}

// Detect runs conflict analysis over archives, applying priorities from the map.
func Detect(archives []*archive.Archive, priorities map[string]int) *Result {
	// Build ModInfo list.
	mods := make([]*ModInfo, len(archives))
	for i, a := range archives {
		mods[i] = &ModInfo{
			Archive:   a,
			Name:      a.Name,
			FileCount: len(a.FileHashes),
			Priority:  priorities[a.Name],
		}
	}

	// Index: hash → list of mods containing it.
	index := make(map[uint64][]*ModInfo)
	for _, m := range mods {
		for _, h := range m.Archive.FileHashes {
			index[h] = append(index[h], m)
		}
	}

	// Apply initial sort before computing wins/losses.
	sortMods(mods)

	// Build conflict entries and win/loss counts.
	var conflicts []*ConflictEntry
	seen := make(map[uint64]bool)
	for _, m := range mods {
		for _, h := range m.Archive.FileHashes {
			if seen[h] {
				continue
			}
			contestants := index[h]
			if len(contestants) < 2 {
				continue
			}
			seen[h] = true
			resource := fmt.Sprintf("0x%016x", h)
			entry := &ConflictEntry{Resource: resource, Mods: contestants}
			conflicts = append(conflicts, entry)

			// Winner is the contestant earliest in load order (lowest index in mods).
			winnerIdx := loadIdx(mods, contestants[0])
			winner := contestants[0]
			for _, c := range contestants[1:] {
				if idx := loadIdx(mods, c); idx < winnerIdx {
					winnerIdx = idx
					winner = c
				}
			}
			for _, c := range contestants {
				c.ConflictCount++
				if c == winner {
					c.Wins++
				} else {
					c.Losses++
					c.ConflictsWith = append(c.ConflictsWith, &ConflictPair{
						Opponent: winner,
						Resource: resource,
					})
				}
			}
		}
	}

	return &Result{Mods: mods, Conflicts: conflicts}
}

// ApplyPriorities re-sorts Mods and resets win/loss counts.
// Call this after any priority change.
func (r *Result) ApplyPriorities() {
	// Reset counters.
	for _, m := range r.Mods {
		m.ConflictCount = 0
		m.Wins = 0
		m.Losses = 0
		m.ConflictsWith = nil
	}
	sortMods(r.Mods)

	// Recompute wins/losses based on new order.
	index := make(map[string][]*ModInfo) // resource → mods
	for _, ce := range r.Conflicts {
		index[ce.Resource] = ce.Mods
	}
	for _, ce := range r.Conflicts {
		winnerIdx := loadIdx(r.Mods, ce.Mods[0])
		winner := ce.Mods[0]
		for _, m := range ce.Mods[1:] {
			if idx := loadIdx(r.Mods, m); idx < winnerIdx {
				winnerIdx = idx
				winner = m
			}
		}
		for _, m := range ce.Mods {
			m.ConflictCount++
			if m == winner {
				m.Wins++
			} else {
				m.Losses++
				m.ConflictsWith = append(m.ConflictsWith, &ConflictPair{
					Opponent: winner,
					Resource: ce.Resource,
				})
			}
		}
	}
}

// Summary returns a one-line human-readable summary of the result.
func (r *Result) Summary() string {
	conflicted := 0
	for _, m := range r.Mods {
		if m.ConflictCount > 0 {
			conflicted++
		}
	}
	return fmt.Sprintf("%d mods loaded · %d conflicts · %d mods affected",
		len(r.Mods), len(r.Conflicts), conflicted)
}

// sortMods sorts in-place:
//  1. Priority mods first (lower number = earlier).
//  2. Tie-break: higher ConflictCount first.
//  3. Unset (0) mods last, alphabetical.
func sortMods(mods []*ModInfo) {
	sort.SliceStable(mods, func(i, j int) bool {
		a, b := mods[i], mods[j]
		aPrio, bPrio := a.Priority > 0, b.Priority > 0
		switch {
		case aPrio && !bPrio:
			return true
		case !aPrio && bPrio:
			return false
		case aPrio && bPrio:
			if a.Priority != b.Priority {
				return a.Priority < b.Priority
			}
			return a.ConflictCount > b.ConflictCount
		default: // both unset
			return a.Name < b.Name
		}
	})
}

func loadIdx(mods []*ModInfo, target *ModInfo) int {
	for i, m := range mods {
		if m == target {
			return i
		}
	}
	return len(mods)
}
