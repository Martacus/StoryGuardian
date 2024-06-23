<script setup lang="ts">

import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Square, Columns2, Columns3, Columns4} from "lucide-vue-next";
import {Button} from "@/components/ui/button";
import {onMounted, ref} from "vue";

const props = defineProps({
  columnSize: String
})


const emit = defineEmits(['updateGridSize']);
const localColumnSize = ref('4');

onMounted(()=>{
  if(props.columnSize){
    localColumnSize.value = props.columnSize + '';
  }
})

function updateGridSize(newSize: string) {
  localColumnSize.value = newSize;
  emit('updateGridSize', newSize);
}

</script>

<template>
  <TextTooltip text="Switch to grid" v-if="localColumnSize === '1'">
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('2')">
      <Square/>
    </Button>
  </TextTooltip>

  <TextTooltip text="Switch to list" v-if="localColumnSize === '2'">
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('3')">
      <Columns2/>
    </Button>
  </TextTooltip>

  <TextTooltip text="Switch to list" v-if="localColumnSize === '3'">
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('4')">
      <Columns3/>
    </Button>
  </TextTooltip>

  <TextTooltip text="Switch to list" v-if="localColumnSize === '4'">
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('1')">
      <Columns4/>
    </Button>
  </TextTooltip>
</template>

<style scoped>

</style>