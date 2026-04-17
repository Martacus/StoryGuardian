<script setup lang="ts">
import { ref, computed, inject, h, watch, onBeforeUnmount, type ShallowRef, type Component } from 'vue'
import { Plus, Link as LinkIcon } from 'lucide-vue-next'
import { useEntityStore } from '@/stores/entityStore'
import { useLinkStore } from '@/stores/linkStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

const moduleActions = inject<ShallowRef<Component | null>>('moduleActions')!
const entityStore = useEntityStore()
const linkStore = useLinkStore()
const navStore = useNavigationStore()

// Load links when selected entity changes
watch(() => entityStore.selectedEntity, (entity) => {
  if (entity) {
    linkStore.loadLinksForEntity(entity.id)
  }
}, { immediate: true })

const predefinedRelationTypes = [
  'Alliance',
  'Rivalry',
  'Parent',
  'Child',
  'Sibling',
  'Member',
  'Owner',
  'Friend',
  'Enemy',
  'Mentor',
  'Student',
  'Serves',
]

// ── Add dialog state ──────────────────────────────────────────────────────
const addDialogOpen = ref(false)
const newType = ref('')
const searchQuery = ref('')
const selectedTargetId = ref<string | null>(null)

const availableEntities = computed(() => {
  const currentId = entityStore.selectedEntity?.id
  if (!currentId) return []
  return entityStore.entities.filter(e => e.id !== currentId)
})

const filteredEntities = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return availableEntities.value
  return availableEntities.value.filter(e =>
    e.name.toLowerCase().includes(q) || e.type?.toLowerCase().includes(q)
  )
})

const selectedTargetName = computed(() => {
  if (!selectedTargetId.value) return ''
  return entityStore.entities.find(e => e.id === selectedTargetId.value)?.name ?? ''
})

/** Resolve the name of the "other" entity in a link. */
function otherEntityName(link: { fromEntityId: string; toEntityId: string }) {
  const currentId = entityStore.selectedEntity?.id
  const otherId = link.fromEntityId === currentId ? link.toEntityId : link.fromEntityId
  return entityStore.entities.find(e => e.id === otherId)?.name ?? 'Unknown'
}

async function handleCreate() {
  if (!selectedTargetId.value || !entityStore.selectedEntity) return
  await linkStore.createLink(entityStore.selectedEntity.id, selectedTargetId.value, newType.value.trim())
  resetDialog()
  addDialogOpen.value = false
}

function resetDialog() {
  newType.value = ''
  searchQuery.value = ''
  selectedTargetId.value = null
}

function navigateToRelation(linkId: string) {
  if (!entityStore.selectedEntity) return
  navStore.navigateTo({
    view: 'entities',
    entityId: entityStore.selectedEntity.id,
    linkId,
  })
}

// ── Module header action button ───────────────────────────────────────────
moduleActions.value = () => h(Button, {
  variant: 'ghost',
  size: 'icon',
  class: 'h-7 w-7 no-drag',
  onClick: () => { addDialogOpen.value = true },
}, () => h(Plus, { class: 'h-4 w-4' }))
onBeforeUnmount(() => { moduleActions.value = null })
</script>

<template>
  <!-- Add Relation dialog -->
  <Dialog v-model:open="addDialogOpen" @update:open="val => { if (!val) resetDialog() }">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Add Relation</DialogTitle>
        <DialogDescription>Create a relation from {{ entityStore.selectedEntity?.name }} to another entity.</DialogDescription>
      </DialogHeader>
      <div class="flex flex-col gap-4 pt-2">
        <!-- Target entity picker -->
        <div class="grid gap-2">
          <Label>Target Entity</Label>
          <Input
            v-model="searchQuery"
            :placeholder="selectedTargetId ? selectedTargetName : 'Search entities…'"
            @focus="selectedTargetId = null"
          />
          <ScrollArea v-if="!selectedTargetId" class="max-h-40 rounded border border-border">
            <div
              v-for="entity in filteredEntities"
              :key="entity.id"
              class="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
              @click="selectedTargetId = entity.id; searchQuery = entity.name"
            >
              <span class="truncate">{{ entity.name }}</span>
              <span v-if="entity.type" class="text-xs text-muted-foreground shrink-0">{{ entity.type }}</span>
            </div>
            <div v-if="filteredEntities.length === 0" class="px-3 py-2 text-sm text-muted-foreground">
              No entities found.
            </div>
          </ScrollArea>
        </div>

        <!-- Relation type -->
        <div class="grid gap-2">
          <Label for="relation-type">Relation Type</Label>
          <Input
            id="relation-type"
            v-model="newType"
            placeholder="e.g. Alliance, Rivalry, Parent…"
            list="relation-type-options"
            @keydown.enter="handleCreate"
          />
          <datalist id="relation-type-options">
            <option v-for="t in predefinedRelationTypes" :key="t" :value="t" />
          </datalist>
        </div>

        <div class="flex justify-end">
          <Button :disabled="!selectedTargetId || linkStore.loading" size="sm" @click="handleCreate">
            Create
          </Button>
        </div>
      </div>
    </DialogContent>
  </Dialog>

  <!-- Module content -->
  <div class="flex flex-col gap-3">
    <div v-if="linkStore.links.length === 0" class="flex flex-col items-center gap-2 py-4 text-muted-foreground">
      <LinkIcon class="h-8 w-8" />
      <p class="text-sm">No relations yet.</p>
    </div>

    <template v-else>
      <p class="text-sm text-muted-foreground">
        {{ linkStore.links.length }} {{ linkStore.links.length === 1 ? 'relation' : 'relations' }}
      </p>
      <div class="flex flex-col gap-1">
        <div
          v-for="link in linkStore.links"
          :key="link.id"
          class="flex items-center gap-2 text-sm rounded px-1 -mx-1 cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
          @click="navigateToRelation(link.id)"
        >
          <span class="truncate font-medium">{{ otherEntityName(link) }}</span>
          <span v-if="link.type" class="text-xs text-muted-foreground shrink-0">{{ link.type }}</span>
        </div>
      </div>
    </template>
  </div>
</template>
