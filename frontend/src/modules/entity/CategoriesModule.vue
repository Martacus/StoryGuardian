<script setup lang="ts">
import { ref, computed, inject, h, onBeforeUnmount, type ShallowRef, type Component } from 'vue'
import { Plus, Tags, X } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { Entity } from '../../../bindings/litguardian/internal/models'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const moduleActions = inject<ShallowRef<Component | null>>('moduleActions')!
const entityStore = useEntityStore()

// ── Dialog state ─────────────────────────────────────────────────────────
const dialogOpen = ref(false)
const searchQuery = ref('')

// Categories on the current entity
const categories = computed<string[]>(() =>
  entityStore.selectedEntity?.categories ?? []
)

// All unique categories across every entity in the world
const allStoryCategories = computed(() => {
  const set = new Set<string>()
  for (const entity of entityStore.entities) {
    for (const cat of entity.categories ?? []) {
      set.add(cat)
    }
  }
  return [...set].sort((a, b) => a.localeCompare(b))
})

// Story categories filtered by search, excluding ones already on this entity
const filteredStoryCategories = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  const existing = new Set(categories.value)
  return allStoryCategories.value
    .filter(cat => !existing.has(cat))
    .filter(cat => !q || cat.toLowerCase().includes(q))
})

// Whether the search query is a new category (doesn't exist anywhere yet and not already assigned)
const isNewCategory = computed(() => {
  const q = searchQuery.value.trim()
  if (!q) return false
  const lower = q.toLowerCase()
  if (categories.value.some(c => c.toLowerCase() === lower)) return false
  return !allStoryCategories.value.some(c => c.toLowerCase() === lower)
})

// Whether the search query matches an existing unassigned category exactly
const isExactMatch = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return false
  return filteredStoryCategories.value.some(c => c.toLowerCase() === q)
})

async function addCategory(cat: string) {
  if (!entityStore.selectedEntity) return
  const value = cat.trim()
  if (!value || categories.value.includes(value)) return

  const updated = new Entity({
    ...entityStore.selectedEntity,
    categories: [...categories.value, value],
  })
  await entityStore.updateEntity(updated)
  searchQuery.value = ''
}

async function removeCategory(cat: string) {
  if (!entityStore.selectedEntity) return

  const updated = new Entity({
    ...entityStore.selectedEntity,
    categories: categories.value.filter(c => c !== cat),
  })
  await entityStore.updateEntity(updated)
}

function handleSearchEnter() {
  const q = searchQuery.value.trim()
  if (!q) return
  addCategory(q)
}

function resetDialog() {
  searchQuery.value = ''
}

// ── Module header action button ──────────────────────────────────────────
moduleActions.value = () => h(Button, {
  variant: 'ghost',
  size: 'icon',
  class: 'h-7 w-7 no-drag',
  onClick: () => { dialogOpen.value = true },
}, () => h(Plus, { class: 'h-4 w-4' }))
onBeforeUnmount(() => { moduleActions.value = null })
</script>

<template>
  <!-- Add Category dialog -->
  <Dialog v-model:open="dialogOpen" @update:open="val => { if (!val) resetDialog() }">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Add Category</DialogTitle>
        <DialogDescription>Pick an existing category or create a new one.</DialogDescription>
      </DialogHeader>
      <div class="flex flex-col gap-4 pt-2">
        <!-- Search bar with + button -->
        <div class="flex gap-2">
          <Input
            v-model="searchQuery"
            placeholder="Search or create category..."
            class="h-9"
            @keydown.enter="handleSearchEnter"
          />
          <Button
            variant="outline"
            size="icon"
            class="h-9 w-9 shrink-0"
            :disabled="(!isNewCategory && !isExactMatch) || entityStore.loading"
            @click="handleSearchEnter"
          >
            <Plus class="h-4 w-4" />
          </Button>
        </div>

        <Separator />

        <!-- Story categories as pills -->
        <div v-if="filteredStoryCategories.length > 0" class="flex flex-wrap gap-1.5">
          <Badge
            v-for="cat in filteredStoryCategories"
            :key="cat"
            variant="outline"
            class="cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
            @click="addCategory(cat)"
          >
            {{ cat }}
          </Badge>
        </div>

        <p v-else-if="searchQuery.trim() && !isNewCategory" class="text-sm text-muted-foreground text-center py-2">
          Already assigned to this entity.
        </p>

        <p v-else-if="allStoryCategories.length === 0 && !searchQuery.trim()" class="text-sm text-muted-foreground text-center py-2">
          No categories in this world yet. Type above to create one.
        </p>

        <p v-else-if="searchQuery.trim() && isNewCategory" class="text-sm text-muted-foreground text-center py-2">
          Press Enter or + to create "{{ searchQuery.trim() }}".
        </p>
      </div>
    </DialogContent>
  </Dialog>

  <!-- Module content: current entity's category pills -->
  <div class="flex flex-col gap-3">
    <div v-if="categories.length === 0" class="flex flex-col items-center gap-2 py-4 text-muted-foreground">
      <Tags class="h-8 w-8" />
      <p class="text-sm">No categories yet.</p>
    </div>

    <div v-else class="flex flex-wrap gap-2">
      <Badge
        v-for="cat in categories"
        :key="cat"
        variant="secondary"
        class="text-sm px-3 py-1.5 gap-2"
      >
        {{ cat }}
        <button
          class="rounded-full p-0.5 hover:bg-destructive/20 hover:text-destructive transition-colors"
          @click="removeCategory(cat)"
        >
          <X class="h-3 w-3" />
        </button>
      </Badge>
    </div>
  </div>
</template>
