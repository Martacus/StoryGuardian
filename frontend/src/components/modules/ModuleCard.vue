<script setup lang="ts">
import { computed, provide } from 'vue'
import { GripVertical, Eye, EyeOff } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/models'
import { Card, CardHeader, CardTitle, CardContent, CardAction } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const props = defineProps<{
  viewId: string
  module: ModuleLayout
}>()

const layoutStore = useLayoutStore()

const actionsTargetId = computed(() => `module-actions-${props.module.id}`)
provide('moduleActionsTarget', actionsTargetId)

const definition = computed(() => moduleRegistry.find(m => m.id === props.module.id))
const label = computed(() => definition.value?.label ?? props.module.id)
const icon = computed(() => definition.value?.icon)

function toggleVisibility() {
  layoutStore.toggleModuleVisibility(props.viewId, props.module.id)
}
</script>

<template>
  <Card :class="[
    'h-full flex flex-col transition-opacity duration-200 overflow-hidden',
    !module.visible && layoutStore.editMode && 'opacity-40',
  ]">
    <CardHeader>
      <!-- Title row: drag handle (edit mode) + icon + label -->
      <CardTitle class="text-sm font-medium">
        <div class="flex items-center gap-2">
          <GripVertical
            v-if="layoutStore.editMode"
            class="drag-handle h-4 w-4 shrink-0 cursor-grab text-muted-foreground active:cursor-grabbing"
          />
          <component
            :is="icon"
            v-if="icon"
            class="h-4 w-4 shrink-0 text-muted-foreground"
          />
          {{ label }}
        </div>
      </CardTitle>

      <!-- Single CardAction: target div always in DOM (v-show) so Teleport never loses its anchor -->
      <CardAction>
        <!-- Teleport target: hidden in edit mode so actions don't show while dragging -->
        <div
          v-show="!layoutStore.editMode"
          :id="actionsTargetId"
          class="flex items-center gap-1 no-drag"
        />
        <!-- Visibility toggle: edit mode only -->
        <Button
          v-if="layoutStore.editMode"
          variant="ghost"
          size="icon-sm"
          class="no-drag"
          :title="module.visible ? 'Hide module' : 'Show module'"
          @click="toggleVisibility"
        >
          <Eye v-if="module.visible" class="h-4 w-4" />
          <EyeOff v-else class="h-4 w-4" />
        </Button>
      </CardAction>
    </CardHeader>

    <!-- Module content — fills remaining card height and scrolls if needed -->
    <CardContent class="flex-1 overflow-auto min-h-0">
      <slot />
    </CardContent>
  </Card>
</template>
