import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { EntityService } from '../../bindings/litguardian'
import type { Entity } from '../../bindings/litguardian'
import { useWorldStore } from './worldStore'
import { useAppToast } from '@/composables/useAppToast'

export const useEntityStore = defineStore('entities', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  const entities = ref<Entity[]>([])
  const selectedEntityId = ref<string | null>(null)
  const loading = ref(false)

  // ── Computed ──────────────────────────────────────────────────────────────
  const selectedEntity = computed(() =>
    entities.value.find(e => e.id === selectedEntityId.value) ?? null
  )

  // ── Actions ───────────────────────────────────────────────────────────────

  /** Load all entities from disk for the open world. */
  async function loadEntities() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      const result = await EntityService.ListEntities(worldStore.currentWorld.path)
      entities.value = result ?? []
    } catch (e) {
      errorToast('Failed to load entities', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Create a new entity and auto-select it. */
  async function createEntity(name: string, entityType: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      const created = await EntityService.CreateEntity(worldStore.currentWorld.path, name, entityType)
      await loadEntities()
      if (created) {
        selectedEntityId.value = created.id
      }
    } catch (e) {
      errorToast('Failed to create entity', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Update an existing entity on disk and refresh the list. */
  async function updateEntity(entity: Entity) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      await EntityService.UpdateEntity(worldStore.currentWorld.path, entity)
      await loadEntities()
    } catch (e) {
      errorToast('Failed to update entity', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Delete an entity, clean up its links, and clear selection if it was selected. */
  async function deleteEntity(id: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      const { useLinkStore } = await import('./linkStore')
      await useLinkStore().deleteLinksForEntity(id)
      await EntityService.DeleteEntity(worldStore.currentWorld.path, id)
      if (selectedEntityId.value === id) {
        selectedEntityId.value = null
      }
      await loadEntities()
    } catch (e) {
      errorToast('Failed to delete entity', String(e))
    }
  }

  /** Select an entity by ID (or null to deselect / return to list). */
  function selectEntity(id: string | null) {
    selectedEntityId.value = id
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    entities.value = []
    selectedEntityId.value = null
    loading.value = false
  }

  return {
    entities,
    selectedEntityId,
    loading,
    selectedEntity,
    loadEntities,
    createEntity,
    updateEntity,
    deleteEntity,
    selectEntity,
    reset,
  }
})
