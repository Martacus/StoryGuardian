# Modular Dashboard System

Free-form 2D grid of module cards. Users drag to any position, resize width/height, show/hide modules. Layout persists per-world in `layout.json` (v2 format).

**Library:** `grid-layout-plus`

## Key Files

| File | Role |
|------|------|
| `layoutservice.go` | Reads/writes `layout.json` atomically. Handles v1→v2 migration. Layout version = 2. |
| `frontend/bindings/litguardian/layoutservice.ts` | Auto-generated — `GetLayout`, `SaveLayout` |
| `frontend/bindings/litguardian/models.ts` | Auto-generated — `ModuleLayout` has `id, x, y, w, h, visible` |
| `frontend/src/modules/registry.ts` | Module registry — `defaultW`, `defaultH`, optional `minW`/`minH`. `getDefaultModules()` auto-positions x/y. |
| `frontend/src/stores/layoutStore.ts` | `handleLayoutUpdate()`, `updateViewLayout()`, `toggleModuleVisibility()` |
| `frontend/src/components/modules/ModuleCard.vue` | Drag handle (`.drag-handle`), visibility toggle (`.no-drag`), `h-full overflow-hidden` |
| `frontend/src/components/modules/ModuleGrid.vue` | `GridLayout` + `GridItem`. Local `gridLayout` ref synced from store. `layout-updated` → `handleLayoutUpdate`. |

## Behaviour

- **No vertical compaction** — modules stay exactly where placed; gaps allowed
- **Edit mode**: all modules shown (hidden ones at 40% opacity), `is-draggable` + `is-resizable` = true
- **View mode**: only visible modules rendered, dragging/resizing disabled
- Drag restricted to `.drag-handle` (GripVertical icon) via `drag-allow-from`
- Visibility toggle has `.no-drag` class to prevent accidental drags
- `getViewModules()` reconciles: keeps saved positions, appends new registry modules below existing, drops removed modules

## Adding a New Module

1. Create `frontend/src/modules/<view>/<Name>Module.vue`
2. Add one entry to `moduleRegistry` in `registry.ts` with `defaultW`, `defaultH`, and `views: ['<viewId>']`
3. Grid, card, store, and Go service handle everything automatically

## Adding a New View / Dashboard

1. Create modules + registry entries with `views: ['newview']`
2. Create `NewView.vue` → `<ModuleGrid view-id="newview">`
3. Add to sidebar nav in `DashboardLayout.vue`
