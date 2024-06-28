<script setup lang="ts">

import IconButton from "@/components/ui/button/IconButton.vue";
import {ChevronDown, ChevronUp} from "lucide-vue-next";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {useVModel} from "@vueuse/core";

const props = defineProps<{
  showCardBody: boolean,
}>()

const emits = defineEmits<{
  (e: 'update:showCardBody', payload: string | number): void,
  (e: 'toggle'): void
}>()

const showCardBody = useVModel(props, 'showCardBody', emits, {
  passive: true,
  defaultValue: true,
})
</script>

<template>
  <TextTooltip text="Minimize" v-if="showCardBody">
    <IconButton
        @click="emits('toggle')">
      <ChevronUp />
    </IconButton>
  </TextTooltip>

  <TextTooltip text="Expand" v-if="!showCardBody">
    <IconButton
        @click="emits('toggle')">
      <ChevronDown />
    </IconButton>
  </TextTooltip>
</template>

<style scoped>

</style>