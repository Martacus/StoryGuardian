<script setup lang="ts">
import { ref } from 'vue'
import { Plus, Trash2, Users } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import EntityTypeInput from '@/components/entity/EntityTypeInput.vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

const entityStore = useEntityStore()
const navStore = useNavigationStore()

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

function resetDialog() {
  newName.value = ''
  newType.value = ''
}

async function handleDelete(e: MouseEvent, id: string) {
  e.stopPropagation()
  await entityStore.deleteEntity(id)
}
</script>

<template>
  <div class="flex flex-col h-full">

    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">Entities</h1>
      <Dialog v-model:open="addDialogOpen" @update:open="val => { if (!val) resetDialog() }">
        <DialogTrigger as-child>
          <Button size="sm">
            <Plus class="h-4 w-4" />
            Add Entity
          </Button>
        </DialogTrigger>
        <DialogContent class="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>New Entity</DialogTitle>
            <DialogDescription>Create a new entity in your world.</DialogDescription>
          </DialogHeader>
          <div class="flex flex-col gap-4 pt-2">
            <div class="grid gap-2">
              <Label for="entity-name">Name</Label>
              <Input
                id="entity-name"
                v-model="newName"
                placeholder="Entity name"
                @keydown.enter="handleCreate"
              />
            </div>
            <div class="grid gap-2">
              <Label for="entity-type">Type</Label>
              <EntityTypeInput
                id="entity-type"
                v-model="newType"
                placeholder="e.g. Character, Location, Item…"
                :options="predefinedTypes"
                @enter="handleCreate"
              />
            </div>
            <div class="flex justify-end">
              <Button :disabled="!newName.trim() || entityStore.loading" size="sm" @click="handleCreate">
                Create
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>

    <!-- Entity list -->
    <ScrollArea class="flex-1 px-6 pb-6">
      <!-- Empty state -->
      <div
        v-if="entityStore.entities.length === 0 && !entityStore.loading"
        class="flex flex-col items-center justify-center gap-3 py-20 text-muted-foreground"
      >
        <Users class="h-10 w-10" />
        <p class="text-sm">No entities yet. Create your first entity to get started.</p>
      </div>

      <!-- Entity rows -->
      <div v-else class="flex flex-col gap-1">
        <button
          v-for="entity in entityStore.entities"
          :key="entity.id"
          class="flex items-center gap-3 rounded-md px-3 py-2.5 text-left transition-colors hover:bg-accent hover:text-accent-foreground group"
          @click="navStore.navigateTo({ view: 'entities', entityId: entity.id, linkId: null, categoryName: null })"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate">{{ entity.name }}</div>
            <div v-if="entity.type" class="text-xs text-muted-foreground truncate">
              {{ entity.type }}
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            class="h-7 w-7 shrink-0 opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-destructive"
            @click="(e: MouseEvent) => handleDelete(e, entity.id)"
          >
            <Trash2 class="h-3.5 w-3.5" />
          </Button>
        </button>
      </div>
    </ScrollArea>

  </div>
</template>
