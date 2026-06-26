import { ref } from 'vue'
import { defineStore } from 'pinia'
import { ThemeService } from '../../bindings/litguardian/internal'
import { ThemeStore, LayoutTheme, ModuleLayout } from '../../bindings/litguardian/internal/models'
import { getDefaultModules } from '@/modules/registry'
import { useWorldStore } from './worldStore'
import { useLayoutStore } from './layoutStore'
import { useAppToast } from '@/composables/useAppToast'

/** The id and label of the built-in, undeletable Standard theme. */
export const STANDARD_THEME_ID = 'standard'

/** Clone a module list into fresh ModuleLayout instances (no shared references). */
function cloneModules(modules: ModuleLayout[]): ModuleLayout[] {
  return modules.map(m => new ModuleLayout({ ...m, config: m.config ? { ...m.config } : undefined }))
}

export const useThemeStore = defineStore('themes', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  const themes = ref<LayoutTheme[]>([])
  const defaultThemeId = ref<string>(STANDARD_THEME_ID)

  // ── Actions ────────────────────────────────────────────────────────────────

  /**
   * Load themes.json for the open world. Guarantees a built-in Standard theme
   * exists, seeding it (on first run / migration) from the existing entity-detail
   * layout so the user's current layout becomes Standard unchanged. Persists when
   * it has to synthesize anything.
   */
  async function loadThemes() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      const saved = await ThemeService.GetThemes(worldStore.currentWorld.path)
      themes.value = saved?.themes ?? []
      defaultThemeId.value = saved?.defaultThemeId || STANDARD_THEME_ID

      let dirty = !saved
      if (!themes.value.some(t => t.id === STANDARD_THEME_ID)) {
        // Migration / first run: seed Standard from the current entity-detail
        // layout (falls back to registry defaults inside getViewModules).
        const seed = useLayoutStore().getViewModules('entity-detail')
        themes.value.unshift(new LayoutTheme({
          id: STANDARD_THEME_ID,
          name: 'Standard',
          builtin: true,
          modules: cloneModules(seed),
        }))
        dirty = true
      }
      if (!themes.value.some(t => t.id === defaultThemeId.value)) {
        defaultThemeId.value = STANDARD_THEME_ID
        dirty = true
      }

      if (dirty) await saveThemes()
    } catch (e) {
      errorToast('Failed to load themes', String(e))
      themes.value = []
      defaultThemeId.value = STANDARD_THEME_ID
    }
  }

  /** Persist current themes to disk. */
  async function saveThemes() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      await ThemeService.SaveThemes(
        worldStore.currentWorld.path,
        new ThemeStore({ version: 1, defaultThemeId: defaultThemeId.value, themes: themes.value }),
      )
    } catch (e) {
      errorToast('Failed to save themes', String(e))
    }
  }

  /** Look up a theme by id. */
  function getTheme(id: string): LayoutTheme | undefined {
    return themes.value.find(t => t.id === id)
  }

  /** The theme used for entities that have no explicit theme set. */
  function getDefaultTheme(): LayoutTheme | undefined {
    return getTheme(defaultThemeId.value) ?? getTheme(STANDARD_THEME_ID) ?? themes.value[0]
  }

  /** Replace a theme's module list and persist. */
  async function updateThemeModules(id: string, modules: ModuleLayout[]) {
    const theme = getTheme(id)
    if (!theme) return
    theme.modules = cloneModules(modules)
    await saveThemes()
  }

  /** Create a new (non-builtin) theme seeded from the given modules. Returns its id. */
  async function createTheme(name: string, seedModules: ModuleLayout[]): Promise<string> {
    const id = crypto.randomUUID()
    themes.value.push(new LayoutTheme({
      id,
      name: name.trim() || 'Untitled theme',
      builtin: false,
      modules: cloneModules(seedModules),
    }))
    await saveThemes()
    return id
  }

  /** Rename a theme. */
  async function renameTheme(id: string, name: string) {
    const theme = getTheme(id)
    if (!theme) return
    theme.name = name.trim() || theme.name
    await saveThemes()
  }

  /** Delete a theme. Built-in themes cannot be deleted. */
  async function deleteTheme(id: string) {
    const theme = getTheme(id)
    if (!theme || theme.builtin) return
    themes.value = themes.value.filter(t => t.id !== id)
    if (defaultThemeId.value === id) defaultThemeId.value = STANDARD_THEME_ID
    await saveThemes()
  }

  /** Set the theme used for newly created entities. */
  async function setDefaultTheme(id: string) {
    if (!getTheme(id)) return
    defaultThemeId.value = id
    await saveThemes()
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    themes.value = []
    defaultThemeId.value = STANDARD_THEME_ID
  }

  return {
    themes,
    defaultThemeId,
    loadThemes,
    saveThemes,
    getTheme,
    getDefaultTheme,
    updateThemeModules,
    createTheme,
    renameTheme,
    deleteTheme,
    setDefaultTheme,
    reset,
  }
})
