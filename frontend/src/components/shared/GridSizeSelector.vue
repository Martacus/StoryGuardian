<script setup lang="ts">

import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Square, Columns2, Columns3, Columns4} from "lucide-vue-next";
import {Button} from "@/components/ui/button";
import {onMounted, ref} from "vue";
import {Popover, PopoverTrigger, PopoverContent} from "@/components/ui/popover";

const props = defineProps<{
  columnSize: String
}>()


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
    <Popover>
      <PopoverTrigger>
        <Button size="icon" aria-label="Toggle italic" variant="outline">
          <Columns2/>
        </Button>
      </PopoverTrigger>
      <PopoverContent side="top" align="center" class="flex flex-1 gap-2 justify-center">
        <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('1')" >
          <Square/>
        </Button>
        <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('2')" >
          <Columns2/>
        </Button>
        <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('3')" >
          <Columns3/>
        </Button>
        <Button size="icon" aria-label="Toggle italic" variant="outline" @click="updateGridSize('4')" >
          <Columns4/>
        </Button>
      </PopoverContent>
    </Popover>

  </TextTooltip>
</template>

<style scoped>

</style>