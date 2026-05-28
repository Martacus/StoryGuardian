<script setup lang="ts">
import { computed } from 'vue'
import { FolderTree, X } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Entity } from '../../../bindings/litguardian/internal/models'

const entityStore = useEntityStore()
const navStore = useNavigationStore()

/** All unique categories across every entity, with entity count. */
const categoriesWithCount = computed(() => {
  const map = new Map<string, number>()
  for (const entity of entityStore.entities) {
    for (const cat of entity.categories ?? []) {
      map.set(cat, (map.get(cat) ?? 0) + 1)
    }
  }
  return [...map.entries()]
    .sort((a, b) => a[0].localeCompare(b[0]))
    .map(([name, count]) => ({ name, count }))
})

function navigateToCategory(name: string) {
  navStore.navigateTo({ view: 'categories', entityId: null, linkId: null, categoryName: name })
}

async function deleteCategory(e: MouseEvent, name: string) {
  e.stopPropagation()
  // Remove this category from every entity that has it
  const affected = entityStore.entities.filter(ent =>
    ent.categories?.includes(name)
  )
  for (const ent of affected) {
    const updated = new Entity({
      ...ent,
      categories: ent.categories.filter(c => c !== name),
    })
    await entityStore.updateEntity(updated)
  }
}
</script>

<template>
  <div class="flex flex-col h-full">

    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">Categories</h1>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-auto px-6 pb-6">

      <!-- Empty state -->
      <div
        v-if="categoriesWithCount.length === 0"
        class="flex flex-col items-center justify-center gap-3 py-20 text-muted-foreground"
      >
        <FolderTree class="h-10 w-10" />
        <p class="text-sm">No categories yet. Add categories to entities to see them here.</p>
      </div>

      <!-- Category pills -->
      <div v-else class="flex flex-wrap gap-2">
        <Badge
          v-for="cat in categoriesWithCount"
          :key="cat.name"
          variant="secondary"
          class="cursor-pointer text-sm px-3 py-1.5 gap-2 hover:bg-accent hover:text-accent-foreground transition-colors"
          @click="navigateToCategory(cat.name)"
        >
          {{ cat.name }}
          <span class="text-xs text-muted-foreground">({{ cat.count }})</span>
          <button
            class="ml-0.5 rounded-full p-0.5 hover:bg-destructive/20 hover:text-destructive transition-colors"
            @click="(e: MouseEvent) => deleteCategory(e, cat.name)"
          >
            <X class="h-3 w-3" />
          </button>
        </Badge>
      </div>

    </div>
  </div>
</template>
