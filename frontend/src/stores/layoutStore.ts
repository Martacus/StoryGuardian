import { ref } from 'vue'
import { defineStore } from 'pinia'
import { LayoutService } from '../../bindings/litguardian'
import { DashboardLayout, ViewLayout, ModuleLayout } from '../../bindings/litguardian/models'
import { moduleRegistry, getDefaultModules } from '@/modules/registry'
import { useWorldStore } from './worldStore'

export const useLayoutStore = defineStore('layout', () => {
  // ── State ─────────────────────────────────────────────────────────────────
  const layout = ref<DashboardLayout>(new DashboardLayout({ version: 1, views: {} }))
  const editMode = ref(false)

  // ── Actions ────────────────────────────────────────────────────────────────

  /** Load layout.json from disk for the open world. Falls back silently to empty state. */
  async function loadLayout() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      const saved = await LayoutService.GetLayout(worldStore.currentWorld.path)
      layout.value = saved ?? new DashboardLayout({ version: 1, views: {} })
    } catch (e) {
      console.error('Failed to load layout:', e)
      layout.value = new DashboardLayout({ version: 1, views: {} })
    }
  }

  /**
   * Returns the reconciled module list for a view.
   * - Saved order and per-module settings are preserved
   * - Falls back to registry defaults when no layout exists for the view
   * - Appends newly registered modules not yet in the saved layout
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

    // Start from saved order, dropping any IDs no longer in the registry
    const reconciled = savedView.modules.filter(m => registryIds.has(m.id))

    // Append new registry modules not yet in the saved layout (added since last open)
    for (const def of registryForView) {
      if (!savedIds.has(def.id)) {
        reconciled.push(new ModuleLayout({ id: def.id, colSpan: def.defaultColSpan, visible: true }))
      }
    }

    return reconciled
  }

  /** Replace the module order for a view. Called after drag-and-drop reorder. */
  function updateModuleOrder(viewId: string, modules: ModuleLayout[]) {
    if (!layout.value.views) layout.value.views = {}
    layout.value.views[viewId] = new ViewLayout({ modules })
  }

  /** Set the column span for a specific module in a view. */
  function setModuleColSpan(viewId: string, moduleId: string, span: number) {
    const modules = getViewModules(viewId).map(m =>
      m.id === moduleId ? new ModuleLayout({ ...m, colSpan: span }) : m
    )
    updateModuleOrder(viewId, modules)
  }

  /** Toggle visibility for a specific module in a view. */
  function toggleModuleVisibility(viewId: string, moduleId: string) {
    const modules = getViewModules(viewId).map(m =>
      m.id === moduleId ? new ModuleLayout({ ...m, visible: !m.visible }) : m
    )
    updateModuleOrder(viewId, modules)
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
      console.error('Failed to save layout:', e)
    }
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    layout.value = new DashboardLayout({ version: 1, views: {} })
    editMode.value = false
  }

  return {
    layout,
    editMode,
    loadLayout,
    getViewModules,
    updateModuleOrder,
    setModuleColSpan,
    toggleModuleVisibility,
    toggleEditMode,
    reset,
  }
})
