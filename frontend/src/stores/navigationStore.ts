import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useEntityStore } from './entityStore'

export interface NavState {
  view: string
  entityId: string | null
  linkId: string | null
  categoryName: string | null
}

export const useNavigationStore = defineStore('navigation', () => {
  const history = ref<NavState[]>([])
  const current = ref<NavState>({ view: 'overview', entityId: null, linkId: null, categoryName: null })

  const canGoBack = computed(() => history.value.length > 0)
  const currentView = computed(() => current.value.view)
  const currentEntityId = computed(() => current.value.entityId)
  const currentLinkId = computed(() => current.value.linkId)
  const currentCategoryName = computed(() => current.value.categoryName)

  function isSameState(a: NavState, b: NavState) {
    return a.view === b.view && a.entityId === b.entityId && a.linkId === b.linkId && a.categoryName === b.categoryName
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
      current.value = { view: 'overview', entityId: null, linkId: null, categoryName: null }
    }
    syncEntityStore()
  }

  function syncEntityStore() {
    const entityStore = useEntityStore()
    entityStore.selectEntity(current.value.entityId)
  }

  function reset() {
    history.value = []
    current.value = { view: 'overview', entityId: null, linkId: null, categoryName: null }
  }

  return {
    history,
    current,
    canGoBack,
    currentView,
    currentEntityId,
    currentLinkId,
    currentCategoryName,
    navigateTo,
    goBack,
    reset,
  }
})
