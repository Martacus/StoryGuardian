<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import { useWorldStore } from '@/stores/worldStore'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { timeAgo } from '@/lib/timeago'

const store = useWorldStore()

// Local editable copies — only committed to disk on Save
const name = ref(store.worldMeta?.name ?? '')
const description = ref(store.worldMeta?.description ?? '')

// Sync when world switches
watch(() => store.worldMeta, (meta) => {
  name.value = meta?.name ?? ''
  description.value = meta?.description ?? ''
}, { immediate: true })

const isDirty = computed(() =>
  name.value !== (store.worldMeta?.name ?? '') ||
  description.value !== (store.worldMeta?.description ?? '')
)

const canSave = computed(() =>
  isDirty.value && name.value.trim().length > 0 && !store.loading
)

async function save() {
  if (!canSave.value) return
  await store.updateWorldMeta(name.value.trim(), description.value)
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
      <Label for="world-name">Name</Label>
      <Input
        id="world-name"
        v-model="name"
        placeholder="World name"
        @keydown.enter="save"
      />
    </div>

    <!-- Description -->
    <div class="grid gap-2">
      <Label for="world-description">Description</Label>
      <Textarea
        id="world-description"
        v-model="description"
        placeholder="Describe your world…"
        class="min-h-28 resize-y"
      />
    </div>

    <!-- Save button -->
    <div class="flex justify-end">
      <Button :disabled="!canSave" size="sm" @click="save">
        <Loader2 v-if="store.loading" class="h-3 w-3 animate-spin" />
        Save
      </Button>
    </div>

    <!-- Metadata -->
    <div class="flex flex-col gap-1 border-t border-border pt-3 text-xs text-muted-foreground">
      <span>Created {{ formatDate(store.worldMeta?.createdAt) }}</span>
      <span>Last modified {{ timeAgo(store.worldMeta?.updatedAt) }}</span>
    </div>

  </div>
</template>
