<script setup lang="ts">
import { computed, provide, shallowRef, type Component } from 'vue'
import { GripVertical, Eye, EyeOff, X } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/internal/models'
import { Card, CardHeader, CardTitle, CardContent, CardAction } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const props = defineProps<{
  viewId: string
  module: ModuleLayout
}>()

const emit = defineEmits<{
  remove: []
}>()

const layoutStore = useLayoutStore()

const moduleActions = shallowRef<Component | null>(null)
provide('moduleActions', moduleActions)

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

      <CardAction>
        <!-- Module-provided actions: hidden in edit mode -->
        <div v-if="moduleActions && !layoutStore.editMode" class="flex items-center gap-1 no-drag">
          <component :is="moduleActions" />
        </div>
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
        <!-- Remove button: edit mode only -->
        <Button
          v-if="layoutStore.editMode"
          variant="ghost"
          size="icon-sm"
          class="no-drag text-muted-foreground hover:text-destructive"
          title="Remove module"
          @click="emit('remove')"
        >
          <X class="h-4 w-4" />
        </Button>
      </CardAction>
    </CardHeader>

    <!-- Module content — fills remaining card height and scrolls if needed -->
    <CardContent class="flex-1 overflow-auto min-h-0">
      <slot />
    </CardContent>
  </Card>
</template>
