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
import { Plus, Settings } from 'lucide-vue-next';
import {onMounted, ref} from "vue";
import {Button} from "@/components/ui/button";
import {GetProject} from "../../bindings/storyguardian/project/projectmanager";
import {useRoute} from "vue-router";
import {Project} from "../../bindings/storyguardian/project";

const addModuleDialogOpened = ref(false);
const story = ref<Project>();
const route = useRoute();

onMounted(() => {
  let projectId: string = route.params['id'] as string
  console.log(projectId)
  GetProject(projectId).then(p => {
    if(p !== null){
      story.value = p
    } else {
      //Alert that project couldnt load
    }
  })
})
</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 flex flex-row p-2 justify-center gap-2">
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
    </Card>
<!--    <Description v-if="story" :description="story.description" @save-description="saveStoryDescription"/>-->
<!--    <EntityList v-if="story" :story="story"/>-->
<!--    <StoryImageModule v-if="story" :story="story"/>-->
  </DashboardLayout>
</template>

<style scoped>

</style>