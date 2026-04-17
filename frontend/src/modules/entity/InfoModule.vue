<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { Entity } from '../../../bindings/litguardian/models'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { timeAgo } from '@/lib/timeago'

const entityStore = useEntityStore()

const predefinedTypes = [
  'Character',
  'Location',
  'Item',
  'Faction',
  'Creature',
  'Event',
  'Concept',
]

// Local editable copies — only committed to disk on Save
const name = ref('')
const description = ref('')
const entityType = ref('')
// Sync when selected entity changes
watch(() => entityStore.selectedEntity, (entity) => {
  name.value = entity?.name ?? ''
  description.value = entity?.description ?? ''
  entityType.value = entity?.type ?? ''
}, { immediate: true })

const isDirty = computed(() => {
  const entity = entityStore.selectedEntity
  if (!entity) return false
  return (
    name.value !== (entity.name ?? '') ||
    description.value !== (entity.description ?? '') ||
    entityType.value !== (entity.type ?? '')
  )
})

const canSave = computed(() =>
  isDirty.value && name.value.trim().length > 0 && !entityStore.loading
)

async function save() {
  if (!canSave.value || !entityStore.selectedEntity) return

  const updated = new Entity({
    ...entityStore.selectedEntity,
    name: name.value.trim(),
    description: description.value,
    type: entityType.value.trim(),
  })

  await entityStore.updateEntity(updated)
}

function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '—'
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <div class="flex flex-col gap-4">

    <!-- Name -->
    <div class="grid gap-2">
      <Label for="entity-name">Name</Label>
      <Input
        id="entity-name"
        v-model="name"
        placeholder="Entity name"
        @keydown.enter="save"
      />
    </div>

    <!-- Type -->
    <div class="grid gap-2">
      <Label for="entity-type">Type</Label>
      <Input
        id="entity-type"
        v-model="entityType"
        placeholder="e.g. Character, Location, Item…"
        list="entity-type-suggestions"
      />
      <datalist id="entity-type-suggestions">
        <option v-for="t in predefinedTypes" :key="t" :value="t" />
      </datalist>
    </div>

    <!-- Description -->
    <div class="grid gap-2">
      <Label for="entity-description">Description</Label>
      <Textarea
        id="entity-description"
        v-model="description"
        placeholder="Describe this entity…"
        class="min-h-28 resize-y"
      />
    </div>

    <!-- Save button -->
    <div class="flex justify-end">
      <Button :disabled="!canSave" size="sm" @click="save">
        <Loader2 v-if="entityStore.loading" class="h-3 w-3 animate-spin" />
        Save
      </Button>
    </div>

    <!-- Metadata -->
    <div class="flex flex-col gap-1 border-t border-border pt-3 text-xs text-muted-foreground">
      <span>Created {{ formatDate(entityStore.selectedEntity?.createdAt) }}</span>
      <span>Last modified {{ timeAgo(entityStore.selectedEntity?.updatedAt) }}</span>
    </div>

  </div>
</template>
