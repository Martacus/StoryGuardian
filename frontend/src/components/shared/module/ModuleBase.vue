<script setup lang="ts">
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {useToggleBody} from "@/composables/useToggleBody";
import {useGridSize} from "@/composables/useGridSize";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import {useItemGridLayout} from "@/composables/useItemGridLayout";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";
import ItemViewSelector from "@/components/shared/button/ItemViewSelector.vue";
import {StoryModule} from "../../../../bindings/storyguardian/src/project";
import {watch} from "vue";

const props = withDefaults(defineProps<{
  title: string;
  moduleConfig: StoryModule;
  bodyToggle?: boolean;
  gridSize?: boolean;
  itemGridLayout?: boolean;
}>(), {
  bodyToggle: true,
  gridSize: true,
  itemGridLayout: false
});

const emit = defineEmits(['configChange', 'update:columnSize', 'update:showCardBody', 'update:itemView']);

const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig);
const {columnSize, changeGridSize} = useGridSize(props.moduleConfig);
const {itemView, changeItemView} = useItemGridLayout(props.moduleConfig);

watch(showCardBody, (newVal) => {
  emit('update:showCardBody', newVal);
});
watch(columnSize, (newVal) => {
  emit('update:columnSize', newVal);
});
watch(itemView, (newVal) => {
  emit('update:itemView', newVal);
});

</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>{{title}}</CardTitle>
      <slot name="center-space" v-if="showCardBody"/>
      <div class="flex flex-row space-x-2">
        <slot name="side-buttons" v-if="showCardBody"></slot>
        <VerticalSeperator/>
        <GridSizeSelector
            v-if="gridSize"
            :column-size="moduleConfig.configuration['columnSize']"
            @update-grid-size="(newSize) => changeGridSize(moduleConfig.name, newSize, emit)"
        />
        <ItemViewSelector
            v-if="itemGridLayout"
            :item-view="itemView"
            :show-card-body="showCardBody"
            @toggle="(payload) => changeItemView(moduleConfig.name, payload, emit)"
        />
        <CardBodyToggler
            v-if="bodyToggle"
            :show-card-body="showCardBody"
            @toggle="toggleCardBody(moduleConfig.name, emit)"
        />
      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <slot name="card-content"/>
    </CardContent>
  </Card>
</template>