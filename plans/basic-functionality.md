# Plan: Dashboard Screen with World Details

## Context
After opening a world, LitGuardian currently shows a placeholder ("World X is open" + Close button). This plan replaces that placeholder with the real dashboard — a modern layout with a sidebar, top nav, and a World Overview page showing editable world details.

This covers the first "inside a world" screen. Future plans will add entity CRUD, links, tags, etc. as additional views within this same dashboard shell.

### Current state
- `App.vue` conditionally renders `WelcomeScreen` vs a placeholder `<div>`
- `WorldMeta` has `name`, `createdAt`, `updatedAt` (no `description`)
- `WorldService` has `CreateWorld`, `OpenWorld`, `SelectFolder` — no way to read or update metadata after creation
- `worldStore` tracks `currentWorld` (a `WorldInfo` with name/path/lastOpened) but not the full `WorldMeta`

---

## Step 1 — Backend: Extend WorldMeta + WorldService

### 1.1 Add `Description` field to `WorldMeta`
```go
// worldservice.go
type WorldMeta struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```
Existing `world.json` files without a `description` key will unmarshal cleanly — Go zero-values the field to `""`.

### 1.2 Add new methods to `WorldService`

| Method | Purpose |
|--------|---------|
| `GetWorldMeta(folderPath string) (*WorldMeta, error)` | Reads `world.json` and returns the full metadata |
| `UpdateWorldMeta(folderPath, name, description string) (*WorldMeta, error)` | Updates name + description, bumps `UpdatedAt`, writes atomically, returns updated meta |

**`GetWorldMeta`** — straightforward read + unmarshal. Returns user-friendly error if `world.json` is missing.

**`UpdateWorldMeta`** — reads existing meta (to preserve `CreatedAt`), updates `Name`, `Description`, sets `UpdatedAt = now`, calls `writeJSONAtomic`. Returns the updated `WorldMeta` so the frontend can refresh without a second round-trip.

Validation:
- `name` must be non-empty
- `folderPath` must be non-empty and contain `world.json`

### 1.3 Regenerate bindings
Run `task dev` (or `wails3 generate bindings`) so new methods appear in `frontend/bindings/litguardian/worldservice.ts` and `WorldMeta` appears in `models.ts`.

---

## Step 2 — Frontend: Install additional shadcn-vue components

```bash
cd frontend && npx shadcn-vue@latest add textarea
```

Only `textarea` is needed — all other components (button, card, input, label, separator, tooltip, scroll-area, dialog, dropdown-menu) are already installed.

---

## Step 3 — Frontend: Update World Store

### 3.1 Add world meta state and actions to `worldStore.ts`

**New state:**
- `worldMeta: ref<WorldMeta | null>` — full metadata of the open world

**New actions:**
| Action | Purpose |
|--------|---------|
| `loadWorldMeta()` | Calls `WorldService.GetWorldMeta(currentWorld.path)`, sets `worldMeta` |
| `updateWorldMeta(name, description)` | Calls `WorldService.UpdateWorldMeta(...)`, updates `worldMeta` + `currentWorld.name` |

**Modified actions:**
- `createWorld()` and `openWorld()` — after setting `currentWorld`, also call `loadWorldMeta()`
- `closeWorld()` — also clear `worldMeta` to `null`

---

## Step 4 — Frontend: Dashboard Layout

### 4.1 `DashboardLayout.vue` — Shell component

```
┌──────────────────────────────────────────────────┐
│  🔍 Search                                      │
├────────┬─────────────────────────────────────────┤
│  Logo  │                                         │
│        │  Content Area                           │
│  ○ Ovw │  (renders the active view component)    │
│  ○ Ent │                                         │
│  ○ Lnk │                                         │
│  ○ Tag │                                         │
│  ○ Cat │                                         │
│        │                                         │
│  [←]   │                                         │
└────────┴─────────────────────────────────────────┘
```

**Top nav bar** (full width, above sidebar + content):
- Search `Input` with `Search` icon — placeholder only for now, no functionality
- Right side: world name displayed as muted text

**Sidebar** (left, fixed width ~200px):
- App logo + "LitGuardian" text at top
- Navigation items with lucide icons + labels:
  - `Globe` → Overview (active in this plan)
  - `Users` → Entities (disabled placeholder)
  - `Link` → Links (disabled placeholder)
  - `Tags` → Tags (disabled placeholder)
  - `FolderTree` → Categories (disabled placeholder)
- Bottom: "Close World" button → `worldStore.closeWorld()`
- Active item uses `bg-sidebar-accent text-sidebar-accent-foreground` from the existing CSS sidebar vars
- Sidebar background uses `bg-sidebar` / `text-sidebar-foreground`

**Content area** — renders the active view. For now only `WorldOverview`. Use a simple reactive `currentView` string as local state in `DashboardLayout.vue`:
```ts
const currentView = ref<'overview' | 'entities' | 'links' | 'tags' | 'categories'>('overview')
```

### 4.2 `SidebarNavItem.vue` — Reusable nav button
- Props: `icon`, `label`, `active`, `disabled`
- Styling: `hover:bg-sidebar-accent`, active state highlight, disabled state muted
- Click emits `select`

---

## Step 5 — Frontend: World Overview Page

### 5.1 `WorldOverview.vue`

Displays and allows editing of world name and description.

**Layout:**
```
┌─────────────────────────────────────────┐
│  World Overview              [Save]     │
├─────────────────────────────────────────┤
│                                         │
│  Name                                   │
│  ┌───────────────────────────────────┐  │
│  │ My Fantasy World                  │  │
│  └───────────────────────────────────┘  │
│                                         │
│  Description                            │
│  ┌───────────────────────────────────┐  │
│  │ A sprawling fantasy world with    │  │
│  │ magic systems and warring         │  │
│  │ kingdoms...                       │  │
│  │                                   │  │
│  └───────────────────────────────────┘  │
│                                         │
└─────────────────────────────────────────┘
```

**Behavior:**
- Loads `worldStore.worldMeta` on mount (should already be populated by `openWorld`/`createWorld`)
- Name `Input` and description `Textarea` are pre-filled from `worldMeta`
- "Save" button enabled only when values differ from stored meta (dirty tracking)
- On save: calls `worldStore.updateWorldMeta(name, description)`
- Error displayed inline if save fails
- Created date shown formatted, last modified shown as relative time via `timeAgo()`

---

## Step 6 — Wire Up App.vue

### 6.1 Replace placeholder in `App.vue`
```vue
<WelcomeScreen v-if="!worldStore.hasOpenWorld" />
<DashboardLayout v-else />
```

Remove the existing placeholder `<div v-else>` and the `Button` import (no longer needed in App.vue).

---

## File Summary

**Modified Go files:**
- `worldservice.go` — Add `Description` to `WorldMeta`, add `GetWorldMeta`, `UpdateWorldMeta`

**New frontend files:**
- `src/components/DashboardLayout.vue` — Shell: sidebar + top nav + content area
- `src/components/SidebarNavItem.vue` — Reusable sidebar nav button
- `src/components/WorldOverview.vue` — World details view (name, description, dates)

**Modified frontend files:**
- `src/stores/worldStore.ts` — Add `worldMeta`, `loadWorldMeta`, `updateWorldMeta`; update `createWorld`/`openWorld`/`closeWorld`
- `src/App.vue` — Render `DashboardLayout` instead of placeholder

**Auto-generated (by Wails, do not edit):**
- `frontend/bindings/litguardian/worldservice.ts` — new method bindings
- `frontend/bindings/litguardian/models.ts` — `WorldMeta` type added

---

## Verification

1. `task dev` — open a world → dashboard layout appears with sidebar + top nav
2. Sidebar shows "Overview" as active, other items show as disabled placeholders
3. World Overview displays name and description fields pre-filled from `world.json`
4. Edit name or description → "Save" button enables → click Save → `world.json` on disk is updated, `UpdatedAt` bumps
5. Created date and last modified time display correctly
6. Search bar is visible but non-functional (placeholder for future)
7. Click "Close World" in sidebar → returns to WelcomeScreen
8. Re-open the world → saved name/description persist
9. World name change reflects in the top nav and in the recents list after closing
