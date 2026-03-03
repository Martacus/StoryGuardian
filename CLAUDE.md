# StoryGuardian

Local-first desktop worldbuilding app for LitRPG and fantasy authors. Users fully own their files (JSON on disk), no cloud dependency.

## Tech Stack

- **Backend**: Go 1.25, Wails 3 (v3.0.0-alpha.74)
- **Frontend**: Vue 3 + TypeScript, Vite 5, Pinia (planned)
- **Runtime**: Wails bindings bridge Go services ↔ TypeScript
- **Storage**: JSON files in user-chosen world folders
- **Build**: Task runner (go-task), Taskfile.yml

## Architecture

```
Vue 3 components → Pinia stores → Wails-generated bindings → Go services → JSON files on disk
```

- Go services handle ALL filesystem I/O — frontend never touches disk directly
- Atomic writes: temp file + rename
- Wails auto-generates TypeScript bindings from Go service structs into `frontend/bindings/`

## Commands

| Command | Description |
|---------|-------------|
| `task dev` | Dev mode with hot-reload |
| `task build` | Build for current platform |
| `task run` | Run the built application |
| `task package` | Create installer (NSIS on Windows) |

## Project Structure

```
├── main.go                  # App entry point, window setup, embedded assets
├── greetservice.go          # Example Go service (template)
├── go.mod                   # Go module (currently named "changeme")
├── Taskfile.yml             # Build tasks
├── frontend/
│   ├── src/                 # Vue 3 app source
│   │   ├── main.ts          # Vue entry point
│   │   ├── App.vue          # Root component
│   │   └── components/      # Vue components
│   ├── bindings/            # Auto-generated Go↔TS bindings (do not edit)
│   ├── index.html           # HTML shell
│   ├── vite.config.ts       # Vite config + Wails plugin
│   └── package.json         # Frontend deps
└── build/                   # Platform-specific build configs
```

## Known Issues

- Go module name is still `changeme` (needs renaming)
- App name has typo: `litguradian` instead of `litguardian` (in Taskfile, main.go, build configs)

## Design Docs (Obsidian Vault)

All design documentation lives in **`C:\Vaults\LitGuardian`**:

| Path | Topic |
|------|-------|
| `01-Product/Vision.md` | Project vision, principles, target users |
| `01-Product/MVP-Scope.md` | MVP v0.1 feature scope, acceptance criteria |
| `02-Architecture/System-Overview.md` | Runtime layers, key constraints |
| `02-Architecture/Data-Model.md` | Entity/Link schemas, world folder structure |
| `02-Architecture/Storage-Strategy.md` | Write policy, atomic ops, error handling |
| `05-Decisions/ADR-0001-Use-Wails.md` | Why Wails over Electron/Tauri |

## Conventions

- One Go service per domain concern, registered in `main.go`
- Frontend state via Vue composables (Pinia stores for shared state)
- All file I/O goes through Go services, never from frontend
- Entity data stored as individual JSON files: `world-folder/entities/{id}.json`
