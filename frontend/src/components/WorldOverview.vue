<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useWorldStore } from '@/stores/worldStore'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Loader2 } from 'lucide-vue-next'
import { timeAgo } from '@/lib/timeago'

const store = useWorldStore()

// Local editable copies — only committed to disk when user clicks Save
const name = ref(store.worldMeta?.name ?? '')
const description = ref(store.worldMeta?.description ?? '')

// Sync local state if worldMeta changes (e.g. world switched)
watch(() => store.worldMeta, (meta) => {
  name.value = meta?.name ?? ''
  description.value = meta?.description ?? ''
}, { immediate: true })

// Enable Save only when something has actually changed
const isDirty = computed(() =>
  name.value !== (store.worldMeta?.name ?? '') ||
  description.value !== (store.worldMeta?.description ?? '')
)

const canSave = computed(() => isDirty.value && name.value.trim().length > 0 && !store.loading)

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
  <div class="max-w-2xl mx-auto px-8 py-8">

    <!-- Page header -->
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-xl font-semibold">World Overview</h1>
      <Button
        :disabled="!canSave"
        @click="save"
      >
        <Loader2 v-if="store.loading" class="mr-2 h-4 w-4 animate-spin" />
        Save
      </Button>
    </div>

    <!-- Error -->
    <p v-if="store.error" class="mb-4 text-sm text-destructive">{{ store.error }}</p>

    <!-- Form -->
    <div class="flex flex-col gap-6">

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
          class="min-h-36 resize-y"
        />
      </div>

      <!-- Metadata -->
      <div class="flex flex-col gap-1 pt-2 border-t border-border text-xs text-muted-foreground">
        <span>Created {{ formatDate(store.worldMeta?.createdAt) }}</span>
        <span>Last modified {{ timeAgo(store.worldMeta?.updatedAt) }}</span>
      </div>
    </div>
  </div>
</template>
