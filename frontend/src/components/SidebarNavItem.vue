<script setup lang="ts">
import type { Component } from 'vue'

const props = defineProps<{
  icon: Component
  label: string
  active?: boolean
  disabled?: boolean
}>()

const emit = defineEmits<{
  select: []
}>()
</script>

<template>
  <button
    :disabled="props.disabled"
    :class="[
      'flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors text-left',
      props.active
        ? 'bg-sidebar-accent text-sidebar-accent-foreground'
        : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
      props.disabled && 'cursor-not-allowed opacity-40 hover:bg-transparent hover:text-sidebar-foreground',
    ]"
    @click="!props.disabled && emit('select')"
  >
    <component :is="props.icon" class="h-4 w-4 shrink-0" />
    <span>{{ props.label }}</span>
  </button>
</template>
