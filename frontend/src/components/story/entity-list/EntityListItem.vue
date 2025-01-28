<script setup lang="ts">
const props = defineProps<{
  entityName: string
}>();

</script>

<template>
  <div class="group bg-muted/30 hover:bg-muted/40 rounded-lg hover:cursor-pointer py-2 flex flex-row items-center
              justify-center relative" >
    <p class="px-4 text-center">
      {{ props.entityName }}
    </p>

    <slot/>
  </div>
</template>