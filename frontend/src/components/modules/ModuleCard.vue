<script setup lang="ts">
import { computed } from 'vue'
import { GripVertical, Eye, EyeOff, ChevronDown } from 'lucide-vue-next'
import { useLayoutStore } from '@/stores/layoutStore'
import { moduleRegistry } from '@/modules/registry'
import type { ModuleLayout } from '../../../bindings/litguardian/models'
import { Card, CardHeader, CardTitle, CardContent, CardAction } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu'

const props = defineProps<{
  viewId: string
  module: ModuleLayout
}>()

const layoutStore = useLayoutStore()

const definition = computed(() => moduleRegistry.find(m => m.id === props.module.id))
const label = computed(() => definition.value?.label ?? props.module.id)
const icon = computed(() => definition.value?.icon)

const colSpanOptions = [
  { label: 'Quarter',    value: 3  },
  { label: 'Third',      value: 4  },
  { label: 'Half',       value: 6  },
  { label: 'Two-thirds', value: 8  },
  { label: 'Full',       value: 12 },
]

const currentWidthLabel = computed(() =>
  colSpanOptions.find(o => o.value === props.module.colSpan)?.label ?? 'Width'
)

function setColSpan(span: number) {
  layoutStore.setModuleColSpan(props.viewId, props.module.id, span)
}

function toggleVisibility() {
  layoutStore.toggleModuleVisibility(props.viewId, props.module.id)
}
</script>

<template>
  <Card :class="['transition-opacity duration-200', !module.visible && layoutStore.editMode && 'opacity-40']">
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

      <!-- Edit-mode controls: resize dropdown + visibility toggle -->
      <CardAction v-if="layoutStore.editMode">
        <div class="flex items-center gap-1">
          <!-- Resize dropdown -->
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="ghost" size="sm" class="h-7 gap-1 px-2 text-xs">
                {{ currentWidthLabel }}
                <ChevronDown class="h-3 w-3" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                v-for="opt in colSpanOptions"
                :key="opt.value"
                @click="setColSpan(opt.value)"
              >
                {{ opt.label }}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <!-- Visibility toggle -->
          <Button
            variant="ghost"
            size="icon-sm"
            :title="module.visible ? 'Hide module' : 'Show module'"
            @click="toggleVisibility"
          >
            <Eye v-if="module.visible" class="h-4 w-4" />
            <EyeOff v-else class="h-4 w-4" />
          </Button>
        </div>
      </CardAction>
    </CardHeader>

    <!-- Module content injected by ModuleGrid -->
    <CardContent>
      <slot />
    </CardContent>
  </Card>
</template>
