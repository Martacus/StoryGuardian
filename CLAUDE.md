# LitGuardian

Local-first desktop worldbuilding app for LitRPG and fantasy authors. Users fully own their files (JSON on disk), no cloud dependency.

## Tech Stack

- **Backend**: Go 1.25, Wails 3 (v3.0.0-alpha.74)
- **Frontend**: Vue 3 + TypeScript, Vite 5, Pinia
- **UI**: shadcn-vue (New York style, Tailwind CSS v4, lucide icons)
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

## Build & Codegen

- **Regen bindings** after any Go service change: `wails3 generate bindings -ts -d frontend/bindings`
- **Go build check**: `go build -v .` (ignore `build/ios` errors — pre-existing Wails platform stubs)
- **Vite build check**: `cd frontend && npx vite build`
- **Bindings import path** from `src/`: use relative `../../bindings/litguardian` (`@/` alias won't reach bindings outside `src/`)

## Project Structure

```
├── main.go                     # App entry point, window setup, service registration
├── appconfig.go                # AppConfig struct
├── appconfigservice.go         # App-level settings (recent worlds, etc.)
├── worldservice.go             # World CRUD — create, open, save, delete
├── layoutservice.go            # Per-world layout.json read/write
├── go.mod
├── Taskfile.yml
├── plans/                      # Implementation plan docs (gitignored)
├── frontend/
│   ├── bindings/               # Auto-generated Go↔TS bindings (do not edit)
│   │   └── litguardian/
│   │       ├── appconfigservice.ts
│   │       ├── worldservice.ts
│   │       ├── layoutservice.ts
│   │       └── models.ts
│   ├── src/
│   │   ├── main.ts             # Vue entry point
│   │   ├── App.vue             # Root component
│   │   ├── assets/index.css    # Tailwind + global styles
│   │   ├── components/
│   │   │   ├── WelcomeScreen.vue       # World picker / launch screen
│   │   │   ├── DashboardLayout.vue     # App shell with sidebar nav
│   │   │   ├── WorldOverview.vue       # Thin wrapper → <ModuleGrid view-id="overview">
│   │   │   ├── CreateWorldDialog.vue
│   │   │   ├── WorldListItem.vue
│   │   │   ├── SidebarNavItem.vue
│   │   │   ├── modules/
│   │   │   │   ├── ModuleGrid.vue      # 12-col grid, VueDraggable, Customize/Done button
│   │   │   │   └── ModuleCard.vue      # Card wrapper — drag handle, resize, visibility toggle
│   │   │   └── ui/                     # shadcn-vue components
│   │   ├── modules/            # Feature module components (lazy-loaded)
│   │   │   ├── registry.ts             # Central module registry (source of truth)
│   │   │   └── overview/
│   │   │       ├── InfoModule.vue      # World name/description form
│   │   │       ├── EntitiesModule.vue  # Placeholder
│   │   │       └── ImagesModule.vue    # Placeholder
│   │   ├── stores/
│   │   │   ├── worldStore.ts           # World open/close/save state
│   │   │   └── layoutStore.ts          # Module layout state, edit mode, reconciliation
│   │   ├── composables/
│   │   │   └── useAppToast.ts
│   │   └── lib/
│   │       ├── utils.ts
│   │       └── timeago.ts
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
└── build/                      # Platform-specific build configs
```

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
- **Frontend orchestrates**: Pinia stores coordinate cross-service workflows (e.g., create world + add to recents). Go services stay decoupled — no service-to-service calls.
- **Go services are independent**: Each service owns its own data and `sync.RWMutex`. No shared state between services.
- **Error handling**: Go services return descriptive errors → Wails rejects the promise → Pinia store catches and sets `error` ref → component displays to user. No silent failures.
- **Atomic writes**: All file mutations use temp file + `os.Rename` to prevent corruption on crash.
- **Cross-store calls in Pinia**: call `useOtherStore()` inside the action body, not at module top-level (avoids circular import issues).
- All file I/O goes through Go services, never from frontend
- Entity data stored as individual JSON files: `world-folder/entities/{id}.json`
- Use shadcn-vue components wherever possible; add new ones via `npx shadcn-vue@latest add <component>`

## Implemented Features

### Modular Dashboard System

A reusable, customizable grid of draggable module cards. Users can reorder, resize (12-col grid), and show/hide modules. Layout persists per-world in `layout.json`.

**Key files:**
- `layoutservice.go` — reads/writes `layout.json` atomically
- `frontend/src/modules/registry.ts` — central module registry; single source of truth for all modules
- `frontend/src/stores/layoutStore.ts` — layout state, edit mode, reconciliation (adds new modules, drops removed ones, preserves user settings)
- `frontend/src/components/modules/ModuleGrid.vue` — 12-col grid + VueDraggable
- `frontend/src/components/modules/ModuleCard.vue` — card wrapper with drag handle, resize dropdown, visibility toggle

**Adding a new module:**
1. Create `frontend/src/modules/<view>/<Name>Module.vue`
2. Add one entry to `moduleRegistry` in `registry.ts` with `views: ['<viewId>']`
3. Grid, card, store, and Go service handle everything automatically

**Adding a new view with its own dashboard:**
1. Create modules + registry entries with `views: ['newview']`
2. Create `NewView.vue` → `<ModuleGrid view-id="newview">`
3. Add to sidebar nav in `DashboardLayout.vue`
