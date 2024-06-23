<script setup lang="ts">

import DashboardLayout from "@/layouts/DashboardLayout.vue";
import {Card} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import TextToolTip from "@/components/ui/tooltip/TextTooltip.vue";
import {Plus, Settings} from 'lucide-vue-next';
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {useRoute} from "vue-router";
import Description from "@/components/shared/Description.vue";
import EntityList from "@/components/story/EntityList.vue";
import {useToast} from "@/components/ui/toast";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import ImageModule from "@/components/story/ImageModule.vue";
import {Story} from "../../bindings/storyguardian/internal/project";
import {GetStory, SetStoryDescription, SetStoryTitle} from "../../bindings/storyguardian/internal/project/storymanager";
import TagList from "@/components/story/TagList.vue";

const addModuleDialogOpened = ref(false);
const story = ref<Story>();
const route = useRoute();
const {toast} = useToast()

onMounted(async () => {
  let projectId: string = route.params['id'] as string
  try{
    const retrievedStory = await GetStory(projectId)
    if(retrievedStory !== null){
      story.value = retrievedStory
      console.log(retrievedStory.modules)
    }
  } catch(error: any) {
    toast({
      title: 'Failed to retrieve story',
      description: error,
    });
  }
})

async function saveStoryDescription(descriptionValue: string) {
  if (!story.value) return;

  try {
    story.value.description = await SetStoryDescription(descriptionValue);
  } catch (error: any) {
    toast({
      title: 'Failed to save story description',
      description: error,
    });
  }
}

async function saveStoryTitle(title: string) {
  if (!story.value) return;
  try {
    await SetStoryTitle(title)
  } catch (error: any) {
    toast({
      title: 'Failed to save story title',
      description: error,
    });
  }
}
</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 flex flex-row p-2 gap-2 items-center col-span-4" v-if="story">
      <EntityTitle :title="story.name" @save-title="saveStoryTitle" class="w-full"/>
      <div class="flex flex-row justify-end mr-2 gap-2">
        <Dialog v-model:open="addModuleDialogOpened">
          <DialogTrigger>
            <TextToolTip text="Add a module">
              <Button class="btn btn-secondary" variant="outline" size="icon">
                <Plus/>
              </Button>
            </TextToolTip>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Select a module</DialogTitle>
            </DialogHeader>
            <DialogDescription>Choose a module to add to your story, you can always remove them later.</DialogDescription>
          </DialogContent>
        </Dialog>
        <TextToolTip text="Story settings">
          <Button class="btn btn-secondary" variant="outline" size="icon">
            <Settings/>
          </Button>
        </TextToolTip>
      </div>
    </Card>

    <Description v-if="story" :module-config="story.modules['description']" :description="story.description" @save-description="saveStoryDescription"/>
    <EntityList v-if="story" :story="story" class="col-span-4"/>
    <TagList v-if="story"  :tags="story.tags"/>
    <ImageModule v-if="story" :story="story" class="col-span-4"/>
  </DashboardLayout>
</template>

<style scoped>

</style>