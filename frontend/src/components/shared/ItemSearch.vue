<script setup lang="ts">

import {Input} from "@/components/ui/input";
import {Search} from "lucide-vue-next";
import {useVModel} from "@vueuse/core";

const props = defineProps<{
  searchInput: string,
  placeholder: string
}>()

const emits = defineEmits<{
  (e: 'update:searchInput', payload: string | number): void
}>()

const searchInput = useVModel(props, 'searchInput', emits, {
  passive: true,
  defaultValue: '',
})

</script>

<template>
  <div class="relative">
    <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground"/>
    <Input
        type="search"
        :placeholder="placeholder"
        class="pl-8"
        v-model="searchInput"
    />
  </div>
</template>

<style scoped>

</style>