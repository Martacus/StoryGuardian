<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft, Users } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'

const entityStore = useEntityStore()
const navStore = useNavigationStore()

const categoryName = computed(() => navStore.currentCategoryName ?? '')

const entitiesInCategory = computed(() =>
  entityStore.entities.filter(e =>
    e.categories?.includes(categoryName.value)
  )
)

function navigateToEntity(entityId: string) {
  navStore.navigateTo({ view: 'entities', entityId, linkId: null, categoryName: null })
}
</script>

<template>
  <div class="flex flex-col h-full">

    <!-- Header -->
    <div class="flex items-center gap-2 px-6 py-4 shrink-0">
      <Button variant="ghost" size="icon" class="h-7 w-7" @click="navStore.goBack()">
        <ArrowLeft class="h-4 w-4" />
      </Button>
      <h1 class="text-xl font-semibold">{{ categoryName }}</h1>
      <span class="text-sm text-muted-foreground">
        {{ entitiesInCategory.length }} {{ entitiesInCategory.length === 1 ? 'entity' : 'entities' }}
      </span>
    </div>

    <!-- Entity list -->
    <ScrollArea class="flex-1 px-6 pb-6">
      <div
        v-if="entitiesInCategory.length === 0"
        class="flex flex-col items-center justify-center gap-3 py-20 text-muted-foreground"
      >
        <Users class="h-10 w-10" />
        <p class="text-sm">No entities in this category.</p>
      </div>

      <div v-else class="flex flex-col gap-1">
        <button
          v-for="entity in entitiesInCategory"
          :key="entity.id"
          class="flex items-center gap-3 rounded-md px-3 py-2.5 text-left transition-colors hover:bg-accent hover:text-accent-foreground"
          @click="navigateToEntity(entity.id)"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate">{{ entity.name }}</div>
            <div v-if="entity.type" class="text-xs text-muted-foreground truncate">
              {{ entity.type }}
            </div>
          </div>
        </button>
      </div>
    </ScrollArea>

  </div>
</template>
