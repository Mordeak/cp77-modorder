package main

// ConflictPairDTO represents one conflict pairing for IPC.
type ConflictPairDTO struct {
	Opponent string `json:"opponent"`
	Resource string `json:"resource"`
}

// ModDTO is the JSON-serialisable representation of a single mod.
type ModDTO struct {
	Name          string            `json:"name"`
	FileCount     int               `json:"fileCount"`
	Priority      int               `json:"priority"`
	ConflictCount int               `json:"conflictCount"`
	Wins          int               `json:"wins"`
	Losses        int               `json:"losses"`
	ConflictsWith []ConflictPairDTO `json:"conflictsWith"` // capped at 50
	HasMore       bool              `json:"hasMore"`
	MoreCount     int               `json:"moreCount"`
	HasXL         bool              `json:"hasXL"`
}

// DisplayRowDTO is one row in the mod list table.
// Mod is nil when the archive appears in modlist.txt but is absent from disk (Missing=true).
// Unlisted=true when the archive is on disk but absent from modlist.txt.
type DisplayRowDTO struct {
	Mod      *ModDTO `json:"mod"`      // nil = MISSING
	Name     string  `json:"name"`
	Unlisted bool    `json:"unlisted"`
	Missing  bool    `json:"missing"`
}

// ConflictDTO describes one resource that is contested by multiple mods.
type ConflictDTO struct {
	Resource string   `json:"resource"`
	Mods     []string `json:"mods"`
}

// ScanResultDTO is the full state returned after every scan or mutation.
type ScanResultDTO struct {
	Rows       []DisplayRowDTO `json:"rows"`
	Conflicts  []ConflictDTO   `json:"conflicts"`
	Summary    string          `json:"summary"`
	HasModlist bool            `json:"hasModlist"`
}

// ConflictGroupDTO holds the ordered conflict-group members for the DnD panel.
type ConflictGroupDTO struct {
	Mods []string `json:"mods"`
}

// ApplyPreviewDTO is the ordered list shown in the Apply dialog.
type ApplyPreviewDTO struct {
	Names   []string `json:"names"`   // new order that will be written
	Current []string `json:"current"` // order at last scan, for diff display
}

// ConfigDTO carries persisted user preferences.
type ConfigDTO struct {
	ModDir       string `json:"modDir"`
	ModStructure string `json:"modStructure"` // "default" | "MO2"
	MO2Dir       string `json:"mo2Dir"`
	MO2Profile   string `json:"mo2Profile"`
	BackupLimit  int    `json:"backupLimit"` // 0 = use default (20)
}
