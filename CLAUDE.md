# LitGuardian

Local-first desktop worldbuilding app for LitRPG and fantasy authors. Users fully own their files (JSON on disk), no cloud dependency.

## Tech Stack

- **Backend**: Go 1.25, Wails 3 (v3.0.0-alpha.74)
- **Frontend**: Vue 3 + TypeScript, Vite 5, Pinia
- **UI**: shadcn-vue (New York style, Tailwind CSS v4, lucide icons)
- **Storage**: JSON files in user-chosen world folders
- **Build**: Task runner (go-task), Taskfile.yml

## Project Structure

```
├── main.go                     # App entry point, window setup, service registration
├── appconfig.go                # AppConfig struct
├── appconfigservice.go         # App-level settings (recent worlds, etc.)
├── worldservice.go             # World CRUD — create, open, save, delete
├── layoutservice.go            # Per-world layout.json read/write
├── go.mod / Taskfile.yml
├── plans/                      # Implementation plan docs (gitignored)
├── frontend/
│   ├── bindings/               # Auto-generated Go↔TS bindings (do not edit)
│   │   └── litguardian/
│   ├── src/
│   │   ├── main.ts / App.vue
│   │   ├── components/
│   │   │   ├── WelcomeScreen.vue
│   │   │   ├── DashboardLayout.vue
│   │   │   ├── WorldOverview.vue
│   │   │   └── modules/
│   │   │       ├── ModuleGrid.vue
│   │   │       └── ModuleCard.vue
│   │   ├── modules/
│   │   │   ├── registry.ts
│   │   │   └── overview/
│   │   ├── stores/
│   │   │   ├── worldStore.ts
│   │   │   └── layoutStore.ts
│   │   └── composables / lib
└── build/
```

## Rules

See `.claude/rules/` for detailed guidance:

- `architecture.md` — data flow, service design, error handling, bindings
- `build.md` — commands, codegen, package manager
- `conventions.md` — Go, frontend, UI conventions; design docs vault map
- `modules.md` — modular dashboard system (grid-layout-plus, add/extend modules)
