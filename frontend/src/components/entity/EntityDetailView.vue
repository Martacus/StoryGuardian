<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ArrowLeft, Save } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { useThemeStore } from '@/stores/themeStore'
import { useLayoutStore } from '@/stores/layoutStore'
import { useAppToast } from '@/composables/useAppToast'
import { useEntityLayoutController } from '@/composables/useLayoutController'
import ModuleGrid from '@/components/modules/ModuleGrid.vue'

const entityStore = useEntityStore()
const navStore = useNavigationStore()
const themeStore = useThemeStore()
const layoutStore = useLayoutStore()
const { successToast } = useAppToast()

const entityName = computed(() => entityStore.selectedEntity?.name ?? 'Entity')

const controller = useEntityLayoutController()

// Local selection ('custom' or a theme id) so the trigger reflects the click
// immediately, while the async setTheme persists in the background. Kept in sync
// with the resolved theme (e.g. after auto-fork to Custom while dragging).
const selectedTheme = ref(controller.effectiveThemeId.value)
watch(() => controller.effectiveThemeId.value, (v) => { selectedTheme.value = v })

function onSelectTheme(value: unknown) { 
  const id = String(value)
  if (!id || id === 'undefined') return
  selectedTheme.value = id 
  controller.setTheme(id)
}

const userThemes = computed(() => themeStore.themes)

// ── Save current custom layout as / into a theme ────────────────────────────
const saveDialogOpen = ref(false)
const newThemeName = ref('')

async function saveAsNewTheme() {
  const id = await themeStore.createTheme(newThemeName.value, controller.getModules())
  saveDialogOpen.value = false
  newThemeName.value = ''
  // Adopt the new theme so future edits to it flow back to this entity.
  await controller.setTheme(id)
  successToast('Theme created', themeStore.getTheme(id)?.name)
}

async function updateExistingTheme(id: string) {
  await themeStore.updateThemeModules(id, controller.getModules())
  successToast('Theme updated', themeStore.getTheme(id)?.name)
}

const editableThemes = computed(() => themeStore.themes.filter(t => !t.builtin))
</script>

<template>
  <ModuleGrid view-id="entity-detail" :controller="controller">
    <template #title>
      <div class="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          class="h-7 w-7"
          @click="navStore.goBack()"
        >
          <ArrowLeft class="h-4 w-4" />
        </Button>
        {{ entityName }}
      </div>
    </template>

    <template #actions>
      <!-- Theme controls are only relevant while customizing -->
      <template v-if="layoutStore.editMode">
        <!-- Theme selector -->
        <Select :model-value="selectedTheme" @update:model-value="onSelectTheme">
          <SelectTrigger size="sm" class="w-[160px]">
            <SelectValue placeholder="Layout theme" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="custom">Custom</SelectItem>
            <SelectItem v-for="t in userThemes" :key="t.id" :value="t.id">
              {{ t.name }}
            </SelectItem>
          </SelectContent>
        </Select>

        <!-- Save / update theme from a custom layout -->
        <DropdownMenu v-if="controller.isCustom.value">
          <DropdownMenuTrigger as-child>
            <Button variant="outline" size="sm">
              <Save class="h-4 w-4" />
              Save theme
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-52">
            <DropdownMenuItem @select="saveDialogOpen = true">
              Save as new theme…
            </DropdownMenuItem>
            <template v-if="editableThemes.length">
              <DropdownMenuSeparator />
              <DropdownMenuLabel>Update existing</DropdownMenuLabel>
              <DropdownMenuItem
                v-for="t in editableThemes"
                :key="t.id"
                @select="updateExistingTheme(t.id)"
              >
                {{ t.name }}
              </DropdownMenuItem>
            </template>
          </DropdownMenuContent>
        </DropdownMenu>
      </template>
    </template>
  </ModuleGrid>

  <!-- New theme name dialog -->
  <Dialog v-model:open="saveDialogOpen">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>Save as new theme</DialogTitle>
      </DialogHeader>
      <div class="grid gap-2">
        <Label for="theme-name">Theme name</Label>
        <Input
          id="theme-name"
          v-model="newThemeName"
          placeholder="e.g. Character, Location…"
          @keydown.enter="saveAsNewTheme()"
        />
      </div>
      <DialogFooter>
        <Button variant="outline" @click="saveDialogOpen = false">Cancel</Button>
        <Button :disabled="!newThemeName.trim()" @click="saveAsNewTheme()">Create</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
