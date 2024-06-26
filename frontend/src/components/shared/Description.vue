<script setup lang="ts">

import { Edit } from 'lucide-vue-next';
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import TipTap from "@/components/shared/TipTap.vue";
import {ref} from "vue";
import {Button} from "@/components/ui/button";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import {useGridSize} from "@/composables/useGridSize";
import {useToggleBody} from "@/composables/useToggleBody";
import VerticalSeparator from "@/components/ui/separator/VerticalSeparator.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";

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
          <IconButton
                  @click="toggleEdit()">
            <Edit />
          </IconButton>
        </TextTooltip>
        <VerticalSeparator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('description', newSize, emit)"/>
        <CardBodyToggler :show-card-body="showCardBody"  @toggle="toggleCardBody('description', emit)"/>
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