<script setup lang="ts">
import { computed, ref } from 'vue'
import { Input } from '@/components/ui/input'

const props = defineProps<{
  id: string
  modelValue: string
  placeholder?: string
  options: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  enter: []
}>()

const open = ref(false)

const visibleOptions = computed(() => {
  const query = props.modelValue.trim().toLowerCase()
  return props.options.filter(option => {
    const normalized = option.toLowerCase()
    return normalized !== query && (!query || normalized.includes(query))
  })
})

function updateValue(value: string) {
  emit('update:modelValue', value)
  open.value = true
}

function selectOption(option: string) {
  emit('update:modelValue', option)
  open.value = false
}

function handleFocus() {
  open.value = visibleOptions.value.length > 0
}

function handleBlur() {
  open.value = false
}

function handleEnter() {
  open.value = false
  emit('enter')
}
</script>

<template>
  <div class="relative">
    <Input
      :id="id"
      :model-value="modelValue"
      :placeholder="placeholder"
      autocomplete="off"
      @update:model-value="value => updateValue(String(value))"
      @focus="handleFocus"
      @blur="handleBlur"
      @keydown.enter="handleEnter"
      @keydown.escape="open = false"
    />
    <div
      v-if="open && visibleOptions.length > 0"
      class="absolute z-50 mt-1 max-h-44 w-full overflow-auto rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
    >
      <button
        v-for="option in visibleOptions"
        :key="option"
        type="button"
        class="flex w-full items-center rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent hover:text-accent-foreground"
        @mousedown.prevent="selectOption(option)"
      >
        {{ option }}
      </button>
    </div>
  </div>
</template>
