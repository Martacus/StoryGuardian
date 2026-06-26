import { ref, computed, type ComputedRef } from 'vue'
import { ModuleLayout, Entity } from '../../bindings/litguardian/internal/models'
import { moduleRegistry, getDefaultModules, type ModuleDefinition } from '@/modules/registry'
import { useLayoutStore, type GridItem } from '@/stores/layoutStore'
import { useEntityStore } from '@/stores/entityStore'
import { useThemeStore } from '@/stores/themeStore'

const ENTITY_VIEW = 'entity-detail'

/**
 * Abstraction over *where* a ModuleGrid's modules are read from and saved to.
 * The grid is agnostic: overview uses the world layout, entity-detail resolves a
 * theme or the entity's custom layout, and the Themes page edits a theme directly.
 */
export interface LayoutController {
  getModules(): ModuleLayout[]
  getAddableModules(): ModuleDefinition[]
  handleLayoutUpdate(items: GridItem[]): void
  addModule(id: string): void | Promise<void>
  removeModule(id: string): void | Promise<void>
  toggleVisibility(id: string): void
  getModuleConfig<T>(id: string, key: string, fallback: T): T
  setModuleConfig(id: string, key: string, value: unknown): void
  save(): Promise<void>
}

/** Injection key used by ModuleGrid (provide) and ModuleCard (inject). */
export const LAYOUT_CONTROLLER_KEY = 'layoutController'

/** Clone a module list into fresh ModuleLayout instances. */
function cloneModules(modules: ModuleLayout[]): ModuleLayout[] {
  return modules.map(m => new ModuleLayout({ ...m, config: m.config ? { ...m.config } : undefined }))
}

/** Drop modules no longer in the registry; fall back to defaults when empty. */
function reconcile(viewId: string, modules: ModuleLayout[] | undefined): ModuleLayout[] {
  if (!modules || modules.length === 0) return getDefaultModules(viewId)
  const registryIds = new Set(moduleRegistry.filter(m => m.views.includes(viewId)).map(m => m.id))
  return modules.filter(m => registryIds.has(m.id))
}

function addableFrom(viewId: string, current: ModuleLayout[]): ModuleDefinition[] {
  const currentIds = new Set(current.map(m => m.id))
  return moduleRegistry.filter(m => m.views.includes(viewId) && !currentIds.has(m.id))
}

// ── World controller (overview & any plain view) ─────────────────────────────
// Thin wrapper over layoutStore — preserves the original behavior exactly.
export function useWorldLayoutController(viewId: string): LayoutController {
  const layoutStore = useLayoutStore()
  return {
    getModules: () => layoutStore.getViewModules(viewId),
    getAddableModules: () => layoutStore.getAddableModules(viewId),
    handleLayoutUpdate: (items) => layoutStore.handleLayoutUpdate(viewId, items),
    addModule: (id) => layoutStore.addModule(viewId, id),
    removeModule: (id) => layoutStore.removeModule(viewId, id),
    toggleVisibility: (id) => layoutStore.toggleModuleVisibility(viewId, id),
    getModuleConfig: (id, key, fallback) => layoutStore.getModuleConfig(viewId, id, key, fallback),
    setModuleConfig: (id, key, value) => layoutStore.setModuleConfig(viewId, id, key, value),
    save: () => layoutStore.saveLayout(),
  }
}

// ── Entity controller (entity-detail) ────────────────────────────────────────
// Resolves the active entity's theme (or its custom layout). Any edit auto-forks
// the entity to Custom, seeded as a copy of whatever it was showing.
export interface EntityLayoutController extends LayoutController {
  effectiveThemeId: ComputedRef<string>
  isCustom: ComputedRef<boolean>
  setTheme: (id: string) => Promise<void>
}

export function useEntityLayoutController(): EntityLayoutController {
  const entityStore = useEntityStore()
  const themeStore = useThemeStore()

  // In-memory edits while customizing; flushed to disk on save().
  const draftModules = ref<ModuleLayout[] | null>(null)
  const draftThemeId = ref<string | null>(null)

  const entity = computed(() => entityStore.selectedEntity)

  const effectiveThemeId = computed(() => {
    if (draftThemeId.value !== null) return draftThemeId.value
    return entity.value?.layoutThemeId || themeStore.defaultThemeId
  })

  const isCustom = computed(() => effectiveThemeId.value === 'custom')

  /** Resolve the persisted module list for the entity's effective theme. */
  function resolvePersisted(): ModuleLayout[] {
    const e = entity.value
    const id = effectiveThemeId.value
    if (id === 'custom') return reconcile(ENTITY_VIEW, e?.customLayout)
    const theme = themeStore.getTheme(id) ?? themeStore.getDefaultTheme()
    return reconcile(ENTITY_VIEW, theme?.modules)
  }

  function getModules(): ModuleLayout[] {
    return draftModules.value ?? resolvePersisted()
  }

  /** Begin a draft, forking the entity to Custom seeded from the current layout. */
  function ensureDraft() {
    if (draftModules.value) return
    draftModules.value = cloneModules(resolvePersisted())
    draftThemeId.value = 'custom'
  }

  async function persist(themeId: string, customLayout: ModuleLayout[]) {
    const e = entity.value
    if (!e) return
    await entityStore.updateEntity(new Entity({ ...e, layoutThemeId: themeId, customLayout }))
  }

  async function save() {
    if (!draftModules.value) return
    // draftThemeId is always 'custom' when a drag/resize draft exists (set by ensureDraft).
    // If setTheme was called since the last drag, it already persisted and cleared the draft,
    // so this branch is unreachable after a theme switch — but be explicit to avoid regression.
    await persist(draftThemeId.value ?? 'custom', cloneModules(draftModules.value))
    draftModules.value = null
    draftThemeId.value = null
  }

  async function setTheme(id: string) {
    const e = entity.value
    if (!e) return
    // Seed custom from the layout currently shown (copy current theme).
    const seed = id === 'custom' ? cloneModules(resolvePersisted()) : (e.customLayout ?? [])
    draftModules.value = null
    draftThemeId.value = null

    console.log('[theme] setTheme ->', id, 'entity:', entity.value?.id) 
    console.log('[theme] persisted, seed length:', seed.length)

    await persist(id, seed)
  }

  return {
    effectiveThemeId,
    isCustom,
    setTheme,
    getModules,
    getAddableModules: () => addableFrom(ENTITY_VIEW, getModules()),
    handleLayoutUpdate(items) {
      ensureDraft()
      draftModules.value = draftModules.value!.map(m => {
        const u = items.find(i => i.i === m.id)
        return u ? new ModuleLayout({ ...m, x: u.x, y: u.y, w: u.w, h: u.h }) : m
      })
    },
    async addModule(id) {
      ensureDraft()
      if (draftModules.value!.some(m => m.id === id)) return
      const def = moduleRegistry.find(m => m.id === id)
      if (!def) return
      const maxY = draftModules.value!.reduce((mx, m) => Math.max(mx, m.y + m.h), 0)
      draftModules.value = [
        ...draftModules.value!,
        new ModuleLayout({ id, x: 0, y: maxY, w: def.defaultW, h: def.defaultH, visible: true }),
      ]
      await save()
    },
    async removeModule(id) {
      ensureDraft()
      draftModules.value = draftModules.value!.filter(m => m.id !== id)
      await save()
    },
    toggleVisibility(id) {
      ensureDraft()
      draftModules.value = draftModules.value!.map(m =>
        m.id === id ? new ModuleLayout({ ...m, visible: !m.visible }) : m,
      )
    },
    getModuleConfig(id, key, fallback) {
      const m = getModules().find(x => x.id === id)
      return (m?.config?.[key] as typeof fallback) ?? fallback
    },
    setModuleConfig(id, key, value) {
      ensureDraft()
      draftModules.value = draftModules.value!.map(m =>
        m.id === id ? new ModuleLayout({ ...m, config: { ...(m.config ?? {}), [key]: value } }) : m,
      )
    },
    save,
  }
}

// ── Theme controller (Themes page editor) ────────────────────────────────────
// Edits a theme's module list directly (shared across all entities on it).
export function useThemeLayoutController(getThemeId: () => string): LayoutController {
  const themeStore = useThemeStore()
  const draft = ref<ModuleLayout[] | null>(null)

  function resolved(): ModuleLayout[] {
    return reconcile(ENTITY_VIEW, themeStore.getTheme(getThemeId())?.modules)
  }
  function getModules(): ModuleLayout[] {
    return draft.value ?? resolved()
  }
  function ensureDraft() {
    if (!draft.value) draft.value = cloneModules(resolved())
  }
  async function save() {
    if (!draft.value) return
    await themeStore.updateThemeModules(getThemeId(), draft.value)
    draft.value = null
  }

  return {
    getModules,
    getAddableModules: () => addableFrom(ENTITY_VIEW, getModules()),
    handleLayoutUpdate(items) {
      ensureDraft()
      draft.value = draft.value!.map(m => {
        const u = items.find(i => i.i === m.id)
        return u ? new ModuleLayout({ ...m, x: u.x, y: u.y, w: u.w, h: u.h }) : m
      })
    },
    async addModule(id) {
      ensureDraft()
      if (draft.value!.some(m => m.id === id)) return
      const def = moduleRegistry.find(m => m.id === id)
      if (!def) return
      const maxY = draft.value!.reduce((mx, m) => Math.max(mx, m.y + m.h), 0)
      draft.value = [
        ...draft.value!,
        new ModuleLayout({ id, x: 0, y: maxY, w: def.defaultW, h: def.defaultH, visible: true }),
      ]
      await save()
    },
    async removeModule(id) {
      ensureDraft()
      draft.value = draft.value!.filter(m => m.id !== id)
      await save()
    },
    toggleVisibility(id) {
      ensureDraft()
      draft.value = draft.value!.map(m =>
        m.id === id ? new ModuleLayout({ ...m, visible: !m.visible }) : m,
      )
    },
    getModuleConfig(id, key, fallback) {
      const m = getModules().find(x => x.id === id)
      return (m?.config?.[key] as typeof fallback) ?? fallback
    },
    setModuleConfig(id, key, value) {
      ensureDraft()
      draft.value = draft.value!.map(m =>
        m.id === id ? new ModuleLayout({ ...m, config: { ...(m.config ?? {}), [key]: value } }) : m,
      )
    },
    save,
  }
}
