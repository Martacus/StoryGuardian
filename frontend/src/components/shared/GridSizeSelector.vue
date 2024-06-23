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
  <TextTooltip text="Switch column width">
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('2')" v-if="localColumnSize === '1'">
      <Square/>
    </Button>
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('3')"  v-if="localColumnSize === '2'">
      <Columns2/>
    </Button>
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('4')" v-if="localColumnSize === '3'">
      <Columns3/>
    </Button>
    <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('1')" v-if="localColumnSize === '4'">
      <Columns4/>
    </Button>
  </TextTooltip>
</template>

<style scoped>

</style>