<!-- frontend/src/components/entity/EntityTagList.vue -->
<script setup lang="ts">
import TagListBase from "@/components/shared/TagListBase.vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])
</script>

<template>
  <TagListBase :tags="props.tags" :moduleConfig="props.moduleConfig" @configChange="emit('configChange', $event)">
    <template #dialog-content>
      <!-- Custom Dialog content for EntityTagList -->
      <p>EntityTagList specific content</p>
    </template>
  </TagListBase>
</template>