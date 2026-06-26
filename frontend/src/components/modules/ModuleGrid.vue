<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount, provide } from 'vue'
import { GridLayout, GridItem } from 'grid-layout-plus'
import type { LayoutItem } from 'grid-layout-plus'
import { Settings2, Check, Plus } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import {
  useWorldLayoutController,
  LAYOUT_CONTROLLER_KEY,
  type LayoutController,
} from '@/composables/useLayoutController'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/internal/models'
import ModuleCard from './ModuleCard.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

const props = defineProps<{
  viewId: string
  // Where modules are read from / saved to. Defaults to the world layout for viewId.
  controller?: LayoutController
}>()

const layoutStore = useLayoutStore()

// Use the provided controller, or fall back to the world layout for this view.
const controller: LayoutController = props.controller ?? useWorldLayoutController(props.viewId)
provide(LAYOUT_CONTROLLER_KEY, controller)

const addModuleOpen = ref(false)

// Modules from the registry that aren't currently on the grid
const addableModules = computed(() => controller.getAddableModules())

// In edit mode show all modules (hidden ones are dimmed); in view mode only visible
const activeModules = computed<ModuleLayout[]>(() => {
  const all = controller.getModules()
  return layoutStore.editMode ? all : all.filter(m => m.visible)
})

// Local mutable array for grid-layout-plus — the library mutates this in place
// during drag/resize, so it must be a ref (not a computed).
const gridLayout = ref<LayoutItem[]>([])

// Suppresses CSS transitions on initial mount so items appear at their correct
// positions immediately instead of sliding in from 0,0. Enabled after first paint.
const isReady = ref(false)
onMounted(() => {
  requestAnimationFrame(() => {
    requestAnimationFrame(() => { isReady.value = true })
  })
})

// suppressRebuild: raised during a drag/resize cycle so the watch doesn't replace
// gridLayout while grid-layout-plus is still cleaning up its placeholder.
// Also raised during an external layout rebuild (theme switch, etc.) so the grid
// firing layout-updated back doesn't re-fork the entity to Custom.
let suppressRebuild = false

// Rebuild gridLayout whenever the source changes (world load, visibility toggle,
// editMode switch, theme switch). Suppressed while a drag/resize cycle is in flight.
watch(
  activeModules,
  (modules) => {
    if (suppressRebuild) return
    suppressRebuild = true
    gridLayout.value = modules.map(m => ({
      i: m.id,
      x: m.x,
      y: m.y,
      w: m.w,
      h: m.h,
    }))
    nextTick(() => { suppressRebuild = false })
  },
  { immediate: true, deep: true },
)

function getModuleDef(moduleId: string) {
  return moduleRegistry.find(m => m.id === moduleId)
}

function getModule(moduleId: string): ModuleLayout | undefined {
  return controller.getModules().find(m => m.id === moduleId)
}

function getComponent(moduleId: string) {
  return getModuleDef(moduleId)?.component
}

function selectModule(moduleId: string) {
  controller.addModule(moduleId)
  addModuleOpen.value = false
}

// Called by grid-layout-plus after a drag or resize completes — sync back to the
// controller. We raise suppressRebuild so the watcher doesn't replace the entire
// gridLayout array while grid-layout-plus is still cleaning up its placeholder.
function onLayoutUpdated(newLayout: LayoutItem[]) {
  if (suppressRebuild) return
  suppressRebuild = true
  controller.handleLayoutUpdate(
    newLayout.map(item => ({
      i: String(item.i),
      x: item.x,
      y: item.y,
      w: item.w,
      h: item.h,
    })),
  )
  nextTick(() => { suppressRebuild = false })
}

// Toggle edit mode. Saving (via the active controller) happens when leaving edit.
async function toggleEdit() {
  if (layoutStore.editMode) {
    await controller.save()
  }
  layoutStore.editMode = !layoutStore.editMode
}

// Flush any in-progress edits if the grid is torn down while still in edit mode
// (e.g. navigating away via the sidebar without clicking Done). The synchronous
// part of save() commits the draft to the store before unmount completes; the
// store mutation is what the next view reads.
onBeforeUnmount(() => {
  if (layoutStore.editMode) {
    controller.save()
    layoutStore.editMode = false
  }
})
</script>

<template>
  <div class="flex flex-col">

    <!-- ── Page header ──────────────────────────────────────────────────── -->
    <div class="flex items-center justify-between px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">
        <slot name="title" />
      </h1>
      <div class="flex items-center gap-2">
        <!-- View-provided controls (e.g. entity theme dropdown) -->
        <slot name="actions" />

        <!-- Add Module dialog -->
        <Dialog v-model:open="addModuleOpen">
          <DialogTrigger as-child>
            <Button
              variant="outline"
              size="sm"
              :disabled="addableModules.length === 0"
            >
              <Plus class="h-4 w-4" />
              Add Module
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Add Module</DialogTitle>
              <DialogDescription>Select a module to add to your dashboard.</DialogDescription>
            </DialogHeader>
            <div class="grid grid-cols-2 gap-3 pt-2">
              <button
                v-for="def in addableModules"
                :key="def.id"
                class="flex flex-col items-center gap-2 rounded-lg border p-4 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                @click="selectModule(def.id)"
              >
                <component :is="def.icon" class="h-8 w-8" />
                {{ def.label }}
              </button>
            </div>
          </DialogContent>
        </Dialog>

        <!-- Customize / Done -->
        <Button
          :variant="layoutStore.editMode ? 'default' : 'outline'"
          size="sm"
          @click="toggleEdit()"
        >
          <Check v-if="layoutStore.editMode" class="h-4 w-4" />
          <Settings2 v-else class="h-4 w-4" />
          {{ layoutStore.editMode ? 'Done' : 'Customize' }}
        </Button>
      </div>
    </div>

    <!-- ── Free-form grid ─────────────────────────────────────────────── -->
    <div class="px-6 pb-6">
      <GridLayout
        :class="{ 'grid-ready': isReady }"
        v-model:layout="gridLayout"
        :col-num="12"
        :row-height="40"
        :margin="[16, 16]"
        :is-draggable="layoutStore.editMode"
        :is-resizable="layoutStore.editMode"
        :vertical-compact="false"
        :use-css-transforms="true"
        :prevent-collision="false"
        @layout-updated="onLayoutUpdated"
      >
        <GridItem
          v-for="item in gridLayout"
          :key="item.i"
          :i="item.i"
          :x="item.x"
          :y="item.y"
          :w="item.w"
          :h="item.h"
          :min-w="getModuleDef(String(item.i))?.minW ?? 3"
          :min-h="getModuleDef(String(item.i))?.minH ?? 3"
          :max-w="12"
          drag-allow-from=".drag-handle"
          drag-ignore-from=".no-drag"
        >
          <ModuleCard
            :view-id="viewId"
            :module="getModule(String(item.i))!"
            class="h-full"
            @remove="controller.removeModule(String(item.i))"
          >
            <component :is="getComponent(String(item.i))" />
          </ModuleCard>
        </GridItem>
      </GridLayout>
    </div>

  </div>
</template>
