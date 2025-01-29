<script setup lang="ts">

import {ScrollArea} from "@/components/ui/scroll-area";
import {computed} from "vue";

const props = defineProps<{
  listHeight: string,
  itemView: string
}>();

const computedHeight = computed(() => props.listHeight);
const computedItemView = computed(() => props.itemView);
</script>

<template>
  <ScrollArea class="w-full" :class="computedHeight">
    <div id="single-entity-list" :class="[computedItemView === 'list' ? 'flex flex-col gap-2' : 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2']" ref="contentRef">
      <slot name="items"/>
    </div>
    <slot name="no-items"/>
  </ScrollArea>
</template>

<style scoped>

</style>