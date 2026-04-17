import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useEntityStore } from './entityStore'

export interface NavState {
  view: string
  entityId: string | null
  linkId: string | null
}

export const useNavigationStore = defineStore('navigation', () => {
  const history = ref<NavState[]>([])
  const current = ref<NavState>({ view: 'overview', entityId: null, linkId: null })

  const canGoBack = computed(() => history.value.length > 0)
  const currentView = computed(() => current.value.view)
  const currentEntityId = computed(() => current.value.entityId)
  const currentLinkId = computed(() => current.value.linkId)

  function isSameState(a: NavState, b: NavState) {
    return a.view === b.view && a.entityId === b.entityId && a.linkId === b.linkId
  }

  function navigateTo(state: NavState) {
    if (isSameState(current.value, state)) return
    history.value.push({ ...current.value })
    current.value = { ...state }
    syncEntityStore()
  }

  function goBack() {
    const prev = history.value.pop()
    if (prev) {
      current.value = prev
    } else {
      current.value = { view: 'overview', entityId: null, linkId: null }
    }
    syncEntityStore()
  }

  function syncEntityStore() {
    const entityStore = useEntityStore()
    entityStore.selectEntity(current.value.entityId)
  }

  function reset() {
    history.value = []
    current.value = { view: 'overview', entityId: null, linkId: null }
  }

  return {
    history,
    current,
    canGoBack,
    currentView,
    currentEntityId,
    currentLinkId,
    navigateTo,
    goBack,
    reset,
  }
})
