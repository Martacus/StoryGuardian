<script setup lang="ts">
import DashboardLayout from "@/layouts/DashboardLayout.vue";
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
import Description from "@/components/shared/Description.vue";
import EntityList from "@/components/story/EntityList.vue";
import {useToast} from "@/components/ui/toast";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import ImageModule from "@/components/story/ImageModule.vue";
import {Story} from "../../bindings/storyguardian/src/project";
import {
  AddStoryModule,
  EditStoryModuleConfig, GetOpenStory,
  GetStory, GetStoryModules,
  SetStoryDescription,
  SetStoryTitle
} from "../../bindings/storyguardian/src/project/storymanager";
import TagList from "@/components/story/TagList.vue";
import ModuleSelectItem from "@/components/story/modules/ModuleSelectItem.vue";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import {LoadRelations} from "../../bindings/storyguardian/src/project/relationmanager";

const {toast} = useToast()
const addModuleDialogOpened = ref(false);
const story = ref<Story>();
const unusedStoryModules = ref<string[]>([])

const isUnused = (moduleName: string) => {
  return unusedStoryModules.value.includes(moduleName);
};

const isUsedModule = (moduleName: string) => {
  if(story.value){
    return moduleName in story.value?.modules;
  }
  return false;
};

onMounted(async () => {
  try {
  const retrievedStory = await GetOpenStory();
  LoadRelations()
  if (retrievedStory !== null) {
    story.value = retrievedStory
  }
} catch (error: any) {
  toast({
    title: 'Failed init story',
    description: error,
  });
}

})

async function retrieveStory(refresh: boolean){
  if(!story.value) return;

  try {
    const retrievedStory = await GetStory(story.value?.id, refresh)
    if (retrievedStory !== null) {
      story.value = retrievedStory
    }
  } catch (error: any) {
    toast({
      title: 'Failed to retrieve story',
      description: error,
    });
  }

  console.log(story.value?.tags)
}

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

function moduleConfigChange(module: string, key: string, value: string) {
  EditStoryModuleConfig(module, key, value).catch((error: string) => {
    toast({
      title: 'Failed to save module config change',
      description: error,
    });
  });
}

function refreshUnusedStoryModules(){
  GetStoryModules(true).then((unusedModules: string[]) => {
    console.log(unusedModules)
    unusedStoryModules.value = unusedModules;
  })
}

function addStoryModule(module: string){
  AddStoryModule(module).then(() => {
    retrieveStory(true);
  }).catch((error: string) => {
    toast({
      title: 'Failed to add module',
      description: error,
    });
  });
  addModuleDialogOpened.value = false;
}
</script>

<template>
  <DashboardLayout>
    <PageHeaderCard v-if="story">
      <EntityTitle :title="story.name" @save-title="saveStoryTitle" class="flex flex-1 justify-center"/>
      <div class="flex flex-row justify-end mr-2 gap-2">
        <Dialog v-model:open="addModuleDialogOpened">
          <DialogTrigger>
            <TextToolTip text="Add a module">
              <IconButton @click="refreshUnusedStoryModules">
                <Plus/>
              </IconButton>
            </TextToolTip>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Select a module</DialogTitle>
            </DialogHeader>
            <DialogDescription>Choose a module to add to your story, you can always remove them later.
            </DialogDescription>
            <ModuleSelectItem v-if="isUnused('description')">
              <p>Description</p>
            </ModuleSelectItem>
            <ModuleSelectItem v-if="isUnused('images')" @click="addStoryModule('images')">
              <p>Images</p>
            </ModuleSelectItem>
            <ModuleSelectItem v-if="isUnused('tagList')" @click="addStoryModule('tagList')">
              <p>Tags</p>
            </ModuleSelectItem>
          </DialogContent>
        </Dialog>
        <TextToolTip text="Story settings">
          <IconButton>
            <Settings/>
          </IconButton>
        </TextToolTip>
      </div>
    </PageHeaderCard>

    <Description
        v-if="story"
        :module-config="story.modules['description']"
        :description="story.description"
        @save-description="saveStoryDescription"
        @config-change="moduleConfigChange"
    />
    <EntityList
        v-if="story"
        :module-config="story.modules['entityList']"
        :story="story"
        @config-change="moduleConfigChange"
    />
    <TagList
        v-if="story && isUsedModule('tagList')"
        :module-config="story.modules['tagList']"
        :tags="story.tags"
        @config-change="moduleConfigChange"
    />
    <ImageModule
        v-if="story && isUsedModule('images')"
        :story="story"
        :module-config="story.modules['images']"
        @config-change="moduleConfigChange"
    />
  </DashboardLayout>
</template>

<style scoped>

</style>