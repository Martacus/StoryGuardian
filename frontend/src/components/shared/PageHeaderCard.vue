<script setup lang="ts">

import {Card} from "@/components/ui/card";
</script>

<template>
  <Card class="bg-muted/30 flex flex-row p-2 gap-2 items-center justify-between col-span-1 sm:col-span-1 md:col-span-2 lg:col-span-3 xl:col-span-4">
    <slot />
  </Card>
</template>

<style scoped>

</style>