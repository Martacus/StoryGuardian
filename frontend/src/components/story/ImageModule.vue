<script setup lang="ts">
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '@/components/ui/table'
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";

import {ChevronDown, ChevronUp, RefreshCcw} from 'lucide-vue-next';
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {useToast} from "@/components/ui/toast";
import {Dialog, DialogContent} from "@/components/ui/dialog";
import {ImageFile, Story} from "../../../bindings/storyguardian/internal/project";
import {GetStoryImages} from "../../../bindings/storyguardian/internal/project/storymanager";
import {OpenProjectFolder} from "../../../bindings/storyguardian/internal/project/applicationmanager";

const showBody = ref(true);
const {toast} = useToast()
const images = ref<ImageFile[]>([])
const openDialogs = ref(images.value.map(() => false));
const imageDialog = ref(false);

const props = defineProps<{
  story: Story
}>();

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

function toggleCard() {
  showBody.value = !showBody.value;
}

function openImageFolder(){
  OpenProjectFolder('images');
}
</script>

<template>
  <Card class="bg-muted/30">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Images</CardTitle>
      <div class="flex flex-row space-x-2">
        <Button class="btn btn-primary" aria-label="Toggle italic" variant="outline" @click="openImageFolder()">
          Open folder
        </Button>
        <Button size="icon" aria-label="Toggle italic" variant="outline" @click="retrieveImages()">
          <RefreshCcw/>
        </Button>
        <TextTooltip text="Expand" v-if="!showBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
            <ChevronDown/>
          </Button>
        </TextTooltip>

        <TextTooltip text="Minimize" v-if="showBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
            <ChevronUp/>
          </Button>
        </TextTooltip>

      </div>
    </CardHeader>
    <CardContent v-if="showBody">
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
          <TableRow v-for="(image, index) in images" @click="openDialogs[index] = true">

            <!-- Table Data -->
            <TableCell class="font-medium">
              {{ image.name }}
            </TableCell>
            <TableCell>{{ image.location }}</TableCell>

            <!--  Image Dialog  -->
            <Dialog v-model:open="openDialogs[index]" v-if="showBody">
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