<script setup lang="ts">

import {ref} from "vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import {Edit} from 'lucide-vue-next';
import ModuleBase from "@/components/shared/module/ModuleBase.vue";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import TipTap from "@/components/shared/TipTap.vue";
import {Button} from "@/components/ui/button";

defineProps<{
  description: String,
  moduleConfig: StoryModule
}>()

const emit = defineEmits(['saveDescription', 'configChange'])

const storyDescriptionEditor = ref();
const isEditing = ref(false);

function toggleEdit() {
  isEditing.value = !isEditing.value;
}

async function save() {
  emit('saveDescription', storyDescriptionEditor.value?.getHTML());
  toggleEdit();
}
</script>

<template>
  <ModuleBase title="Description" :module-config="moduleConfig" @config-change="(payload) => emit('configChange', payload)">
    <template #side-buttons>
      <TextTooltip text="Edit">
        <IconButton
            @click="toggleEdit()">
          <Edit />
        </IconButton>
      </TextTooltip>
    </template>
    <template #card-content>
      <div v-html="description" v-if="!isEditing"></div>
      <div v-if="isEditing" class="flex flex-col items-center">
        <TipTap v-bind:content="description" ref="storyDescriptionEditor"></TipTap>
        <Button class="btn btn-primary mt-2 max-w-40" @click="save()">Save Description</Button>
      </div>
    </template>
  </ModuleBase>
</template>

<style scoped>

</style>