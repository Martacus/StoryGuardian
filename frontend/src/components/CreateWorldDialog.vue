<script setup lang="ts">
import { ref, computed } from 'vue'
import { useWorldStore } from '@/stores/worldStore'
import { WorldService } from '../../bindings/litguardian'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { FolderOpen, Loader2 } from 'lucide-vue-next'

const open = defineModel<boolean>('open', { default: false })

const store = useWorldStore()
const name = ref('')
const folderPath = ref('')
const picking = ref(false)
// Local error — isolated from the global store so stale errors from other
// flows (open world, update meta, etc.) never bleed into this dialog.
const dialogError = ref<string | null>(null)

const canCreate = computed(() => name.value.trim().length > 0 && folderPath.value.length > 0)

async function pickFolder() {
  picking.value = true
  try {
    const path = await WorldService.SelectFolder()
    if (path) folderPath.value = path
  } finally {
    picking.value = false
  }
}

async function handleCreate() {
  if (!canCreate.value || store.loading) return
  dialogError.value = null
  await store.createWorld(name.value.trim(), folderPath.value)
  if (store.error) {
    dialogError.value = store.error
    return
  }
  open.value = false
  name.value = ''
  folderPath.value = ''
}

function handleCancel() {
  open.value = false
  name.value = ''
  folderPath.value = ''
  dialogError.value = null
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Create New World</DialogTitle>
        <DialogDescription>
          Give your world a name and choose an empty folder to save it in.
        </DialogDescription>
      </DialogHeader>

      <div class="grid gap-4 py-2">
        <!-- World name -->
        <div class="grid gap-2">
          <Label for="world-name">World Name</Label>
          <Input
            id="world-name"
            v-model="name"
            placeholder="My Fantasy World"
            @keydown.enter="handleCreate"
          />
        </div>

        <!-- Folder picker -->
        <div class="grid gap-2">
          <Label>Folder</Label>
          <div class="flex gap-2">
            <Input
              :value="folderPath"
              readonly
              placeholder="No folder selected"
              class="flex-1 cursor-default text-muted-foreground"
            />
            <Button variant="outline" :disabled="picking" @click="pickFolder">
              <Loader2 v-if="picking" class="h-4 w-4 animate-spin" />
              <FolderOpen v-else class="h-4 w-4" />
              <span class="ml-2">Browse</span>
            </Button>
          </div>
        </div>

        <!-- Error -->
        <p v-if="dialogError" class="text-sm text-destructive">{{ dialogError }}</p>
      </div>

      <DialogFooter>
        <Button variant="ghost" @click="handleCancel">Cancel</Button>
        <Button :disabled="!canCreate || store.loading" @click="handleCreate">
          <Loader2 v-if="store.loading" class="mr-2 h-4 w-4 animate-spin" />
          Create World
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
