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
import StoryTitle from "@/components/story/StoryTitle.vue";
import Description from "@/components/shared/Description.vue";
import EntityList from "@/components/story/EntityList.vue";
import {Story} from "../../bindings/storyguardian/project";
import {GetStory, SetStoryDescription} from "../../bindings/storyguardian/project/storymanager";

const addModuleDialogOpened = ref(false);
const story = ref<Story>();
const route = useRoute();

onMounted(() => {
  let projectId: string = route.params['id'] as string
  console.log(projectId)
  GetStory(projectId).then(p => {
    if(p !== null){
      story.value = p
    } else {
      //Alert that project couldnt load
    }
  })
})

async function saveStoryDescription(descriptionValue: string) {
  if (!story.value) return;

  try {
    story.value.description = await SetStoryDescription(story.value.id, descriptionValue);
  } catch (error) {
    //Alert
    console.error('Failed to save story description:', error);
  }
}
</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 flex flex-row p-2 gap-2 items-center" v-if="story">
      <StoryTitle :story="story" class="w-full"/>
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
    <Description v-if="story" :description="story.description" @save-description="saveStoryDescription"/>
    <EntityList v-if="story" :story="story"/>
<!--    <StoryImageModule v-if="story" :story="story"/>-->
  </DashboardLayout>
</template>

<style scoped>

</style>