<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { GridLayout, GridItem } from 'grid-layout-plus'
import type { LayoutItem } from 'grid-layout-plus'
import { Settings2, Check, Plus } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/models'
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
}>()

const layoutStore = useLayoutStore()

const addModuleOpen = ref(false)

// Modules from the registry that aren't currently on the grid
const addableModules = computed(() => layoutStore.getAddableModules(props.viewId))

// In edit mode show all modules (hidden ones are dimmed); in view mode only visible
const activeModules = computed<ModuleLayout[]>(() => {
  const all = layoutStore.getViewModules(props.viewId)
  return layoutStore.editMode ? all : all.filter(m => m.visible)
})

// Local mutable array for grid-layout-plus — the library mutates this in place
// during drag/resize, so it must be a ref (not a computed).
const gridLayout = ref<LayoutItem[]>([])

// Guard flag: prevents the watch from replacing gridLayout mid-drag.
// When layout-updated fires we update the store, which causes activeModules to
// recompute. Without this guard the watch would immediately replace gridLayout
// with a fresh array, orphaning the placeholder <div> in the DOM and leaving it
// sitting on top of cards (blocking all interaction) until the page is reloaded.
let suppressRebuild = false

// Rebuild gridLayout whenever the store changes (world load, visibility toggle,
// editMode switch). Suppressed while a drag/resize cycle is in flight.
watch(
  activeModules,
  (modules) => {
    if (suppressRebuild) return
    gridLayout.value = modules.map(m => ({
      i: m.id,
      x: m.x,
      y: m.y,
      w: m.w,
      h: m.h,
    }))
  },
  { immediate: true, deep: true },
)

function getModuleDef(moduleId: string) {
  return moduleRegistry.find(m => m.id === moduleId)
}

function getModule(moduleId: string): ModuleLayout | undefined {
  return layoutStore.getViewModules(props.viewId).find(m => m.id === moduleId)
}

function getComponent(moduleId: string) {
  return getModuleDef(moduleId)?.component
}

// Called by grid-layout-plus after a drag or resize completes — sync back to store.
// We raise suppressRebuild so the watcher doesn't replace the entire gridLayout
// array while grid-layout-plus is still cleaning up its placeholder element.
// The flag is cleared on the next tick, after grid-layout-plus has finished.
function selectModule(moduleId: string) {
  layoutStore.addModule(props.viewId, moduleId)
  addModuleOpen.value = false
}

function onLayoutUpdated(newLayout: LayoutItem[]) {
  suppressRebuild = true
  layoutStore.handleLayoutUpdate(
    props.viewId,
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
</script>

<template>
  <div class="flex flex-col">

    <!-- ── Page header ──────────────────────────────────────────────────── -->
    <div class="flex items-center justify-between px-6 py-4 shrink-0">
      <h1 class="text-xl font-semibold">
        <slot name="title" />
      </h1>
      <div class="flex items-center gap-2">
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
          @click="layoutStore.toggleEditMode()"
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
            @remove="layoutStore.removeModule(viewId, String(item.i))"
          >
            <component
              :is="getComponent(String(item.i))"
              :view-id="viewId"
              :module-id="String(item.i)"
            />
          </ModuleCard>
        </GridItem>
      </GridLayout>
    </div>

  </div>
</template>
