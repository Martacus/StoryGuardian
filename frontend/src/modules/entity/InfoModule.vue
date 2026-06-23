<script setup lang="ts">
import { ref, computed, watch, inject, h, onBeforeUnmount, type ShallowRef, type Component } from 'vue'
import { Loader2, Pencil } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { Entity } from '../../../bindings/litguardian/internal/models'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { timeAgo } from '@/lib/timeago'

const entityStore = useEntityStore()
const moduleActions = inject<ShallowRef<Component | null>>('moduleActions')!

const predefinedTypes = [
  'Character',
  'Location',
  'Item',
  'Faction',
  'Creature',
  'Event',
  'Concept',
]

const isEditing = ref(false)
const name = ref('')
const description = ref('')
const entityType = ref('')

function syncFromStore() {
  const entity = entityStore.selectedEntity
  name.value = entity?.name ?? ''
  description.value = entity?.description ?? ''
  entityType.value = entity?.type ?? ''
}

watch(() => entityStore.selectedEntity, () => {
  syncFromStore()
  isEditing.value = false
  updateModuleActions()
}, { immediate: true })

updateModuleActions()
onBeforeUnmount(() => { moduleActions.value = null })

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

function updateModuleActions() {
  moduleActions.value = isEditing.value ? null : () => h(Button, {
    variant: 'ghost',
    size: 'icon',
    class: 'h-7 w-7 no-drag',
    onClick: startEdit,
  }, () => h(Pencil, { class: 'h-4 w-4' }))
}

function startEdit() {
  isEditing.value = true
  updateModuleActions()
}

function cancel() {
  syncFromStore()
  isEditing.value = false
  updateModuleActions()
}

async function save() {
  if (!canSave.value || !entityStore.selectedEntity) return

  const updated = new Entity({
    ...entityStore.selectedEntity,
    name: name.value.trim(),
    description: description.value,
    type: entityType.value.trim(),
  })

  await entityStore.updateEntity(updated)
  isEditing.value = false
  updateModuleActions()
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

    <!-- View mode -->
    <template v-if="!isEditing">
      <div class="grid gap-4">
        <div class="grid gap-1">
          <Label class="text-xs text-muted-foreground">Name</Label>
          <p class="text-sm font-medium">{{ name || '—' }}</p>
        </div>
        <div class="grid gap-1">
          <Label class="text-xs text-muted-foreground">Type</Label>
          <p class="text-sm">{{ entityType || '—' }}</p>
        </div>
        <div class="grid gap-1">
          <Label class="text-xs text-muted-foreground">Description</Label>
          <p class="text-sm text-muted-foreground whitespace-pre-wrap">{{ description || '—' }}</p>
        </div>
      </div>
    </template>

    <!-- Edit mode -->
    <template v-else>
      <div class="grid gap-2">
        <Label for="entity-name">Name</Label>
        <Input
          id="entity-name"
          v-model="name"
          placeholder="Entity name"
          @keydown.enter="save"
        />
      </div>

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

      <div class="grid gap-2">
        <Label for="entity-description">Description</Label>
        <Textarea
          id="entity-description"
          v-model="description"
          placeholder="Describe this entity…"
          class="min-h-28 resize-y"
        />
      </div>

      <div class="flex justify-end gap-2">
        <Button variant="outline" size="sm" @click="cancel">Cancel</Button>
        <Button :disabled="!canSave" size="sm" @click="save">
          <Loader2 v-if="entityStore.loading" class="h-3 w-3 animate-spin" />
          Save
        </Button>
      </div>
    </template>

    <!-- Metadata -->
    <div class="flex flex-col gap-1 border-t border-border pt-3 text-xs text-muted-foreground">
      <span>Created {{ formatDate(entityStore.selectedEntity?.createdAt) }}</span>
      <span>Last modified {{ timeAgo(entityStore.selectedEntity?.updatedAt) }}</span>
    </div>

  </div>
</template>
