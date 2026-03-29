# CP77 Mod Order — GUI

A Cyberpunk 2077 `.archive` mod conflict manager with a graphical interface.
Scans your mods folder, detects file conflicts between archives, and lets you control the load order by writing `modlist.txt`.

## Features

- Scans all `.archive` files and detects resource conflicts
- Drag-to-reorder mods freely in the main list
- Group related (conflicting) mods into a contiguous block automatically
- Drag-to-reorder within a conflict group from the detail panel
- Apply writes `modlist.txt`; previous file is backed up to `modlist.old/`
- Restore any backup from `modlist.old/` in one click
- Conflict resolution wizard for newly detected mods on rescan

## Requirements

- Windows 10 or 11 (x64)
- [Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) — already installed on most Windows 10/11 systems via Windows Update

## Prerequisites (building from source)

| Tool | Version | Notes |
|------|---------|-------|
| [Go](https://go.dev/dl/) | 1.22+ | |
| [Wails CLI](https://wails.io/docs/gettingstarted/installation) | v2 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| [Node.js](https://nodejs.org/) | 18+ | npm included |
| GCC | any recent | Windows: install [MSYS2](https://www.msys2.org/) and run `pacman -S mingw-w64-ucrt-x86_64-gcc` |

> **Windows path**: the Makefile expects GCC at `C:/msys64/ucrt64/bin/gcc.exe`.
> Adjust `CC` in the Makefile if your MSYS2 is installed elsewhere.

## Install dependencies

```bash
make tidy
# equivalent to:
go mod tidy
cd frontend && npm install
```

## Development

Runs the Go backend and Vite dev server with hot-reload:

```bash
make dev
# equivalent to:
wails dev
```

The app window opens automatically. Frontend changes are reflected instantly; Go changes trigger a backend restart.

## Production build

```bash
make build
# equivalent to:
wails build -platform windows/amd64 -o CP77-modorder.exe
```

Output: `CP77-modorder.exe` in the project root.
The `-H windowsgui` linker flag (set in `wails.json`) suppresses the console window.

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Production build → `CP77-modorder.exe` |
| `make dev` | Development mode with hot-reload |
| `make generate` | Regenerate Wails JS bindings (`wails generate module`) |
| `make tidy` | `go mod tidy` + `npm install` |
| `make clean` | Remove the binary and `frontend/dist/` |

## Usage

1. Launch `CP77-modorder.exe`
2. Click **Open** and select your mods folder — typically:
   `C:\Program Files (x86)\Steam\steamapps\common\Cyberpunk 2077\archive\pc\mod`
3. The app scans all `.archive` files and shows conflicts in the table
4. Reorder mods by dragging rows, or click a mod to use the detail panel
5. Click **Group** to automatically consolidate conflicting mods together
6. Click **Apply** to write `modlist.txt` — the game reads this on startup

## Project layout

```
cp77-modorder-gui/
├── main.go                  # Wails entry point
├── app.go                   # Exported backend methods (bound to JS)
├── dto.go                   # Data transfer types
├── wails.json               # Wails configuration
├── Makefile
├── internal/
│   ├── archive/             # RED4/RDAR binary parser
│   ├── conflict/            # Conflict detection + priority sort
│   ├── config/              # JSON config persistence
│   └── modlist/             # modlist.txt writer
└── frontend/
    └── src/
        ├── App.vue
        ├── components/      # Toolbar, ModListTable, DetailPanel, dialogs…
        ├── composables/     # useWails — IPC bridge
        └── stores/          # Pinia app state
```
