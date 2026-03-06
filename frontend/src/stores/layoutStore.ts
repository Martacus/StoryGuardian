import { ref } from 'vue'
import { defineStore } from 'pinia'
import { LayoutService } from '../../bindings/litguardian'
import { DashboardLayout, ViewLayout, ModuleLayout } from '../../bindings/litguardian/models'
import { moduleRegistry, getDefaultModules } from '@/modules/registry'
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
   * - Saved positions (x/y/w/h) and visibility are preserved
   * - Falls back to registry defaults when no layout exists for the view
   * - Appends newly registered modules below existing ones
   * - Drops modules that have been removed from the registry
   */
  function getViewModules(viewId: string): ModuleLayout[] {
    const savedView = layout.value.views?.[viewId]
    const registryForView = moduleRegistry.filter(m => m.views.includes(viewId))

    if (!savedView || savedView.modules.length === 0) {
      return getDefaultModules(viewId)
    }

    const registryIds = new Set(registryForView.map(m => m.id))
    const savedIds = new Set(savedView.modules.map(m => m.id))

    // Keep saved modules that still exist in registry (preserves x/y/w/h)
    const reconciled = savedView.modules.filter(m => registryIds.has(m.id))

    // Find the bottom edge of all existing modules to place new ones below
    let maxY = 0
    for (const m of reconciled) {
      maxY = Math.max(maxY, m.y + m.h)
    }

    // Append new registry modules not yet in the saved layout
    let x = 0
    let rowMaxH = 0
    for (const def of registryForView) {
      if (!savedIds.has(def.id)) {
        if (x + def.defaultW > 12) {
          x = 0
          maxY += rowMaxH
          rowMaxH = 0
        }
        reconciled.push(new ModuleLayout({
          id: def.id,
          x,
          y: maxY,
          w: def.defaultW,
          h: def.defaultH,
          visible: true,
        }))
        rowMaxH = Math.max(rowMaxH, def.defaultH)
        x += def.defaultW
      }
    }

    return reconciled
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
    toggleEditMode,
    reset,
  }
})
