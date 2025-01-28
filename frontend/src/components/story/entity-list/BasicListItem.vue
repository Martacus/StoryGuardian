<script setup lang="ts">
const props = defineProps<{
  text: string
}>();

</script>

<template>
  <div class="group bg-muted/30 hover:bg-muted/40 rounded-lg hover:cursor-pointer py-2 flex flex-row items-center
              justify-center relative" >
    <p class="px-4 text-center">
      {{ props.text }}
    </p>

    <slot name="item-action"/>
  </div>
</template>