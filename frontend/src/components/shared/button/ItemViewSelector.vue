<script setup lang="ts">

import IconButton from "@/components/ui/button/IconButton.vue";
import {LayoutGrid, StretchHorizontal} from "lucide-vue-next";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";

const props = defineProps<{
  itemView: String,
  showCardBody: boolean
}>()

const emits = defineEmits<{
  (e: 'update:showCardBody', payload: string | number): void,
  (e: 'toggle', payload: string): void
}>()
</script>

<template>
  <TextTooltip text="Switch to grid" v-if="itemView === 'list' && showCardBody">
    <IconButton
        @click="emits('toggle', 'grid')" >
      <StretchHorizontal/>
    </IconButton>
  </TextTooltip>

  <TextTooltip text="Switch to list" v-if="itemView === 'grid' && showCardBody">
    <IconButton
        @click="emits('toggle', 'list')" >
      <LayoutGrid/>
    </IconButton>
  </TextTooltip>
</template>

<style scoped>

</style>