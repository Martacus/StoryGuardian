<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
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

// Local copy of the module list, used as VueDraggable's v-model.
// Synced from the store whenever the layout changes (world load, colSpan, visibility).
const modules = ref<ModuleLayout[]>([])

watch(
  () => layoutStore.layout,
  () => { modules.value = layoutStore.getViewModules(props.viewId) },
  { immediate: true, deep: true },
)

// View mode only renders visible modules
const visibleModules = computed(() => modules.value.filter(m => m.visible))

// Called after a drag-and-drop reorder — VueDraggable has already updated modules.value
function onDragEnd() {
  layoutStore.updateModuleOrder(props.viewId, [...modules.value])
}

function getComponent(moduleId: string) {
  return moduleRegistry.find(m => m.id === moduleId)?.component
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

    <!-- ── Edit mode: draggable grid (all modules shown, hidden ones dimmed) ── -->
    <VueDraggable
      v-if="layoutStore.editMode"
      v-model="modules"
      handle=".drag-handle"
      :animation="150"
      class="grid grid-cols-12 gap-4 px-6 pb-6 content-start"
      @end="onDragEnd"
    >
      <div
        v-for="mod in modules"
        :key="mod.id"
        :style="{ gridColumn: `span ${mod.colSpan}` }"
      >
        <ModuleCard :view-id="viewId" :module="mod">
          <component :is="getComponent(mod.id)" />
        </ModuleCard>
      </div>
    </VueDraggable>

    <!-- ── View mode: static grid (visible modules only) ─────────────────── -->
    <div
      v-else
      class="grid grid-cols-12 gap-4 px-6 pb-6 content-start"
    >
      <div
        v-for="mod in visibleModules"
        :key="mod.id"
        :style="{ gridColumn: `span ${mod.colSpan}` }"
      >
        <ModuleCard :view-id="viewId" :module="mod">
          <component :is="getComponent(mod.id)" />
        </ModuleCard>
      </div>
    </div>

  </div>
</template>
