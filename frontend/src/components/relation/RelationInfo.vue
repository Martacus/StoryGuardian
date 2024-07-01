<script setup lang="ts">

import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import {useGridSize} from "@/composables/useGridSize";
import {useToggleBody} from "@/composables/useToggleBody";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";

const props = defineProps<{
  moduleConfig: StoryModule
}>()

const emit = defineEmits(['configChange'])


const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig)

</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Information</CardTitle>
      <div class="flex flex-row space-x-2">
        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('description', newSize, emit)"/>
        <CardBodyToggler :show-card-body="showCardBody"  @toggle="toggleCardBody('relationInfo', emit)"/>
      </div>

    </CardHeader>
    <CardContent v-if="showCardBody">

    </CardContent>
  </Card>
</template>

<style scoped>

</style>