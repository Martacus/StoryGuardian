# Plan: Welcome/Launcher Screen (Obsidian-style World Selector)

## Context
LitGuardian needs a launcher screen — the first thing a user sees when opening the app. Like Obsidian's vault picker: left panel lists recently opened worlds, right panel has "Create New World" and "Open Existing Folder" buttons. This is explicitly in the MVP scope ("World selector: create new, open existing, view recent worlds").

Currently the app shows a template demo (HelloWorld + GreetService). We need to replace it with the real launcher.

---

## Step 0 — Housekeeping

### 0.1 Fix app name typo
`litguradian` → `litguardian` in:
- `Taskfile.yml` line 12 (`APP_NAME`)
- `main.go` line 38 (`Name`)
- `build/windows/` files if needed

### 0.2 Update CLAUDE.md
Add shadcn-vue to the Tech Stack and Conventions sections so future sessions know the UI library.

---

## Step 1 — Go Backend: AppConfigService

Persists app-level config (recent worlds list) to the OS config directory.

### 1.1 Create data types (`appconfig.go`)
```go
type WorldInfo struct {
    Name       string    `json:"name"`
    Path       string    `json:"path"`
    LastOpened time.Time `json:"lastOpened"`
}

type AppConfig struct {
    RecentWorlds []WorldInfo `json:"recentWorlds"`
}
```

### 1.2 Create service (`appconfigservice.go`)
- Uses `os.UserConfigDir()` → `%APPDATA%/LitGuardian/config.json` on Windows
- Implements `ServiceStartup(ctx, options)` — loads config from disk on app start
- Atomic writes (temp file + rename) per Storage Strategy doc

**Exposed methods (become Wails bindings):**
| Method | Purpose |
|--------|---------|
| `GetRecentWorlds() []WorldInfo` | Returns list sorted by LastOpened desc |
| `AddRecentWorld(name, path string) error` | Upserts by path, caps at 10, persists |
| `RemoveRecentWorld(path string) error` | Removes entry, persists |

### 1.3 Register in `main.go`
Add `application.NewService(&AppConfigService{})` to the Services slice.

---

## Step 2 — Go Backend: WorldService

Handles world folder operations and native OS folder picker.

### 2.1 Create service (`worldservice.go`)
```go
type WorldMeta struct {
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

**Exposed methods:**
| Method | Purpose |
|--------|---------|
| `SelectFolder() (string, error)` | Native OS directory picker via Wails dialog API |
| `CreateWorld(name, folderPath string) (*WorldInfo, error)` | Creates folder structure: `world.json`, `entities/`, `links.json`, `tags.json`, `categories.json` |
| `OpenWorld(folderPath string) (*WorldInfo, error)` | Validates folder has `world.json`, reads name, returns WorldInfo |

**Wails 3 dialog API usage:**
```go
dialog := application.Get().Dialog.OpenFile()
dialog.CanChooseDirectories(true)
dialog.CanChooseFiles(false)
dialog.CanCreateDirectories(true)
dialog.SetTitle("Select World Folder")
result, err := dialog.PromptForSingleSelection()
```

### 2.2 Register in `main.go`
Add `application.NewService(&WorldService{})` to the Services slice.

### 2.3 Clean up template code in `main.go`
- Remove `init()` function (event registration)
- Remove time-emitting goroutine (lines 69–75)
- Update window title, app name, description

---

## Step 3 — Frontend: Pinia + shadcn-vue Components

### 3.1 Install Pinia
```bash
cd frontend && npm install pinia
```
Update `main.ts` to initialize Pinia.

### 3.2 Install shadcn-vue components
```bash
npx shadcn-vue@latest add card scroll-area dialog input label dropdown-menu separator tooltip
```

---

## Step 4 — Frontend: World Store

### 4.1 Create Pinia store (`frontend/src/stores/worldStore.ts`)
- `currentWorld: ref<WorldInfo | null>` — null means show welcome screen
- `recentWorlds: ref<WorldInfo[]>`
- `hasOpenWorld: computed` — drives App.vue conditional rendering
- Actions call the Wails-generated bindings (`AppConfigService.*`, `WorldService.*`)
- `createWorld()` orchestrates: WorldService.CreateWorld → AppConfigService.AddRecentWorld → set currentWorld
- `openWorld()` orchestrates: WorldService.OpenWorld → AppConfigService.AddRecentWorld → set currentWorld

---

## Step 5 — Frontend: Welcome Screen UI

### 5.1 `WelcomeScreen.vue` — Main layout
```
┌───────────────────────────────────────────────────┐
│                  LitGuardian                     │
├─────────────────────────┬─────────────────────────┤
│  Recent Worlds          │  Get Started            │
│  (ScrollArea)           │                         │
│                         │  [✚ Create New World]   │
│  ▸ My Fantasy World     │                         │
│    C:\Worlds\fantasy    │  [📂 Open Folder]       │
│                         │                         │
│  ▸ LitRPG Draft         │                         │
│    C:\Worlds\litrpg     │                         │
│                         │                         │
│  (empty state if none)  │                         │
├─────────────────────────┴─────────────────────────┤
└───────────────────────────────────────────────────┘
```
- Full viewport (`h-screen w-screen`), two-column grid
- Uses: `ScrollArea`, `Button`, `Separator`
- Icons from `lucide-vue-next` (already a dep)
- Calls `worldStore.loadRecentWorlds()` on mount

### 5.2 `WorldListItem.vue` — Each entry in the recents list
- Props: `world: WorldInfo`
- Shows name (prominent) + path (muted, truncated with `Tooltip`)
- Click → `worldStore.openWorld(world.path)`
- `DropdownMenu` (⋯ button) → "Remove from recents"
- Hover styling: `hover:bg-accent rounded-lg p-3`

### 5.3 `CreateWorldDialog.vue` — Modal for new world creation
- Uses shadcn `Dialog`, `Input`, `Label`, `Button`
- Flow: enter name → click "Choose Folder" (calls `WorldService.SelectFolder()`) → "Create"
- Validates: name non-empty, folder selected
- On success: closes dialog, world opens

### 5.4 `lib/timeago.ts` — Small utility
Pure function converting ISO timestamps to "2 hours ago" strings. No external dependency needed.

### 5.5 Update `App.vue` — Conditional rendering
```vue
<WelcomeScreen v-if="!worldStore.hasOpenWorld" />
<div v-else>
  <!-- Future: MainEditor -->
  <p>World "{{ worldStore.currentWorld?.name }}" is open.</p>
  <Button @click="worldStore.closeWorld()">Close</Button>
</div>
```

### 5.6 Clean up
- Delete `HelloWorld.vue`
- Update `index.html` title to "LitGuardian"

---

## File Summary

**New Go files:**
- `appconfig.go` — WorldInfo, AppConfig structs
- `appconfigservice.go` — Config persistence service
- `worldservice.go` — World folder ops + dialog

**New frontend files:**
- `src/stores/worldStore.ts` — Pinia store
- `src/components/WelcomeScreen.vue` — Launcher layout
- `src/components/WorldListItem.vue` — World list entry
- `src/components/CreateWorldDialog.vue` — Create world modal
- `src/lib/timeago.ts` — Relative time helper

**Modified:**
- `main.go` — Register services, remove template code, fix name
- `Taskfile.yml` — Fix APP_NAME
- `CLAUDE.md` — Add shadcn-vue to tech stack
- `frontend/src/main.ts` — Add Pinia
- `frontend/src/App.vue` — Replace template with WelcomeScreen routing

**Deleted:**
- `frontend/src/components/HelloWorld.vue`

**Auto-generated (by Wails, do not edit):**
- `frontend/bindings/changeme/appconfigservice.ts`
- `frontend/bindings/changeme/worldservice.ts`

---

## Verification
1. `task dev` — app launches, shows WelcomeScreen (empty recents)
2. Click "Create New World" → enter name, pick folder → world folder created on disk with `world.json`, `entities/`, etc.
3. App shows placeholder "World is open" view
4. Close world → back to WelcomeScreen, world appears in recents
5. Click recent world → opens it again
6. Dropdown → "Remove from recents" → entry disappears
7. "Open Existing Folder" → native picker → select folder with `world.json` → opens
8. Close + reopen app → recents persist from `%APPDATA%/LitGuardian/config.json`
9. Open a folder without `world.json` → user-friendly error
