<script setup lang="ts">

import { ChevronDown, ChevronUp, Edit } from 'lucide-vue-next';
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import TipTap from "@/components/shared/TipTap.vue";
import {ref} from "vue";
import {Button} from "@/components/ui/button";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import GridSizeSelector from "@/components/shared/GridSizeSelector.vue";
import {useGridSize} from "@/composables/useGridSize";
import {useToggleBody} from "@/composables/useToggleBody";
import VerticalSeperator from "@/components/ui/separator/VerticalSeperator.vue";

const props = defineProps<{
  description: String,
  moduleConfig: StoryModule
}>()

const emit = defineEmits(['saveDescription', 'configChange'])

const storyDescriptionEditor = ref();
const isEditing = ref(false);

const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig)

function toggleEdit() {
  isEditing.value = !isEditing.value;
}

async function save() {
  emit('saveDescription', storyDescriptionEditor.value?.getHTML());
  toggleEdit();
}
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Description</CardTitle>
      <div class="flex flex-row space-x-2">
        <TextTooltip text="Edit" v-if="showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleEdit()">
            <Edit />
          </Button>
        </TextTooltip>
        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('description', newSize, emit)"/>
        <TextTooltip text="Minimize" v-if="showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleCardBody('description', emit)">
            <ChevronUp />
          </Button>
        </TextTooltip>

        <TextTooltip text="Expand" v-if="!showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleCardBody('description', emit)">
            <ChevronDown />
          </Button>
        </TextTooltip>
      </div>

    </CardHeader>
    <CardContent v-if="showCardBody">
      <div v-html="description" v-if="!isEditing"></div>
      <div v-if="isEditing" class="flex flex-col items-center">
        <TipTap v-bind:content="description" ref="storyDescriptionEditor"></TipTap>
        <Button class="btn btn-primary mt-2 max-w-40" @click="save()">Save Description</Button>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>

</style>