import { ref } from 'vue'
import { defineStore } from 'pinia'
import { LinkService } from '../../bindings/litguardian/internal'
import type { Link } from '../../bindings/litguardian/internal'
import { useWorldStore } from './worldStore'
import { useAppToast } from '@/composables/useAppToast'

export const useLinkStore = defineStore('links', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  const links = ref<Link[]>([])
  const loading = ref(false)

  // ── Actions ───────────────────────────────────────────────────────────────

  /** Load all links for a specific entity. */
  async function loadLinksForEntity(entityId: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      const result = await LinkService.GetLinksForEntity(worldStore.currentWorld.path, entityId)
      links.value = result ?? []
    } catch (e) {
      errorToast('Failed to load relations', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Create a new link between two entities. */
  async function createLink(fromEntityId: string, toEntityId: string, linkType: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      await LinkService.CreateLink(worldStore.currentWorld.path, fromEntityId, toEntityId, linkType)
      await loadLinksForEntity(fromEntityId)
    } catch (e) {
      errorToast('Failed to create relation', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Update an existing link (type, description). */
  async function updateLink(link: Link, reloadForEntityId?: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      await LinkService.UpdateLink(worldStore.currentWorld.path, link)
      if (reloadForEntityId) {
        await loadLinksForEntity(reloadForEntityId)
      }
    } catch (e) {
      errorToast('Failed to update relation', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Delete a link by ID. */
  async function deleteLink(id: string, reloadForEntityId?: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      await LinkService.DeleteLink(worldStore.currentWorld.path, id)
      if (reloadForEntityId) {
        await loadLinksForEntity(reloadForEntityId)
      }
    } catch (e) {
      errorToast('Failed to delete relation', String(e))
    }
  }

  /** Delete all links for an entity (cleanup on entity delete). */
  async function deleteLinksForEntity(entityId: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      await LinkService.DeleteLinksForEntity(worldStore.currentWorld.path, entityId)
    } catch (e) {
      errorToast('Failed to clean up relations', String(e))
    }
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    links.value = []
    loading.value = false
  }

  return {
    links,
    loading,
    loadLinksForEntity,
    createLink,
    updateLink,
    deleteLink,
    deleteLinksForEntity,
    reset,
  }
})
