<script setup lang="ts">

import { ChevronDown, ChevronUp, Edit } from 'lucide-vue-next';
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import TipTap from "@/components/shared/TipTap.vue";
import {ref} from "vue";
import {Button} from "@/components/ui/button";


const storyDescriptionEditor = ref();
const editDescription = ref(false);
const showDescription = ref(true);

const emit = defineEmits(['saveDescription'])

defineProps({
  description: String
})

function activateEditDescription() {
  editDescription.value = !editDescription.value;
}

function toggleDescription() {
  showDescription.value = !showDescription.value;
}

async function save() {
  emit('saveDescription', storyDescriptionEditor.value?.getHTML());
  activateEditDescription();
}
</script>

<template>
  <Card class="bg-muted/30">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Description</CardTitle>
      <div class="flex flex-row space-x-2">
        <TextTooltip text="Edit">
          <Button size="icon" aria-label="Toggle italic" variant="outline" v-if="showDescription"
                  @click="activateEditDescription()">
            <Edit />
          </Button>
        </TextTooltip>

        <TextTooltip text="Minimize">
          <Button size="icon" aria-label="Toggle italic" variant="outline" v-if="showDescription"
                  @click="toggleDescription()">
            <ChevronUp />
          </Button>
        </TextTooltip>

        <TextTooltip text="Expand">
          <Button size="icon" aria-label="Toggle italic" variant="outline" v-if="!showDescription"
                  @click="toggleDescription()">
            <ChevronDown />
          </Button>
        </TextTooltip>
      </div>

    </CardHeader>
    <CardContent v-if="showDescription">
      <div v-html="description" v-if="!editDescription"></div>
      <div v-if="editDescription" class="flex flex-col items-center">
        <TipTap v-bind:content="description" ref="storyDescriptionEditor"></TipTap>
        <Button class="btn btn-primary mt-2 max-w-40" @click="save()">Save Description</Button>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>

</style>