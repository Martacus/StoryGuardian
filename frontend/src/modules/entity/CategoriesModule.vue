<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Tags } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { Entity } from '../../../bindings/litguardian/models'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const entityStore = useEntityStore()

const newCategory = ref('')

const categories = computed<string[]>(() =>
  entityStore.selectedEntity?.categories ?? []
)

async function addCategory() {
  const value = newCategory.value.trim()
  if (!value || !entityStore.selectedEntity) return
  if (categories.value.includes(value)) {
    newCategory.value = ''
    return
  }

  const updated = new Entity({
    ...entityStore.selectedEntity,
    categories: [...categories.value, value],
  })
  await entityStore.updateEntity(updated)
  newCategory.value = ''
}

async function removeCategory(cat: string) {
  if (!entityStore.selectedEntity) return

  const updated = new Entity({
    ...entityStore.selectedEntity,
    categories: categories.value.filter(c => c !== cat),
  })
  await entityStore.updateEntity(updated)
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <!-- Category pills -->
    <div v-if="categories.length === 0" class="flex flex-col items-center gap-2 py-4 text-muted-foreground">
      <Tags class="h-8 w-8" />
      <p class="text-sm">No categories yet.</p>
    </div>

    <div v-else class="flex flex-wrap gap-1.5">
      <Badge
        v-for="cat in categories"
        :key="cat"
        variant="secondary"
        class="cursor-pointer hover:bg-destructive/20 hover:text-destructive transition-colors"
        @click="removeCategory(cat)"
      >
        {{ cat }}
        <span class="ml-1 text-xs opacity-60">&times;</span>
      </Badge>
    </div>

    <!-- Add input -->
    <div class="flex gap-2">
      <Input
        v-model="newCategory"
        placeholder="Add category..."
        class="h-8 text-sm"
        @keydown.enter="addCategory"
      />
      <Button
        size="sm"
        variant="outline"
        class="h-8 shrink-0"
        :disabled="!newCategory.trim() || entityStore.loading"
        @click="addCategory"
      >
        Add
      </Button>
    </div>
  </div>
</template>
