# cp77-conflict — Cyberpunk 2077 Mod Conflict Manager

Detects hash conflicts between `.archive` mods and manages the `modlist.txt` load order.
Built with **Wails v2** (Go backend + Vue 3 frontend), targeting Windows amd64.

## Project layout

```
cp77-conflict/
├── main.go                        # Wails app init; embeds frontend/dist
├── app.go                         # All Go methods bound to JS (App struct)
├── dto.go                         # Data transfer objects for IPC
├── wails.json                     # Wails project config
├── Makefile                       # Build automation
├── internal/
│   ├── archive/archive.go         # RDAR binary parser (FNV-1a hashes + LXRS paths)
│   ├── conflict/conflict.go       # Conflict detection + priority sort
│   ├── config/config.go           # JSON persistence (%APPDATA%/cp77-modorder)
│   └── modlist/modlist.go         # modlist.txt writer with timestamped backups
└── frontend/
    ├── src/
    │   ├── App.vue                # Root layout (split panel + dialogs)
    │   ├── main.ts                # Vue + Pinia entry point
    │   ├── style.css              # Global CSS variables (--cp-* theme tokens)
    │   ├── utils.ts               # truncate() and other helpers
    │   ├── components/
    │   │   ├── ModListTable.vue   # Mod list (vuedraggable + filter/search)
    │   │   ├── DetailPanel.vue    # Right pane: selected mod detail
    │   │   ├── Toolbar.vue        # Top action buttons
    │   │   ├── PathBar.vue        # Folder path + scan controls
    │   │   ├── StatusBar.vue      # Bottom summary line
    │   │   ├── ConflictGraph.vue  # Modal: all conflicts table
    │   │   ├── ApplyDialog.vue    # Modal: preview + write modlist.txt
    │   │   ├── ConflictResolutionDialog.vue  # Modal: handle new conflicting mods
    │   │   ├── RestoreDialog.vue  # Modal: restore modlist backup
    │   │   ├── LoadOrderDnd.vue   # Drag-and-drop subcomponent for conflict groups
    │   │   └── PriorityBadge.vue  # Small badge showing priority number
    │   ├── composables/
    │   │   └── useWails.ts        # Wraps all Go IPC calls + registers events
    │   └── stores/
    │       └── app.ts             # Pinia store (global scan result, selection, dialogs)
    ├── wailsjs/                   # Auto-generated Wails JS bindings (do not edit)
    │   ├── go/main/App.js / App.d.ts
    │   └── go/models.ts
    ├── vite.config.ts
    ├── tsconfig.json
    └── package.json
```

## Dependencies

**Go:**
- `github.com/wailsapp/wails/v2 v2.12.0`

**Frontend (package.json):**
- Vue 3.4, Pinia 2.1, vuedraggable 4.1
- Vite 5.2, TypeScript 5.3, vue-tsc

Do not add CGO-heavy libs, web renderers, or anything that complicates Windows cross-compilation.

## Build instructions

```bash
# Dev mode (hot-reload frontend + Go IPC)
make dev          # or: wails dev

# Release build
make build        # or: wails build -platform windows/amd64 -o CP77-modorder.exe -trimpath -ldflags="-s -w"

# Regenerate Wails JS bindings after changing App struct methods
make generate     # or: wails generate module

# Dependency hygiene
make tidy         # go mod tidy + npm install
```

CGO is required. The Makefile sets `CC=C:/msys64/ucrt64/bin/gcc.exe`. The `-H windowsgui` flag is applied automatically by Wails via `wails.json` platforms config.

## Architecture

### Go backend — app.go / dto.go

`App` struct holds all mutable state (config, conflict result, modlist order). All exported methods are automatically bound to JavaScript by Wails.

**Public methods (callable from JS via `useWails.ts`):**

| Method | Description |
|---|---|
| `GetConfig()` | Returns current mod folder path |
| `PickFolder()` | Opens native folder picker |
| `Scan(dir)` | Scans for .archive files, detects conflicts, returns `ScanResultDTO` |
| `GetApplyPreview()` | Preview of modlist.txt order |
| `SetPriority(name, p)` | Set priority 1–99 (0 = unset), recomputes, returns updated result |
| `GetConflictGroup(name)` | Mods sharing conflicts with the named mod |
| `ReorderConflictGroup(names)` | Reorder a conflict group, preserve global order |
| `GroupConflicts()` | Cluster connected conflict components together |
| `SetModlistOrder(names)` | Apply drag-drop order (in-memory until WriteModlist) |
| `WriteModlist()` | Write modlist.txt to disk (auto-backup to modlist.old/) |
| `ListBackups()` | List timestamped backups |
| `RestoreBackup(filename)` | Restore a backup and recompute conflicts |

Events emitted via `runtime.EventsEmit`: `scan:progress`.

### Frontend — Vue 3 + Pinia

**Data flow:**
1. `useWails.ts` calls a Go method → receives `ScanResultDTO`
2. Calls `store.setScanResult()` → Pinia updates global state
3. Components read from the store reactively

**Store (stores/app.ts) key state:**
- `scanResult` — the current `ScanResultDTO` (rows, conflicts, summary)
- `selectedIndex` — currently selected row
- `modDir` — current mod folder
- Dialog visibility flags

**useWails.ts** is the single point of contact with Go — all IPC goes through it.

## Data transfer objects — dto.go

```go
ScanResultDTO  { Rows, Conflicts, Summary, HasModlist }
DisplayRowDTO  { Mod *ModDTO, Name, Unlisted, Missing }
ModDTO         { Name, FileCount, Priority, ConflictCount, Wins, Losses, ConflictsWith []ConflictPairDTO, HasMore }
ConflictDTO    { Resource, Mods []string }
ConflictPairDTO { Opponent, Resource }
ConflictGroupDTO { Mods []string }
ApplyPreviewDTO  { Names []string }
ConfigDTO        { ModDir string }
```

TypeScript types are auto-generated into `frontend/wailsjs/go/models.ts` — do not edit that file directly.

## Visual design

Theme is implemented via CSS custom properties in `frontend/src/style.css`:

| Token | Hex | Role |
|---|---|---|
| `--cp-bg` | `#0D0D0D` | Background |
| `--cp-fg` | `#EEEEEE` | Text |
| `--cp-primary` | `#F0C000` | Accent / priority |
| `--cp-focus` | `#00E5FF` | Selection / focus |
| `--cp-error` | `#FF3366` | Conflicts / losses |
| `--cp-success` | `#39FF14` | Wins |
| `--cp-input` | `#1A1A1A` | Input / button background |
| `--cp-dim` | (muted grey) | Secondary text |

## Core logic rules

The four `internal/` packages encode correct game behaviour. Edit carefully and test thoroughly.

1. **Load order semantics**: first entry in modlist.txt loads first and wins conflicts.
2. **Sort rules in `conflict.go`**:
   - Priority mods first; lower number = earlier load.
   - Tie-break: higher ConflictCount (maximises wins).
   - Unset (0) mods go last, sorted alphabetically.
3. **`modlist.Write()`** always backs up the existing file before overwriting.
4. **Archive format**: magic bytes `RDAR`, reads FNV1a-64 hashes from index starting at `indexOffset`; skip 24 bytes per entry after the hash. Optional LXRS footer at offset `0xAC` provides human-readable resource paths.
5. **Silent parse failures**: one bad archive must not abort the scan — log and continue.

## State management

All mutable server-side state is owned by the `App` struct in `app.go`. Frontend state lives in the Pinia store (`stores/app.ts`). Do not use globals in either layer.

After any operation that changes mod order or priorities, `app.go` must call `buildScanResult()` and return the new `ScanResultDTO` to the frontend.

## What NOT to do

- Do not edit files under `frontend/wailsjs/` — they are code-generated by `wails generate module`.
- Do not call `os.Exit` anywhere except `main.go` on fatal startup error.
- Do not add a splash screen, auto-updater, or telemetry.
- Do not hard-code the mod directory — always load from `config.Config.ModDir`.
- Do not bypass the DTO layer — Go types must not be exposed directly to JS.
