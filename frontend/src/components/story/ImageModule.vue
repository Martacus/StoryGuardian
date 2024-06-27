<script setup lang="ts">
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '@/components/ui/table'
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";

import {ChevronDown, ChevronUp, RefreshCcw} from 'lucide-vue-next';
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {useToast} from "@/components/ui/toast";
import {Dialog, DialogContent} from "@/components/ui/dialog";
import {ImageFile, Story, StoryModule} from "../../../bindings/storyguardian/src/project";
import {GetStoryImages} from "../../../bindings/storyguardian/src/project/storymanager";
import {OpenProjectFolder} from "../../../bindings/storyguardian/src/project/applicationmanager";
import {useToggleBody} from "@/composables/useToggleBody";
import {useGridSize} from "@/composables/useGridSize";
import GridSizeSelector from "@/components/shared/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeperator.vue";
import IconButton from "@/components/ui/button/IconButton.vue";

const props = defineProps<{
  story: Story
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast()
const images = ref<ImageFile[]>([])
const openDialogs = ref(images.value.map(() => false));


const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig);
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig);


onMounted(async () => {
  if (!props.story) {
    return
  }

  await retrieveImages()
});

async function retrieveImages(){
  try {
    images.value = await GetStoryImages();
    openDialogs.value = images.value.map(() => false);
  } catch (error: any) {
    toast({
      title: 'error loading images',
      description: error
    })
  }
}

function openImageFolder(){
  OpenProjectFolder('images');
}
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Images</CardTitle>
      <div class="flex flex-row space-x-2">
        <Button class="btn btn-primary" variant="outline" @click="openImageFolder()">
          Open folder
        </Button>
        <IconButton @click="retrieveImages()">
          <RefreshCcw/>
        </IconButton>
        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('images', newSize, emit)"/>
        <TextTooltip text="Expand" v-if="!showCardBody">
          <IconButton @click="toggleCardBody('images', emit)">
            <ChevronDown/>
          </IconButton>
        </TextTooltip>

        <TextTooltip text="Minimize" v-if="showCardBody">
          <IconButton @click="toggleCardBody('images', emit)">
            <ChevronUp/>
          </IconButton>
        </TextTooltip>

      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead class="w-[100px]">
              Name
            </TableHead>
            <TableHead>Location</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(image, index) in images" @click="openDialogs[index] = true" class="hover:cursor-pointer">

            <!-- Table Data -->
            <TableCell class="font-medium">
              {{ image.name }}
            </TableCell>
            <TableCell>{{ image.location }}</TableCell>

            <!--  Image Dialog  -->
            <Dialog v-model:open="openDialogs[index]" v-if="showCardBody">
              <DialogContent>
                <img :src="'/images/' + image.location" alt="Error loading image"/>
              </DialogContent>
            </Dialog>

          </TableRow>
        </TableBody>
      </Table>
    </CardContent>
  </Card>
</template>

<style scoped>

</style>