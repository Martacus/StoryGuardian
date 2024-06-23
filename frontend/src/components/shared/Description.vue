<script setup lang="ts">

import { ChevronDown, ChevronUp, Edit } from 'lucide-vue-next';
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import TipTap from "@/components/shared/TipTap.vue";
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {StoryModule} from "../../../bindings/storyguardian/internal/project";
import GridSizeSelector from "@/components/shared/GridSizeSelector.vue";


const storyDescriptionEditor = ref();
const isEditing = ref(false);
const showCardBody = ref(true);
const columnSize = ref('col-span-4')

const emit = defineEmits(['saveDescription', 'configChange'])

const props = defineProps({
  description: String,
  moduleConfig: StoryModule
})

function toggleEdit() {
  isEditing.value = !isEditing.value;
}

function toggleCardBody() {
  showCardBody.value = !showCardBody.value;
}

async function save() {
  emit('saveDescription', storyDescriptionEditor.value?.getHTML());
  toggleEdit();
}

onMounted(()=> {
  columnSize.value = columnSize.value.slice(0, -1) + props.moduleConfig?.configuration['columnSize'];
})

function changeGridSize(newColumnSize: string){
  if(props.moduleConfig){
    props.moduleConfig.configuration['columnSize'] = newColumnSize;
    columnSize.value = columnSize.value.slice(0, -1) + newColumnSize;
  }
  emit('configChange', 'columnSize', newColumnSize)
}
</script>

<template>
  <Card class="bg-muted/30" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Description</CardTitle>
      <div class="flex flex-row space-x-2">
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="changeGridSize"/>
        <TextTooltip text="Edit" v-if="showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleEdit()">
            <Edit />
          </Button>
        </TextTooltip>

        <TextTooltip text="Minimize" v-if="showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleCardBody()">
            <ChevronUp />
          </Button>
        </TextTooltip>

        <TextTooltip text="Expand" v-if="!showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline"
                  @click="toggleCardBody()">
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