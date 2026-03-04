# Plan: Modular Dashboard System

## Context
The user wants a customizable, modular UI where content is displayed as independent card modules. Users can reorder modules via drag-and-drop, resize them on a 12-column grid, and hide/show them. Layout preferences persist per-world in `layout.json`. This system will be the foundation for all views in the app going forward — starting with the World Overview.

Currently `WorldOverview.vue` is a monolithic form. This plan replaces it with a grid of draggable module cards, and creates the reusable infrastructure so future views (entities, links, etc.) can adopt the same pattern by simply registering modules.

---

## Step 1 — Go Backend: LayoutService

**New file: `layoutservice.go`**

```go
type ModuleLayout struct {
    ID      string `json:"id"`
    ColSpan int    `json:"colSpan"`
    Visible bool   `json:"visible"`
}
type ViewLayout struct { Modules []ModuleLayout `json:"modules"` }
type DashboardLayout struct { Version int `json:"version"`; Views map[string]ViewLayout `json:"views"` }
```

| Method | Purpose |
|--------|---------|
| `GetLayout(folderPath) (*DashboardLayout, error)` | Reads `layout.json`, returns `nil` if file doesn't exist (frontend uses defaults) |
| `SaveLayout(folderPath, layout) error` | Writes `layout.json` atomically via `writeJSONAtomic` |

Register `LayoutService` in `main.go`. Regenerate bindings with `wails3 generate bindings -ts`.

---

## Step 2 — Install vue-draggable-plus

```bash
cd frontend && npm install vue-draggable-plus
```

Vue 3 drag-and-drop wrapping SortableJS. Provides `<VueDraggable>` component with `v-model` array binding.

---

## Step 3 — Module Registry

**New file: `frontend/src/modules/registry.ts`**

Central definition of all modules:
```ts
interface ModuleDefinition {
  id: string; label: string; icon: Component
  defaultColSpan: number; component: Component; views: string[]
}
```

Initial overview modules:
- `overview-info` (name/description form, default 6 cols)
- `overview-entities` (placeholder, default 6 cols)
- `overview-images` (placeholder, default 12 cols)

Also exports `getDefaultModules(viewId)` — returns default layout array for worlds with no saved layout.

---

## Step 4 — Layout Store

**New file: `frontend/src/stores/layoutStore.ts`**

| State/Action | Purpose |
|---|---|
| `layout: ref<DashboardLayout>` | Full layout from disk (or defaults) |
| `editMode: ref<boolean>` | Toggles customize mode |
| `loadLayout()` | Reads `layout.json` via Go binding |
| `getViewModules(viewId)` | Returns modules array, falls back to registry defaults, reconciles new/removed modules |
| `updateModuleOrder(viewId, modules)` | Sets new order after drag |
| `setModuleColSpan(viewId, moduleId, span)` | Updates a module's column width |
| `toggleModuleVisibility(viewId, moduleId)` | Show/hide a module |
| `toggleEditMode()` | Toggle edit mode; auto-saves layout on exit |

**Wire into worldStore.ts:**
- `openWorld()` / `createWorld()` → call `layoutStore.loadLayout()` after `loadWorldMeta()`
- `closeWorld()` → reset layout store state

---

## Step 5 — ModuleCard.vue

**New file: `frontend/src/components/modules/ModuleCard.vue`**

Wraps each module in a shadcn `Card`. Two modes:

**View mode:** Clean card with title + content only.

**Edit mode:** Shows drag handle (`GripVertical`), resize dropdown (Quarter/Third/Half/Two-thirds/Full), and visibility toggle (`Eye`/`EyeOff`). Hidden modules show at reduced opacity with card shell visible so users can re-enable them.

Uses existing shadcn components: `Card`, `CardHeader`, `CardTitle`, `CardContent`, `CardAction`, `DropdownMenu`.

---

## Step 6 — ModuleGrid.vue

**New file: `frontend/src/components/modules/ModuleGrid.vue`**

Props: `viewId: string`. Slot: `#title` for the page heading.

- Renders a `grid-cols-12` CSS Grid
- `VueDraggable` wraps the grid items (drag disabled when not in edit mode, uses `.drag-handle` selector)
- Each module's `gridColumn: span N` sets its width
- "Customize" / "Done" button toggles edit mode
- In view mode: only visible modules shown. In edit mode: all modules shown (hidden ones dimmed)

---

## Step 7 — Module Content Components

Three components in `frontend/src/modules/overview/`:

**`InfoModule.vue`** — Extracted from current `WorldOverview.vue`. Contains name `Input`, description `Textarea`, dirty tracking, save button, and created/modified metadata. Uses `worldStore` directly.

**`EntitiesModule.vue`** — Placeholder stub ("No entities yet").

**`ImagesModule.vue`** — Placeholder stub ("No images yet").

---

## Step 8 — Replace WorldOverview.vue

Replace `WorldOverview.vue` contents with a thin wrapper:
```vue
<ModuleGrid view-id="overview">
  <template #title>World Overview</template>
</ModuleGrid>
```

`DashboardLayout.vue` and `App.vue` need **no changes** — they still render `<WorldOverview>`, which now internally uses ModuleGrid.

---

## File Summary

| File | Action |
|------|--------|
| `layoutservice.go` | New — Go layout persistence service |
| `main.go` | Edit — register LayoutService |
| `frontend/src/stores/layoutStore.ts` | New — layout Pinia store |
| `frontend/src/stores/worldStore.ts` | Edit — wire layout loading |
| `frontend/src/modules/registry.ts` | New — module definitions + defaults |
| `frontend/src/modules/overview/InfoModule.vue` | New — name/description form |
| `frontend/src/modules/overview/EntitiesModule.vue` | New — entities placeholder |
| `frontend/src/modules/overview/ImagesModule.vue` | New — images placeholder |
| `frontend/src/components/modules/ModuleCard.vue` | New — card wrapper with edit controls |
| `frontend/src/components/modules/ModuleGrid.vue` | New — 12-col grid + drag-and-drop |
| `frontend/src/components/WorldOverview.vue` | Replace — thin ModuleGrid wrapper |

---

## Future Extensibility

Adding a new module to any view:
1. Create `frontend/src/modules/<view>/<Name>Module.vue`
2. Add one entry to `moduleRegistry` in `registry.ts`
3. Done — the grid, card, store, and Go service handle everything automatically

Adding a new view with its own module dashboard:
1. Create modules + registry entries with `views: ['newview']`
2. Create `NewView.vue` → `<ModuleGrid view-id="newview">`
3. Add to sidebar nav in `DashboardLayout.vue`

---

## Verification

1. `task dev` → open a world → Overview shows 3 module cards in a grid
2. Click "Customize" → drag handles, resize dropdowns, visibility toggles appear
3. Drag a card → order changes live
4. Resize via dropdown → card width changes on the 12-col grid
5. Hide a module → it dims in edit mode, disappears in view mode
6. Click "Done" → edit controls disappear, `layout.json` written to world folder
7. Close and reopen the world → layout persists
8. Create a new world (no `layout.json`) → defaults from registry apply
9. Info module's name/description form still saves to `world.json` correctly
