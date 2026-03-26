<script setup lang="ts">
import { ref, computed, inject, h, onBeforeUnmount, type ShallowRef, type Component } from 'vue'
import { Plus, Users } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'

const navStore = useNavigationStore()
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const moduleActions = inject<ShallowRef<Component | null>>('moduleActions')!

const entityStore = useEntityStore()

const entityCount = computed(() => entityStore.entities.length)
const previewEntities = computed(() => entityStore.entities.slice(0, 8))

const addDialogOpen = ref(false)
const newName = ref('')
const newType = ref('')

const predefinedTypes = [
  'Character',
  'Location',
  'Item',
  'Faction',
  'Creature',
  'Event',
  'Concept',
]

async function handleCreate() {
  if (!newName.value.trim()) return
  await entityStore.createEntity(newName.value.trim(), newType.value.trim())
  newName.value = ''
  newType.value = ''
  addDialogOpen.value = false
}

moduleActions.value = () => h(Button, {
  variant: 'ghost',
  size: 'icon',
  class: 'h-7 w-7 no-drag',
  onClick: () => { addDialogOpen.value = true },
}, () => h(Plus, { class: 'h-4 w-4' }))
onBeforeUnmount(() => { moduleActions.value = null })

function resetDialog() {
  newName.value = ''
  newType.value = ''
}
</script>

<template>
  <!-- Add Entity dialog (opened by header action button) -->
  <Dialog v-model:open="addDialogOpen" @update:open="val => { if (!val) resetDialog() }">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>New Entity</DialogTitle>
        <DialogDescription>Create a new entity in your world.</DialogDescription>
      </DialogHeader>
      <div class="flex flex-col gap-4 pt-2">
        <div class="grid gap-2">
          <Label for="overview-entity-name">Name</Label>
          <Input
            id="overview-entity-name"
            v-model="newName"
            placeholder="Entity name"
            @keydown.enter="handleCreate"
          />
        </div>
        <div class="grid gap-2">
          <Label for="overview-entity-type">Type</Label>
          <Input
            id="overview-entity-type"
            v-model="newType"
            placeholder="e.g. Character, Location, Item…"
            list="overview-entity-type-options"
          />
          <datalist id="overview-entity-type-options">
            <option v-for="t in predefinedTypes" :key="t" :value="t" />
          </datalist>
        </div>
        <div class="flex justify-end">
          <Button :disabled="!newName.trim() || entityStore.loading" size="sm" @click="handleCreate">
            Create
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>

  <!-- Module content -->
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
          class="flex items-center gap-2 text-sm rounded px-1 -mx-1 cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
          @click="() => navStore.navigateTo({ view: 'entities', entityId: entity.id })"
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
