<script setup lang="ts">
import { computed } from 'vue'
import { Users } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'

const entityStore = useEntityStore()

const entityCount = computed(() => entityStore.entities.length)
const previewEntities = computed(() => entityStore.entities.slice(0, 8))
</script>

<template>
  <div class="flex flex-col gap-3">
    <div v-if="entityCount === 0" class="flex flex-col items-center gap-2 py-4 text-muted-foreground">
      <Users class="h-8 w-8" />
      <p class="text-sm">No entities yet.</p>
    </div>

    <template v-else>
      <p class="text-sm text-muted-foreground">
        {{ entityCount }} {{ entityCount === 1 ? 'entity' : 'entities' }}
      </p>
      <div class="flex flex-col gap-1">
        <div
          v-for="entity in previewEntities"
          :key="entity.id"
          class="flex items-center gap-2 text-sm"
        >
          <span class="truncate">{{ entity.name }}</span>
          <span v-if="entity.type" class="text-xs text-muted-foreground shrink-0">
            {{ entity.type }}
          </span>
        </div>
        <p v-if="entityCount > 8" class="text-xs text-muted-foreground">
          and {{ entityCount - 8 }} more…
        </p>
      </div>
    </template>
  </div>
</template>
