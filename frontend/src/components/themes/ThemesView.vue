<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowLeft, LayoutGrid, Pencil, Trash2, Plus } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useThemeStore, STANDARD_THEME_ID } from '@/stores/themeStore'
import { useLayoutStore } from '@/stores/layoutStore'
import { getDefaultModules } from '@/modules/registry'
import { useThemeLayoutController } from '@/composables/useLayoutController'
import ModuleGrid from '@/components/modules/ModuleGrid.vue'

const themeStore = useThemeStore()
const layoutStore = useLayoutStore()

// null = list view; otherwise the id of the theme being edited.
const editingId = ref<string | null>(null)
const editingTheme = computed(() => (editingId.value ? themeStore.getTheme(editingId.value) : undefined))
const themeController = useThemeLayoutController(() => editingId.value ?? '')

/** Leave the editor, flushing any in-flight edits and clearing edit mode. */
async function closeEditor() {
  await themeController.save()
  layoutStore.editMode = false
  editingId.value = null
}

// ── Create ───────────────────────────────────────────────────────────────────
const createOpen = ref(false)
const createName = ref('')

async function createTheme() {
  const seed = themeStore.getTheme(STANDARD_THEME_ID)?.modules ?? getDefaultModules('entity-detail')
  const id = await themeStore.createTheme(createName.value, seed)
  createOpen.value = false
  createName.value = ''
  editingId.value = id
}

// ── Rename ───────────────────────────────────────────────────────────────────
const renameOpen = ref(false)
const renameId = ref<string | null>(null)
const renameName = ref('')

function openRename(id: string, name: string) {
  renameId.value = id
  renameName.value = name
  renameOpen.value = true
}

async function confirmRename() {
  if (renameId.value) await themeStore.renameTheme(renameId.value, renameName.value)
  renameOpen.value = false
  renameId.value = null
}

// ── Delete ───────────────────────────────────────────────────────────────────
const deleteOpen = ref(false)
const deleteId = ref<string | null>(null)
const deleteName = computed(() => (deleteId.value ? themeStore.getTheme(deleteId.value)?.name : ''))

function openDelete(id: string) {
  deleteId.value = id
  deleteOpen.value = true
}

async function confirmDelete() {
  if (deleteId.value) await themeStore.deleteTheme(deleteId.value)
  deleteOpen.value = false
  deleteId.value = null
}
</script>

<template>
  <!-- ── Editor ─────────────────────────────────────────────────────────── -->
  <ModuleGrid
    v-if="editingId && editingTheme"
    view-id="entity-detail"
    :controller="themeController"
  >
    <template #title>
      <div class="flex items-center gap-2">
        <Button variant="ghost" size="icon" class="h-7 w-7" @click="closeEditor()">
          <ArrowLeft class="h-4 w-4" />
        </Button>
        {{ editingTheme.name }}
        <Badge v-if="editingTheme.builtin" variant="secondary" class="text-xs">Built-in</Badge>
      </div>
    </template>
  </ModuleGrid>

  <!-- ── List ───────────────────────────────────────────────────────────── -->
  <div v-else class="flex flex-col h-full">
    <div class="flex items-center justify-between px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">Themes</h1>
      <Button size="sm" @click="createOpen = true">
        <Plus class="h-4 w-4" />
        New theme
      </Button>
    </div>

    <p class="px-6 -mt-2 pb-2 text-sm text-muted-foreground shrink-0">
      Themes are reusable entity layouts. Pick one as the default for new entities.
    </p>

    <div class="flex-1 overflow-auto px-6 pb-6">
      <div class="flex flex-col gap-2">
        <Card
          v-for="theme in themeStore.themes"
          :key="theme.id"
          class="flex flex-row items-center gap-3 p-3"
        >
          <LayoutGrid class="h-5 w-5 shrink-0 text-muted-foreground" />

          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span class="font-medium truncate">{{ theme.name }}</span>
              <Badge v-if="theme.builtin" variant="secondary" class="text-xs">Built-in</Badge>
            </div>
            <span class="text-xs text-muted-foreground">
              {{ theme.modules.length }} module{{ theme.modules.length === 1 ? '' : 's' }}
            </span>
          </div>

          <!-- Default-for-new-entities -->
          <label class="flex items-center gap-1.5 text-xs text-muted-foreground cursor-pointer select-none">
            <input
              type="radio"
              name="default-theme"
              class="accent-primary"
              :checked="themeStore.defaultThemeId === theme.id"
              @change="themeStore.setDefaultTheme(theme.id)"
            />
            Default for new entities
          </label>

          <Button variant="outline" size="sm" @click="editingId = theme.id">
            <Pencil class="h-4 w-4" />
            Edit layout
          </Button>
          <Button variant="ghost" size="sm" @click="openRename(theme.id, theme.name)">
            Rename
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            class="text-muted-foreground hover:text-destructive"
            :disabled="theme.builtin"
            title="Delete theme"
            @click="openDelete(theme.id)"
          >
            <Trash2 class="h-4 w-4" />
          </Button>
        </Card>
      </div>
    </div>
  </div>

  <!-- New theme dialog -->
  <Dialog v-model:open="createOpen">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>New theme</DialogTitle>
        <DialogDescription>Starts as a copy of the Standard layout.</DialogDescription>
      </DialogHeader>
      <div class="grid gap-2">
        <Label for="new-theme-name">Theme name</Label>
        <Input
          id="new-theme-name"
          v-model="createName"
          placeholder="e.g. Character, Location…"
          @keydown.enter="createName.trim() && createTheme()"
        />
      </div>
      <DialogFooter>
        <Button variant="outline" @click="createOpen = false">Cancel</Button>
        <Button :disabled="!createName.trim()" @click="createTheme()">Create</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- Rename dialog -->
  <Dialog v-model:open="renameOpen">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>Rename theme</DialogTitle>
      </DialogHeader>
      <div class="grid gap-2">
        <Label for="rename-theme-name">Theme name</Label>
        <Input
          id="rename-theme-name"
          v-model="renameName"
          @keydown.enter="renameName.trim() && confirmRename()"
        />
      </div>
      <DialogFooter>
        <Button variant="outline" @click="renameOpen = false">Cancel</Button>
        <Button :disabled="!renameName.trim()" @click="confirmRename()">Save</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>

  <!-- Delete confirm dialog -->
  <Dialog v-model:open="deleteOpen">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>Delete theme?</DialogTitle>
        <DialogDescription>
          “{{ deleteName }}” will be removed. Entities using it fall back to the default theme.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button variant="outline" @click="deleteOpen = false">Cancel</Button>
        <Button variant="destructive" @click="confirmDelete()">Delete</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
