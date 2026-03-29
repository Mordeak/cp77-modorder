package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/Mordeak/cp77-modorder-gui/internal/archive"
	"github.com/Mordeak/cp77-modorder-gui/internal/conflict"
	"github.com/Mordeak/cp77-modorder-gui/internal/config"
	"github.com/Mordeak/cp77-modorder-gui/internal/modlist"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails application struct. All exported methods are auto-bound to JS.
type App struct {
	ctx          context.Context
	cfg          *config.Config
	cfgPath      string
	result       *conflict.Result
	modDir       string
	modlistOrder []string
	modlistSet   map[string]bool
	pathMap      map[string]string // "0x<hex>" → human-readable resource path (from LXRS footers)
}

// NewApp creates the App with a loaded config.
func NewApp() *App {
	cfgPath := config.DefaultPath()
	cfg, err := config.Load(cfgPath)
	if err != nil {
		cfg = &config.Config{Priorities: make(map[string]int)}
	}
	return &App{
		cfg:     cfg,
		cfgPath: cfgPath,
		modDir:  cfg.ModDir,
	}
}

// startup is called when Wails is ready; saves the context for runtime calls.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ---- Bound methods --------------------------------------------------------

// GetConfig returns persisted config for use on startup.
func (a *App) GetConfig() ConfigDTO {
	return ConfigDTO{ModDir: a.cfg.ModDir}
}

// PickFolder opens the native OS folder picker and returns the chosen path.
func (a *App) PickFolder() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Cyberpunk 2077 mods folder",
	})
	if err != nil {
		return ""
	}
	return dir
}

// Scan scans dir for .archive files, detects conflicts, and returns the full result.
// It also saves modDir to config. Progress events are emitted as "scan:progress".
func (a *App) Scan(dir string) (ScanResultDTO, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return ScanResultDTO{}, fmt.Errorf("no folder path provided")
	}

	runtime.EventsEmit(a.ctx, "scan:progress", "Scanning "+dir+" …")

	archives, err := archive.Scan(dir)
	if err != nil {
		return ScanResultDTO{}, err
	}
	if len(archives) == 0 {
		// Count raw .archive files to give a better error message.
		var rawCount int
		if entries, readErr := os.ReadDir(dir); readErr == nil {
			for _, e := range entries {
				if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".archive") {
					rawCount++
				}
			}
		}
		if rawCount > 0 {
			return ScanResultDTO{}, fmt.Errorf(
				"found %d .archive file(s) in %s but none could be read as valid RED4 archives",
				rawCount, dir,
			)
		}
		return ScanResultDTO{}, fmt.Errorf("no .archive files found in %s", dir)
	}

	a.result = conflict.Detect(archives, a.cfg.Priorities)

	// Build a merged hash→path map from all archives that have LXRS data.
	pathMap := make(map[string]string)
	for _, ar := range archives {
		for hash, path := range ar.FilePaths {
			pathMap[fmt.Sprintf("0x%016x", hash)] = path
		}
	}
	a.pathMap = pathMap

	order, _ := a.readModlistOrder(dir)
	a.modlistOrder = order
	a.modlistSet = make(map[string]bool, len(order))
	for _, n := range order {
		a.modlistSet[n] = true
	}

	// Persist the directory.
	a.modDir = dir
	a.cfg.ModDir = dir
	_ = a.cfg.Save(a.cfgPath)

	runtime.EventsEmit(a.ctx, "scan:progress", a.result.Summary())
	return a.buildScanResult(), nil
}

// SetPriority sets the priority for mod name (0 = clear) and returns the updated result.
func (a *App) SetPriority(name string, p int) (ScanResultDTO, error) {
	if a.result == nil {
		return ScanResultDTO{}, fmt.Errorf("no scan results — run a scan first")
	}
	if p < 0 || p > 99 {
		return ScanResultDTO{}, fmt.Errorf("priority must be 0–99")
	}
	for _, m := range a.result.Mods {
		if m.Name == name {
			m.Priority = p
			if p == 0 {
				delete(a.cfg.Priorities, name)
			} else {
				a.cfg.Priorities[name] = p
			}
			break
		}
	}
	_ = a.cfg.Save(a.cfgPath)
	a.result.ApplyPriorities()
	return a.buildScanResult(), nil
}

// GetConflictGroup returns the ordered list of mods that share at least one
// conflicting resource with the named mod.
func (a *App) GetConflictGroup(name string) (ConflictGroupDTO, error) {
	if a.result == nil {
		return ConflictGroupDTO{}, fmt.Errorf("no scan results")
	}
	var anchor *conflict.ModInfo
	for _, m := range a.result.Mods {
		if m.Name == name {
			anchor = m
			break
		}
	}
	if anchor == nil {
		return ConflictGroupDTO{}, fmt.Errorf("mod %q not found", name)
	}
	group := a.conflictGroup(anchor)
	names := make([]string, len(group))
	for i, m := range group {
		names[i] = m.Name
	}
	return ConflictGroupDTO{Mods: names}, nil
}

// ReorderConflictGroup reorders the conflict-group mods relative to each
// other without moving them to the top of the global load order.
//
// Strategy: the group mods collectively occupy a set of "slots" in the
// current sorted list. We redistribute those slots to the mods in the
// user's requested order instead of blindly assigning priorities 1, 2, 3.
//
// Special case: if every group mod is currently unset (priority 0) we give
// them explicit priorities placed just after the last explicitly-prioritised
// non-group mod, so they stay near the end of the priority zone rather than
// jumping to the very top.
func (a *App) ReorderConflictGroup(names []string) (ScanResultDTO, error) {
	if a.result == nil {
		return ScanResultDTO{}, fmt.Errorf("no scan results")
	}

	// Build name→mod lookup.
	nameIdx := make(map[string]*conflict.ModInfo, len(a.result.Mods))
	for _, m := range a.result.Mods {
		nameIdx[m.Name] = m
	}

	// Current position of every mod in the already-sorted slice.
	posOf := make(map[string]int, len(a.result.Mods))
	for i, m := range a.result.Mods {
		posOf[m.Name] = i
	}

	// Sort the group mods by their current global position so we know which
	// priority "slot" each one occupies.
	type slot struct {
		name string
		pos  int
	}
	slots := make([]slot, len(names))
	for i, name := range names {
		slots[i] = slot{name, posOf[name]}
	}
	sort.Slice(slots, func(i, j int) bool { return slots[i].pos < slots[j].pos })

	// Collect the priority values that the slots currently hold, in position order.
	slotPrios := make([]int, len(slots))
	for i, s := range slots {
		slotPrios[i] = nameIdx[s.name].Priority
	}

	// If every group mod is unset (priority 0) we cannot preserve relative
	// ordering through 0 alone (sortMods would just re-alphabetise them).
	// Instead, anchor the group just after the last explicitly-prioritised
	// non-group mod so they stay at the bottom of the priority zone.
	allZero := true
	for _, p := range slotPrios {
		if p != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		groupSet := make(map[string]bool, len(names))
		for _, name := range names {
			groupSet[name] = true
		}
		maxPrio := 0
		for _, m := range a.result.Mods {
			if !groupSet[m.Name] && m.Priority > maxPrio {
				maxPrio = m.Priority
			}
		}
		for i := range slotPrios {
			slotPrios[i] = maxPrio + i + 1
		}
	}

	// Redistribute: the i-th mod in the user's new order inherits the priority
	// that belonged to the i-th slot in the current (position-sorted) order.
	for i, name := range names {
		if m, ok := nameIdx[name]; ok {
			m.Priority = slotPrios[i]
			a.cfg.Priorities[name] = slotPrios[i]
		}
	}

	_ = a.cfg.Save(a.cfgPath)
	a.result.ApplyPriorities()
	return a.buildScanResult(), nil
}

// GetApplyPreview returns the ordered mod names that will be written to modlist.txt.
func (a *App) GetApplyPreview() (ApplyPreviewDTO, error) {
	if a.result == nil {
		return ApplyPreviewDTO{}, fmt.Errorf("no scan results — run a scan first")
	}
	names := make([]string, len(a.result.Mods))
	for i, m := range a.result.Mods {
		names[i] = m.Name
	}
	return ApplyPreviewDTO{Names: names}, nil
}

// WriteModlist writes modlist.txt to the current mod directory.
func (a *App) WriteModlist() error {
	if a.result == nil {
		return fmt.Errorf("no scan results — run a scan first")
	}
	if a.modDir == "" {
		return fmt.Errorf("no mod directory set")
	}
	return modlist.Write(a.modDir, a.result.Mods)
}

// ---- Private helpers -------------------------------------------------------

// buildScanResult converts internal state to a ScanResultDTO.
// Row order follows modlist.txt when present; otherwise result.Mods order.
func (a *App) buildScanResult() ScanResultDTO {
	if a.result == nil {
		return ScanResultDTO{}
	}

	var rows []DisplayRowDTO

	if len(a.modlistOrder) == 0 {
		// No modlist.txt — use result.Mods order.
		for _, m := range a.result.Mods {
			rows = append(rows, DisplayRowDTO{Mod: modToDTO(m, a.pathMap), Name: m.Name})
		}
	} else {
		modByName := make(map[string]*conflict.ModInfo, len(a.result.Mods))
		for _, m := range a.result.Mods {
			modByName[m.Name] = m
		}
		// modlist.txt entries first.
		for _, name := range a.modlistOrder {
			if m, ok := modByName[name]; ok {
				rows = append(rows, DisplayRowDTO{Mod: modToDTO(m, a.pathMap), Name: name})
			} else {
				rows = append(rows, DisplayRowDTO{Mod: nil, Name: name, Missing: true})
			}
		}
		// Unlisted mods (on disk, not in modlist.txt) at the bottom.
		for _, m := range a.result.Mods {
			if !a.modlistSet[m.Name] {
				rows = append(rows, DisplayRowDTO{Mod: modToDTO(m, a.pathMap), Name: m.Name, Unlisted: true})
			}
		}
	}

	// Build conflict list, resolving hex hashes to human-readable paths where available.
	conflicts := make([]ConflictDTO, 0, len(a.result.Conflicts))
	for _, ce := range a.result.Conflicts {
		mods := make([]string, len(ce.Mods))
		for i, m := range ce.Mods {
			mods[i] = m.Name
		}
		res := ce.Resource
		if human, ok := a.pathMap[res]; ok {
			res = human
		}
		conflicts = append(conflicts, ConflictDTO{Resource: res, Mods: mods})
	}

	return ScanResultDTO{
		Rows:       rows,
		Conflicts:  conflicts,
		Summary:    a.result.Summary(),
		HasModlist: len(a.modlistOrder) > 0,
	}
}

// modToDTO converts a ModInfo to ModDTO, capping ConflictsWith at 50 entries.
// pathMap resolves hex resource keys to human-readable paths where available.
func modToDTO(m *conflict.ModInfo, pathMap map[string]string) *ModDTO {
	pairs := m.ConflictsWith
	hasMore := len(pairs) > 50
	moreCount := 0
	if hasMore {
		moreCount = len(pairs) - 50
		pairs = pairs[:50]
	}
	cw := make([]ConflictPairDTO, len(pairs))
	for i, p := range pairs {
		res := p.Resource
		if human, ok := pathMap[res]; ok {
			res = human
		}
		cw[i] = ConflictPairDTO{Opponent: p.Opponent.Name, Resource: res}
	}
	return &ModDTO{
		Name:          m.Name,
		FileCount:     m.FileCount,
		Priority:      m.Priority,
		ConflictCount: m.ConflictCount,
		Wins:          m.Wins,
		Losses:        m.Losses,
		ConflictsWith: cw,
		HasMore:       hasMore,
		MoreCount:     moreCount,
	}
}

// readModlistOrder reads modlist.txt from dir and returns the ordered mod names.
// Returns nil (not an error) when the file does not exist.
func (a *App) readModlistOrder(dir string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "modlist.txt"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			names = append(names, line)
		}
	}
	return names, nil
}

// conflictGroup returns all mods that share at least one conflicting resource with m,
// sorted by their current position in result.Mods (load order).
func (a *App) conflictGroup(m *conflict.ModInfo) []*conflict.ModInfo {
	seen := make(map[*conflict.ModInfo]bool)
	for _, ce := range a.result.Conflicts {
		if slices.Contains(ce.Mods, m) {
			for _, cm := range ce.Mods {
				seen[cm] = true
			}
		}
	}
	out := make([]*conflict.ModInfo, 0, len(seen))
	for mod := range seen {
		out = append(out, mod)
	}
	posOf := make(map[*conflict.ModInfo]int, len(a.result.Mods))
	for i, mod := range a.result.Mods {
		posOf[mod] = i
	}
	sort.Slice(out, func(i, j int) bool {
		return posOf[out[i]] < posOf[out[j]]
	})
	return out
}
