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
import {onMounted, ref, watch} from "vue";
import Description from "@/components/shared/Description.vue";
import EntityList from "@/components/story/entity-list/EntityList.vue";
import {useToast} from "@/components/ui/toast";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import ImageModule from "@/components/story/ImageModule.vue";
import {Story, StoryModule} from "../../bindings/storyguardian/src/project";
import {
  AddStoryModule,
  EditStoryModuleConfig,
  GetOpenStory,
  GetStory,
  GetStoryModules,
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

const isUnused = (moduleName: string) => {
  if(!story.value) return false;
  let unused = true;
  story.value.modules.forEach((module) => {
    if (module.name === moduleName) {
      unused = false;
    }
  });
  return unused;
};

onMounted(async () => {
  try {
    const retrievedStory = await GetOpenStory();
    LoadRelations()

    if (retrievedStory !== null) {
      story.value = retrievedStory
    }

    //Sort modules by position
    if (story.value) {
      story.value.modules.sort((a, b) => {
        if (a.configuration["position"] < b.configuration["position"]) {
          return -1;
        }
        if (a.configuration["position"] > b.configuration["position"]) {
          return 1;
        }
        return 0;
      });
    }
  } catch (error: any) {
    toast({
      title: 'Failed init story',
      description: error,
    });
  }
});

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

function moduleConfigChange({module, key, value}: { module: string, key: string, value: string }) {
  EditStoryModuleConfig(module, key, value).catch((error: string) => {
    toast({
      title: 'Failed to save module config change',
      description: error,
    });
  });
}

function addStoryModule(module: string) {
  AddStoryModule(module).then((module) => {
    story.value?.modules.push(module);
  }).catch((error: string) => {
    toast({
      title: 'Failed to add module',
      description: error,
    });
  });
  addModuleDialogOpened.value = false;
}

function updateTagsFromTagList(tags: string[]) {
  if (story.value) {
    story.value.tags = tags;
  }
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
              <IconButton>
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

    <template v-if="story" v-for="module in story.modules">
      <Description
          v-if="module && module.name === 'description'"
          :module-config="module"
          :description="story.description"
          @save-description="saveStoryDescription"
          @config-change="moduleConfigChange"
      />
      <EntityList
          v-if="module && module.name === 'entities'"
          :module-config="module"
          :story="story"
          @config-change="moduleConfigChange"
      />
      <TagList
          v-if="module && module.name === 'tagList'"
          :module-config="module"
          :tags="story.tags"
          @config-change="moduleConfigChange"
          @update:tags="updateTagsFromTagList"
      />
      <ImageModule
          v-if="module && module.name === 'images'"
          :story="story"
          :module-config="module"
          @config-change="moduleConfigChange"
      />
    </template>
  </DashboardLayout>
</template>

<style scoped>

</style>