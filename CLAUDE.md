# cp77-modorder — GUI Edition
Cyberpunk 2077 .archive conflict manager with priority-based auto load order.
GUI built with Fyne v2, targeting Windows (amd64) as the primary platform.

## Project layout
```
cp77-modorder-gui/
├── CLAUDE.md
├── go.mod
├── main.go                        # Entry point: app init, window creation
├── internal/
│   ├── archive/archive.go         # RED4 binary parser (DO NOT MODIFY logic)
│   ├── conflict/conflict.go       # Conflict detection + priority sort (DO NOT MODIFY logic)
│   ├── config/config.go           # JSON persistence (DO NOT MODIFY logic)
│   └── modlist/modlist.go         # modlist.txt writer (DO NOT MODIFY logic)
└── gui/
    ├── app.go                     # Root Fyne app + window wiring
    ├── theme.go                   # Custom Cyberpunk-inspired Fyne theme
    ├── modlist_view.go            # Main mod list panel (left pane)
    ├── detail_view.go             # Conflict detail panel (right pane)
    ├── conflict_view.go           # Full conflict graph tab
    ├── apply_view.go              # Apply / write modlist.txt tab
    └── widgets/
        ├── priority_badge.go      # Small coloured badge showing priority
        └── conflict_bar.go        # Win/Lose bar widget
```
Copy the four internal/ packages verbatim from the TUI project — they have no TUI dependencies and require zero changes.

## Dependencies
- `fyne.io/fyne/v2 v2.5.x`
- Add nothing else. Do not add a web renderer, CGO-heavy libs, or anything that complicates cross-compilation to Windows.

## go.mod bootstrap
```
module github.com/Mordeak/cp77-modorder-gui

go 1.22

require fyne.io/fyne/v2 v2.5.3
```
After writing go.mod, run:
```bash
go mod tidy
```

## Build instructions
```bash
# Local build (runs on current OS)
go build -o cp77-modorder-gui.exe ./main.go

# Windows cross-compile from Linux
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-H windowsgui" -o cp77-modorder-gui.exe ./main.go
```
The `-H windowsgui` linker flag suppresses the console window on Windows — always include it in release builds.

## Visual design

### Theme — gui/theme.go
Implement a custom `fyne.Theme` named `CyberpunkTheme`. Key colours:

| Role | Hex |
|---|---|
| Background | #0D0D0D |
| Foreground / text | #EEEEEE |
| Primary (accent) | #F0C000 (yellow) |
| Focus / selection | #00E5FF (cyan) |
| Error / conflict | #FF3366 (red) |
| Success / win | #39FF14 (green) |
| Button | #1A1A1A |
| Input background | #1A1A1A |

Font: use Fyne's built-in monospace for data rows; default sans-serif elsewhere. Do not bundle external fonts unless asked.

Apply the theme immediately in main.go:
```go
a := app.New()
a.Settings().SetTheme(gui.NewCyberpunkTheme())
```

## Screen / layout specification

### Main window — gui/app.go
- Fixed minimum size: 1100 × 650.
- Layout: `container.NewHSplit` with left pane (mod list) and right pane (detail).
- HSplit ratio: 60 % left / 40 % right.
- Bottom: a `widget.Label` status bar showing `result.Summary()`.
- Top: a `widget.Toolbar` with these actions:

| Icon | Label | Action |
|---|---|---|
| `theme.FolderOpenIcon()` | Open folder | Open native folder picker → scan |
| `theme.ViewRefreshIcon()` | Rescan | Re-run scan on current dir |
| `theme.DocumentSaveIcon()` | Apply | Write modlist.txt |
| `theme.InfoIcon()` | Conflicts | Switch right pane to conflict graph |

### Left pane — Mod list — gui/modlist_view.go
Use `widget.NewTable` (not `widget.List`) — it handles large mod counts efficiently.

**Columns**

| # | Header | Width | Content |
|---|---|---|---|
| 0 | # | 45 | Load position (1-based) |
| 1 | Priority | 75 | Priority badge widget or — |
| 2 | Mod | 340 | Archive name (truncated if >45 chars) |
| 3 | Files | 65 | Total internal file count |
| 4 | Conflicts | 90 | Conflict count, red if >0 |
| 5 | W / L | 90 | Win / Lose counts, coloured |

**Behaviour**
- Clicking a row selects it and updates the right detail pane.
- Double-clicking a row opens an `inline dialog.NewForm` for setting the priority:
  - One `widget.Entry` labelled "Priority (1–99, 0 = unset)".
  - On confirm: update `ModInfo.Priority`, save config, call `result.ApplyPriorities()`, refresh the table.
- Rows with `ConflictCount > 0` render the mod name in the conflict colour (#FF3366).
- The selected row is highlighted with the focus colour (#00E5FF) background.
- After `ApplyPriorities()` re-sorts the list, scroll the table so the previously selected mod stays visible.

### Right pane — Detail panel — gui/detail_view.go
Shown when a mod row is selected. Use `container.NewVBox` with:

1. Header: mod name in bold, large text.
2. Stats grid (`widget.Form`):
   - Files in archive
   - Conflicting files (red if > 0)
   - Wins (green)
   - Losses (red)
   - Priority (yellow, editable via button next to it)
3. Conflict list (`widget.List`): each entry shows:
   - The competing mod name (red)
   - The resource path or `0x<hex>` hash (dim, monospace)
   - Show max 50 entries; if more, append a "…and N more" label.
4. A "Set Priority" button that opens the same priority form as double-click on the mod list.
5. A "Clear Priority" button (only enabled when Priority > 0).

### Conflict graph tab — gui/conflict_view.go
Shown when the toolbar "Conflicts" button is pressed (replaces the right pane content, or opens in a `dialog.NewCustom` — your choice, prefer the dialog for simplicity).

- Title: "Conflict Graph — N conflicting files".
- A `widget.Table` with columns:

| # | Header | Width | Content |
|---|---|---|---|
| 0 | Resource | 420 | Path or `0x<hex>`, monospace, dim |
| 1 | Mods | 500 | All conflicting mod names joined by ✗, red |

- Scrollable, no row limit.
- A `widget.SearchEntry` above the table to filter by mod name or resource path (case-insensitive substring).

### Apply dialog — gui/apply_view.go
Opened by the toolbar "Apply" button as `dialog.NewCustomWithButtons`.

- Shows the full ordered list in a `widget.List` (scrollable):
  `1. modname.archive, 2. … etc.`
- Buttons: "Write modlist.txt" (primary) and "Cancel".
- On success: show `dialog.ShowInformation` — "Done! modlist.txt written. Previous file backed up."
- On error: show `dialog.ShowError`.

## Core logic rules — NEVER change these
The four `internal/` packages encode correct game behaviour. Do not alter:

1. **Load order semantics**: first entry in modlist.txt loads first and wins conflicts.
2. **`ApplyPriorities()` sort rules**:
   - Priority mods first, lower number = earlier load.
   - Tie-break: higher ConflictCount wins (maximises wins).
   - Unset (0) mods go last, sorted alphabetically.
3. **`modlist.Write()`** always backs up the existing file before overwriting.
4. **Archive parsing**: reads the RED4 index (FNV1a-64 hashes) starting at `indexOffset`; skip 24 bytes per entry after the hash.

## State management
All mutable state lives in `gui/app.go` as fields of an `AppState` struct:
```go
type AppState struct {
    cfg     *config.Config
    result  *conflict.Result   // nil until first scan
    modDir  string
    selected int               // currently selected row index (-1 = none)
}
```
Pass `*AppState` to every view constructor — do not use globals.
After any operation that changes `result.Mods` order (priority change, rescan), call `modListTable.Refresh()` and update the status bar label.

## Error handling
- Scan errors (bad directory, unreadable files): `dialog.ShowError(err, window)`.
- Non-fatal parse failures (one bad archive): log with `fyne.LogError` and continue — do not abort the scan.
- Config save errors: silent (best-effort persistence).

## What NOT to do
- Do not use `widget.List` for the mod table — use `widget.Table` for performance.
- Do not embed the TUI bubbletea code; the GUI replaces it entirely.
- Do not call `os.Exit` anywhere except `main.go` on fatal startup error.
- Do not add a splash screen, auto-updater, or telemetry.
- Do not use `layout.NewGridLayout` for the main split — use `container.NewHSplit`.
- Do not hard-code the mod directory path — always load from `config.Config.ModDir`.
