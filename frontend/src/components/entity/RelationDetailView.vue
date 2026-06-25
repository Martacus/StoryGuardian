<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ArrowLeft, ArrowLeftRight, Loader2, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useEntityStore } from '@/stores/entityStore'
import { useLinkStore } from '@/stores/linkStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { LinkService } from '../../../bindings/litguardian/internal'
import { Link } from '../../../bindings/litguardian/internal/models'
import { useWorldStore } from '@/stores/worldStore'
import ComboInput from '@/components/ui/combo-input/ComboInput.vue'

const entityStore = useEntityStore()
const linkStore = useLinkStore()
const navStore = useNavigationStore()
const worldStore = useWorldStore()

// ── Load link data ────────────────────────────────────────────────────────
const link = ref<Link | null>(null)
const loading = ref(false)

async function loadLink() {
  const linkId = navStore.currentLinkId
  if (!linkId || !worldStore.currentWorld) return
  loading.value = true
  try {
    link.value = await LinkService.GetLink(worldStore.currentWorld.path, linkId)
  } catch {
    link.value = null
  } finally {
    loading.value = false
  }
}

watch(() => navStore.currentLinkId, loadLink, { immediate: true })

// ── Local editable state ──────────────────────────────────────────────────
const editType = ref('')
const editDescription = ref('')

watch(link, (l) => {
  editType.value = l?.type ?? ''
  editDescription.value = l?.description ?? ''
}, { immediate: true })

const isDirty = computed(() => {
  if (!link.value) return false
  return editType.value !== (link.value.type ?? '') ||
    editDescription.value !== (link.value.description ?? '')
})

const canSave = computed(() => isDirty.value && !linkStore.loading)

// ── Entity name lookups ───────────────────────────────────────────────────
const fromEntity = computed(() =>
  entityStore.entities.find(e => e.id === link.value?.fromEntityId) ?? null
)
const toEntity = computed(() =>
  entityStore.entities.find(e => e.id === link.value?.toEntityId) ?? null
)

// ── Actions ───────────────────────────────────────────────────────────────
async function save() {
  if (!canSave.value || !link.value) return
  const updated = new Link({
    ...link.value,
    type: editType.value.trim(),
    description: editDescription.value,
  })
  await linkStore.updateLink(updated, entityStore.selectedEntity?.id)
  // Refresh local link state
  await loadLink()
}

async function deleteRelation() {
  if (!link.value) return
  const entityId = entityStore.selectedEntity?.id
  await linkStore.deleteLink(link.value.id, entityId)
  navStore.goBack()
}

function navigateToEntity(entityId: string) {
  navStore.navigateTo({ view: 'entities', entityId, linkId: null, categoryName: null })
}

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
</script>

<template>
  <div class="flex flex-col">
    <!-- ── Page header (matches ModuleGrid layout) ─────────────────────── -->
    <div class="flex items-center gap-2 px-6 py-4 shrink-0">
      <Button variant="ghost" size="icon" class="h-7 w-7" @click="navStore.goBack()">
        <ArrowLeft class="h-4 w-4" />
      </Button>
      <h1 class="text-xl font-semibold">
        {{ link?.type ? `Relation: ${link.type}` : 'Relation' }}
      </h1>
    </div>

    <!-- ── Content ─────────────────────────────────────────────────────── -->
    <div class="flex-1 overflow-auto px-6 pb-6">
      <div class="max-w-2xl">
        <div v-if="loading" class="flex items-center justify-center py-12">
          <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
        </div>

        <div v-else-if="!link" class="text-center py-12 text-muted-foreground">
          Relation not found.
        </div>

        <template v-else>
          <!-- Entities card -->
          <Card class="mb-4">
            <CardHeader>
              <CardTitle class="text-sm font-medium">Linked Entities</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="flex items-center gap-2 flex-wrap">
                <Button
                  variant="link"
                  class="h-auto p-0 text-sm font-medium"
                  @click="fromEntity && navigateToEntity(fromEntity.id)"
                >
                  {{ fromEntity?.name ?? 'Unknown entity' }}
                </Button>
                <span v-if="fromEntity?.type" class="text-xs text-muted-foreground">{{ fromEntity.type }}</span>
                <ArrowLeftRight class="h-4 w-4 text-muted-foreground mx-1 shrink-0" />
                <Button
                  variant="link"
                  class="h-auto p-0 text-sm font-medium"
                  @click="toEntity && navigateToEntity(toEntity.id)"
                >
                  {{ toEntity?.name ?? 'Unknown entity' }}
                </Button>
                <span v-if="toEntity?.type" class="text-xs text-muted-foreground">{{ toEntity.type }}</span>
              </div>
            </CardContent>
          </Card>

          <!-- Edit form -->
          <Card>
            <CardHeader>
              <CardTitle class="text-sm font-medium">Details</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="flex flex-col gap-4">
                <!-- Type -->
                <div class="grid gap-2">
                  <Label for="relation-detail-type">Type</Label>
                  <ComboInput
                    id="relation-detail-type"
                    v-model="editType"
                    placeholder="e.g. Alliance, Rivalry, Parent…"
                    :options="predefinedRelationTypes"
                  />
                </div>

                <!-- Description -->
                <div class="grid gap-2">
                  <Label for="relation-detail-description">Description</Label>
                  <Textarea
                    id="relation-detail-description"
                    v-model="editDescription"
                    placeholder="Describe this relation…"
                    class="min-h-28 resize-y"
                  />
                </div>

                <!-- Actions -->
                <div class="flex items-center justify-between">
                  <Button
                    variant="ghost"
                    size="sm"
                    class="text-destructive hover:text-destructive"
                    @click="deleteRelation"
                  >
                    <Trash2 class="h-3 w-3 mr-1" />
                    Delete
                  </Button>
                  <Button :disabled="!canSave" size="sm" @click="save">
                    <Loader2 v-if="linkStore.loading" class="h-3 w-3 animate-spin" />
                    Save
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </template>
      </div>
    </div>
  </div>
</template>
