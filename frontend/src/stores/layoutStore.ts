import { ref } from 'vue'
import { defineStore } from 'pinia'
import { LayoutService } from '../../bindings/litguardian'
import { DashboardLayout, ViewLayout, ModuleLayout } from '../../bindings/litguardian/models'
import { moduleRegistry, getDefaultModules, type ModuleDefinition } from '@/modules/registry'
import { useWorldStore } from './worldStore'
import { useAppToast } from '@/composables/useAppToast'

/** Minimal shape that grid-layout-plus provides in the layout-updated event. */
export interface GridItem {
  i: string
  x: number
  y: number
  w: number
  h: number
}

export const useLayoutStore = defineStore('layout', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  const layout = ref<DashboardLayout>(new DashboardLayout({ version: 2, views: {} }))
  const editMode = ref(false)

  // ── Actions ────────────────────────────────────────────────────────────────

  /** Load layout.json from disk for the open world. Falls back silently to empty state. */
  async function loadLayout() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      const saved = await LayoutService.GetLayout(worldStore.currentWorld.path)
      layout.value = saved ?? new DashboardLayout({ version: 2, views: {} })
    } catch (e) {
      errorToast('Failed to load layout', String(e))
      layout.value = new DashboardLayout({ version: 2, views: {} })
    }
  }

  /**
   * Returns the reconciled module list for a view.
   * - Saved modules array is the source of truth — only what's saved is shown
   * - Falls back to registry defaults when no layout exists for the view
   * - Drops modules that have been removed from the registry (stale entries)
   * - Does NOT auto-append new registry modules; use addModule() explicitly
   */
  function getViewModules(viewId: string): ModuleLayout[] {
    const savedView = layout.value.views?.[viewId]
    const registryForView = moduleRegistry.filter(m => m.views.includes(viewId))

    if (!savedView || savedView.modules.length === 0) {
      return getDefaultModules(viewId)
    }

    const registryIds = new Set(registryForView.map(m => m.id))
    // Keep saved modules that still exist in registry (preserves x/y/w/h)
    return savedView.modules.filter(m => registryIds.has(m.id))
  }

  /** Replace the full module list for a view. */
  function updateViewLayout(viewId: string, modules: ModuleLayout[]) {
    if (!layout.value.views) layout.value.views = {}
    layout.value.views[viewId] = new ViewLayout({ modules })
  }

  /**
   * Called by ModuleGrid when grid-layout-plus emits layout-updated after a
   * drag or resize. Merges the new x/y/w/h back into the store.
   */
  function handleLayoutUpdate(viewId: string, items: GridItem[]) {
    const current = getViewModules(viewId)
    const modules = current.map(m => {
      const updated = items.find(u => u.i === m.id)
      if (updated) {
        return new ModuleLayout({ ...m, x: updated.x, y: updated.y, w: updated.w, h: updated.h })
      }
      return m
    })
    updateViewLayout(viewId, modules)
  }

  /** Toggle visibility for a specific module in a view. */
  function toggleModuleVisibility(viewId: string, moduleId: string) {
    const modules = getViewModules(viewId).map(m =>
      m.id === moduleId ? new ModuleLayout({ ...m, visible: !m.visible }) : m
    )
    updateViewLayout(viewId, modules)
  }

  /** Remove a module from the grid entirely. It can be re-added via addModule(). */
  async function removeModule(viewId: string, moduleId: string) {
    const modules = getViewModules(viewId).filter(m => m.id !== moduleId)
    updateViewLayout(viewId, modules)
    await saveLayout()
  }

  /**
   * Add a module to the grid at a default position below existing modules.
   * If the module is already on the grid, this is a no-op.
   */
  async function addModule(viewId: string, moduleId: string) {
    const current = getViewModules(viewId)
    if (current.some(m => m.id === moduleId)) return

    const def = moduleRegistry.find(m => m.id === moduleId)
    if (!def) return

    const maxY = current.reduce((max, m) => Math.max(max, m.y + m.h), 0)

    const modules = [
      ...current,
      new ModuleLayout({
        id: moduleId,
        x: 0,
        y: maxY,
        w: def.defaultW,
        h: def.defaultH,
        visible: true,
      }),
    ]
    updateViewLayout(viewId, modules)
    await saveLayout()
  }

  /**
   * Returns registry modules for a view that are not currently on the grid.
   * These are the candidates for the "Add Module" dropdown.
   */
  function getAddableModules(viewId: string): ModuleDefinition[] {
    const currentIds = new Set(getViewModules(viewId).map(m => m.id))
    return moduleRegistry.filter(m => m.views.includes(viewId) && !currentIds.has(m.id))
  }

  /** Set a single config key for a module. Merges into existing config. */
  function setModuleConfig(viewId: string, moduleId: string, key: string, value: unknown) {
    const modules = getViewModules(viewId).map(m => {
      if (m.id !== moduleId) return m
      const config = { ...(m.config ?? {}), [key]: value }
      return new ModuleLayout({ ...m, config })
    })
    updateViewLayout(viewId, modules)
  }

  /** Read a config value for a module, with a typed default fallback. */
  function getModuleConfig<T>(viewId: string, moduleId: string, key: string, fallback: T): T {
    const view = layout.value.views?.[viewId]
    const mod = view?.modules?.find(m => m.id === moduleId)
    const val = mod?.config?.[key]
    return (val as T) ?? fallback
  }

  /** Toggle edit mode. Auto-saves layout.json when exiting edit mode. */
  async function toggleEditMode() {
    if (editMode.value) {
      await saveLayout()
    }
    editMode.value = !editMode.value
  }

  /** Persist current layout to disk. */
  async function saveLayout() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      await LayoutService.SaveLayout(worldStore.currentWorld.path, layout.value)
    } catch (e) {
      errorToast('Failed to save layout', String(e))
    }
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    layout.value = new DashboardLayout({ version: 2, views: {} })
    editMode.value = false
  }

  return {
    layout,
    editMode,
    loadLayout,
    getViewModules,
    updateViewLayout,
    handleLayoutUpdate,
    toggleModuleVisibility,
    removeModule,
    addModule,
    getAddableModules,
    setModuleConfig,
    getModuleConfig,
    toggleEditMode,
    saveLayout,
    reset,
  }
})
