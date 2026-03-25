# Conventions

## Go

- One service per domain concern, registered in `main.go`
- Each service owns its own `sync.RWMutex` — no shared state
- No service-to-service calls
- Return descriptive errors — no silent failures
- All file writes use temp file + `os.Rename` (atomic)

## Frontend / Pinia

- **Cross-store calls**: call `useOtherStore()` inside the action body, not at module top-level (avoids circular import issues)
- Pinia stores are the orchestration layer — they coordinate multi-service workflows
- Components never call Go bindings directly; always go through a store

## UI

- Use shadcn-vue components wherever possible
- shadcn-vue style: New York, Tailwind CSS v4, lucide icons
- Add new components via `pnpm dlx shadcn-vue@latest add <component>`

## Design Docs (Obsidian Vault)

All design docs live in `C:\Vaults\LitGuardian`:

| Path | Topic |
|------|-------|
| `01-Product/Vision.md` | Project vision, principles, target users |
| `01-Product/MVP-Scope.md` | MVP v0.1 feature scope, acceptance criteria |
| `02-Architecture/System-Overview.md` | Runtime layers, key constraints |
| `02-Architecture/Data-Model.md` | Entity/Link schemas, world folder structure |
| `02-Architecture/Storage-Strategy.md` | Write policy, atomic ops, error handling |
| `05-Decisions/ADR-0001-Use-Wails.md` | Why Wails over Electron/Tauri |

Plans live in `C:\Projects\StoryGuardian\plans\` — read the relevant plan file before implementing.
