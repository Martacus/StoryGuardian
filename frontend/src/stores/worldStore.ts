import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { AppConfigService, WorldService } from '../../bindings/litguardian'
import type { WorldInfo, WorldMeta } from '../../bindings/litguardian'
import { useLayoutStore } from './layoutStore'
import { useImageStore } from './imageStore'
import { useEntityStore } from './entityStore'
import { useLinkStore } from './linkStore'
import { useNavigationStore } from './navigationStore'
import { useAppToast } from '@/composables/useAppToast'

export const useWorldStore = defineStore('world', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  // null = no world open → show WelcomeScreen
  const currentWorld = ref<WorldInfo | null>(null)
  // Full metadata from world.json — populated whenever a world is open
  const worldMeta = ref<WorldMeta | null>(null)
  const recentWorlds = ref<WorldInfo[]>([])
  const loading = ref(false)

  // ── Computed ───────────────────────────────────────────────────────────────
  // Drives App.vue conditional rendering: WelcomeScreen vs dashboard
  const hasOpenWorld = computed(() => currentWorld.value !== null)

  // ── Actions ────────────────────────────────────────────────────────────────

  /** Fetch the recents list from disk and populate the store. */
  async function loadRecentWorlds() {
    try {
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
    } catch (e) {
      errorToast('Failed to load recent worlds', String(e))
    }
  }

  /** Read world.json for the currently open world and cache in worldMeta. */
  async function loadWorldMeta() {
    if (!currentWorld.value) return
    try {
      worldMeta.value = await WorldService.GetWorldMeta(currentWorld.value.path)
    } catch (e) {
      errorToast('Failed to load world metadata', String(e))
    }
  }

  /**
   * Update world name + description on disk and refresh worldMeta.
   * Also syncs the name into currentWorld so the top nav stays up to date.
   */
  async function updateWorldMeta(name: string, description: string) {
    if (!currentWorld.value) return
    loading.value = true
    try {
      const updated = await WorldService.UpdateWorldMeta(currentWorld.value.path, name, description)
      if (!updated) return
      worldMeta.value = updated
      // Keep currentWorld.name in sync so the top nav reflects the new name
      currentWorld.value = { ...currentWorld.value, name: updated.name }
    } catch (e) {
      errorToast('Failed to update world', String(e))
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new world on disk, add it to recents, and open it.
   * Orchestrates: WorldService.CreateWorld → AppConfigService.AddRecentWorld → set currentWorld → load meta
   */
  async function createWorld(name: string, folderPath: string) {
    loading.value = true
    try {
      const worldInfo = await WorldService.CreateWorld(name, folderPath)
      if (!worldInfo) return
      await AppConfigService.AddRecentWorld(worldInfo.name, worldInfo.path)
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
      currentWorld.value = worldInfo
      await loadWorldMeta()
      await useLayoutStore().loadLayout()
      await useImageStore().loadImages()
      await useEntityStore().loadEntities()
    } catch (e) {
      errorToast('Failed to create world', String(e))
    } finally {
      loading.value = false
    }
  }

  /**
   * Open an existing world folder, add it to recents, and set it as current.
   * Orchestrates: WorldService.OpenWorld → AppConfigService.AddRecentWorld → set currentWorld → load meta
   */
  async function openWorld(folderPath: string) {
    loading.value = true
    try {
      const worldInfo = await WorldService.OpenWorld(folderPath)
      if (!worldInfo) return
      await AppConfigService.AddRecentWorld(worldInfo.name, worldInfo.path)
      recentWorlds.value = await AppConfigService.GetRecentWorlds()
      currentWorld.value = worldInfo
      await loadWorldMeta()
      await useLayoutStore().loadLayout()
      await useImageStore().loadImages()
      await useEntityStore().loadEntities()
    } catch (e) {
      errorToast('Failed to open world', String(e))
    } finally {
      loading.value = false
    }
  }

  /**
   * Remove a world entry from recents (does not delete files on disk).
   * Optimistically removes from local state immediately.
   */
  async function removeRecentWorld(path: string) {
    // Optimistic update — remove immediately so UI responds instantly
    recentWorlds.value = recentWorlds.value.filter(w => w.path !== path)
    try {
      await AppConfigService.RemoveRecentWorld(path)
    } catch (e) {
      // Roll back on failure
      await loadRecentWorlds()
      errorToast('Failed to remove world from recents', String(e))
    }
  }

  /** Close the current world and return to the WelcomeScreen. */
  function closeWorld() {
    currentWorld.value = null
    worldMeta.value = null
    useLayoutStore().reset()
    useImageStore().reset()
    useEntityStore().reset()
    useLinkStore().reset()
    useNavigationStore().reset()
  }

  return {
    // State
    currentWorld,
    worldMeta,
    recentWorlds,
    loading,
    // Computed
    hasOpenWorld,
    // Actions
    loadRecentWorlds,
    loadWorldMeta,
    updateWorldMeta,
    createWorld,
    openWorld,
    removeRecentWorld,
    closeWorld,
  }
})
