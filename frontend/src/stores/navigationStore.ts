import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useEntityStore } from './entityStore'

export interface NavState {
  view: string
  entityId: string | null
}

export const useNavigationStore = defineStore('navigation', () => {
  const history = ref<NavState[]>([])
  const current = ref<NavState>({ view: 'overview', entityId: null })

  const canGoBack = computed(() => history.value.length > 0)
  const currentView = computed(() => current.value.view)
  const currentEntityId = computed(() => current.value.entityId)

  function isSameState(a: NavState, b: NavState) {
    return a.view === b.view && a.entityId === b.entityId
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
      current.value = { view: 'overview', entityId: null }
    }
    syncEntityStore()
  }

  function syncEntityStore() {
    const entityStore = useEntityStore()
    entityStore.selectEntity(current.value.entityId)
  }

  function reset() {
    history.value = []
    current.value = { view: 'overview', entityId: null }
  }

  return {
    history,
    current,
    canGoBack,
    currentView,
    currentEntityId,
    navigateTo,
    goBack,
    reset,
  }
})
