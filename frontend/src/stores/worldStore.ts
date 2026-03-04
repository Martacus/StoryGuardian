import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { AppConfigService, WorldService } from '../../bindings/litguardian'
import type { WorldInfo } from '../../bindings/litguardian'

export const useWorldStore = defineStore('world', () => {
  // ── State ─────────────────────────────────────────────────────────────────
  // null = no world open → show WelcomeScreen
  const currentWorld = ref<WorldInfo | null>(null)
  const recentWorlds = ref<WorldInfo[]>([])
  const error = ref<string | null>(null)
  const loading = ref(false)

  // ── Computed ───────────────────────────────────────────────────────────────
  // Drives App.vue conditional rendering: WelcomeScreen vs editor
  const hasOpenWorld = computed(() => currentWorld.value !== null)

  // ── Actions ────────────────────────────────────────────────────────────────

  /** Fetch the recents list from disk and populate the store. */
  async function loadRecentWorlds() {
    error.value = null
    try {
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
    } catch (e) {
      error.value = String(e)
    }
  }

  /**
   * Create a new world on disk, add it to recents, and open it.
   * Orchestrates: WorldService.CreateWorld → AppConfigService.AddRecentWorld → set currentWorld
   */
  async function createWorld(name: string, folderPath: string) {
    error.value = null
    loading.value = true
    try {
      const worldInfo = await WorldService.CreateWorld(name, folderPath)
      if (!worldInfo) return
      await AppConfigService.AddRecentWorld(worldInfo.name, worldInfo.path)
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
      currentWorld.value = worldInfo
    } catch (e) {
      error.value = String(e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Open an existing world folder, add it to recents, and set it as current.
   * Orchestrates: WorldService.OpenWorld → AppConfigService.AddRecentWorld → set currentWorld
   */
  async function openWorld(folderPath: string) {
    error.value = null
    loading.value = true
    try {
      const worldInfo = await WorldService.OpenWorld(folderPath)
      if (!worldInfo) return
      await AppConfigService.AddRecentWorld(worldInfo.name, worldInfo.path)
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
      currentWorld.value = worldInfo
    } catch (e) {
      error.value = String(e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Remove a world entry from recents (does not delete files on disk).
   * Optimistically removes from local state immediately.
   */
  async function removeRecentWorld(path: string) {
    error.value = null
    // Optimistic update — remove immediately so UI responds instantly
    recentWorlds.value = recentWorlds.value.filter(w => w.path !== path)
    try {
      await AppConfigService.RemoveRecentWorld(path)
    } catch (e) {
      // Roll back on failure
      await loadRecentWorlds()
      error.value = String(e)
    }
  }

  /** Close the current world and return to the WelcomeScreen. */
  function closeWorld() {
    currentWorld.value = null
    error.value = null
  }

  return {
    // State
    currentWorld,
    recentWorlds,
    error,
    loading,
    // Computed
    hasOpenWorld,
    // Actions
    loadRecentWorlds,
    createWorld,
    openWorld,
    removeRecentWorld,
    closeWorld,
  }
})
