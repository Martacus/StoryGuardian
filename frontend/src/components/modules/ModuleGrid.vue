<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { GridLayout, GridItem } from 'grid-layout-plus'
import type { LayoutItem } from 'grid-layout-plus'
import { Settings2, Check } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/models'
import ModuleCard from './ModuleCard.vue'
import { Button } from '@/components/ui/button'

const props = defineProps<{
  viewId: string
}>()

const layoutStore = useLayoutStore()

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
          >
            <component :is="getComponent(String(item.i))" />
          </ModuleCard>
        </GridItem>
      </GridLayout>
    </div>

  </div>
</template>
