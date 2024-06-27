<script setup lang="ts">

import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Square, Columns2, Columns3, Columns4} from "lucide-vue-next";
import {onMounted, ref} from "vue";
import {Popover, PopoverTrigger, PopoverContent} from "@/components/ui/popover";
import IconButton from "@/components/ui/button/IconButton.vue";

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
        <IconButton>
          <Columns2/>
        </IconButton>
      </PopoverTrigger>
      <PopoverContent side="top" align="center" class="flex flex-1 gap-2 justify-center">
        <IconButton @click="updateGridSize('1')" >
          <Square/>
        </IconButton>
        <IconButton @click="updateGridSize('2')" >
          <Columns2/>
        </IconButton>
        <IconButton @click="updateGridSize('3')" >
          <Columns3/>
        </IconButton>
        <IconButton @click="updateGridSize('4')" >
          <Columns4/>
        </IconButton>
      </PopoverContent>
    </Popover>

  </TextTooltip>
</template>

<style scoped>

</style>