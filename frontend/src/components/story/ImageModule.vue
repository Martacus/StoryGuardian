<script setup lang="ts">
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow,} from '@/components/ui/table'
import {RefreshCcw} from 'lucide-vue-next';
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {useToast} from "@/components/ui/toast";
import {Dialog, DialogContent} from "@/components/ui/dialog";
import {ImageFile, Story, StoryModule} from "../../../bindings/storyguardian/src/project";
import {GetStoryImages} from "../../../bindings/storyguardian/src/project/storymanager";
import {OpenProjectFolder} from "../../../bindings/storyguardian/src/project/applicationmanager";
import IconButton from "@/components/ui/button/IconButton.vue";
import ModuleBase from "@/components/shared/module/ModuleBase.vue";

const props = defineProps<{
  story: Story
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast()
const images = ref<ImageFile[]>([])
const openDialogs = ref(images.value.map(() => false));

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
  <ModuleBase title="Images" :module-config="moduleConfig" @config-change="(payload) => emit('configChange', payload)">
    <template #side-buttons>
      <Button class="btn btn-primary" variant="outline" @click="openImageFolder()">
        Open folder
      </Button>
      <IconButton @click="retrieveImages()">
        <RefreshCcw/>
      </IconButton>
    </template>
    <template #card-content>
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
            <Dialog v-model:open="openDialogs[index]">
              <DialogContent>
                <img :src="'/images/' + image.location" alt="Error loading image"/>
              </DialogContent>
            </Dialog>

          </TableRow>
        </TableBody>
      </Table>
    </template>
  </ModuleBase>
</template>

<style scoped>

</style>