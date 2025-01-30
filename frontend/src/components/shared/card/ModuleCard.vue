<script setup lang="ts">
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";

defineProps({
  title: String
})
</script>

<template>
  <Card class="bg-muted/30 w-full">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>{{title}}</CardTitle>
      <slot name="center-space"/>
      <div class="flex flex-row space-x-2">
        <slot name="side-buttons"></slot>
      </div>
    </CardHeader>
    <CardContent>
      <slot name="card-content"/>
    </CardContent>
  </Card>
</template>

<style scoped>

</style>