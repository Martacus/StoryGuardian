<script setup lang="ts">

import DashboardLayout from "@/layouts/DashboardLayout.vue";
import TextToolTip from "@/components/ui/tooltip/TextTooltip.vue";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import {Card, CardHeader} from "@/components/ui/card";
import {Plus, Settings, ArrowLeft} from "lucide-vue-next";
import {Button} from "@/components/ui/button";
import {onMounted, ref} from "vue";
import {useRoute, useRouter} from "vue-router";
import {
  AddEntityModule, EditEntityModuleConfig,
  GetEntity, GetEntityModules, RefreshEntity,
  SetEntityDescription,
  SetEntityName
} from "../../bindings/storyguardian/src/project/entitymanager";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import {Entity} from "../../bindings/storyguardian/src/project";
import Description from "@/components/shared/Description.vue";
import {useToast} from "@/components/ui/toast";
import RelationModule from "@/components/entity/RelationModule.vue";
import ModuleSelectItem from "@/components/story/modules/ModuleSelectItem.vue";

const route = useRoute();
const router = useRouter()
const entity = ref<Entity>();
const addModuleDialogOpened = ref(false)
const unusedModules = ref<string[]>([])
const {toast} = useToast()
const entityId: string = route.params['id'] as string

onMounted(async () => {
  await retrieveEntity();
});

const isUnused = (moduleName: string) => {
  return unusedModules.value.includes(moduleName);
};

const isUsedModule = (moduleName: string) => {
  if(entity.value){
    return moduleName in entity.value?.modules;
  }
  return false;
};

async function retrieveEntity(){
  let projectId: string = route.params['id'] as string
  const retrievedEntity = await GetEntity(projectId)
  if (retrievedEntity) {
    entity.value = retrievedEntity;
  }
}

async function saveStoryTitle(title: string) {
  if (!entity.value) return;
  try {
    await SetEntityName(entity.value.id, title)
  } catch (error: any) {
    toast({
      title: 'Failed to save entity title',
      description: error,
    });
  }
}

async function saveDescription(descriptionValue: string) {
  if (!entity.value) return;

  try {
    entity.value.description = await SetEntityDescription(entity.value.id, descriptionValue);
  } catch (error: any) {
    toast({
      title: 'Failed to save entity description',
      description: error,
    });
  }
}

function addEntityModule(module: string){
  AddEntityModule(entityId, module).then(() => {
    RefreshEntity(entityId).then(newEntity => {
      if(newEntity){
        entity.value = newEntity;
      }
    }).catch((error: string) => {
      toast({
        title: 'Failed to refresh entity',
        description: error,
      });
    });
  }).catch((error: string) => {
    toast({
      title: 'Failed to add module',
      description: error,
    });
  });
  addModuleDialogOpened.value = false;
}

function refreshUnusedEntityModules(){
  GetEntityModules(entityId, true).then((unusedModulesList: string[]) => {
    unusedModules.value = unusedModulesList;
  })
}

function moduleConfigChange(module: string, key: string, value: string) {
  EditEntityModuleConfig(entityId, module, key, value).catch((error: string) => {
    toast({
      title: 'Failed to save module config change',
      description: error,
    });
  });
}
</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 col-span-4" v-if="entity">
      <CardHeader class="flex flex-row items-center justify-between py-4">
        <Button class="btn btn-secondary" variant="outline" size="icon" @click="router.back()">
          <ArrowLeft />
        </Button>
        <EntityTitle :title="entity.name" @save-title="saveStoryTitle" class="flex flex-1 justify-center"/>
        <div class="flex flex-row gap-2">
          <Dialog v-model:open="addModuleDialogOpened">
            <DialogTrigger>
              <TextToolTip text="Add a module">
                <Button class="btn btn-secondary" variant="outline" size="icon" @click="refreshUnusedEntityModules">
                  <Plus/>
                </Button>
              </TextToolTip>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Select a module</DialogTitle>
              </DialogHeader>
              <DialogDescription>Choose a module to add to your entity, you can always remove them.</DialogDescription>
              <ModuleSelectItem v-if="isUnused('tagList')" @click="addEntityModule('tagList')">
                <p>Tags</p>
              </ModuleSelectItem>
            </DialogContent>
          </Dialog>
          <TextToolTip text="Story settings">
            <Button class="btn btn-secondary" variant="outline" size="icon">
              <Settings/>
            </Button>
          </TextToolTip>
        </div>
      </CardHeader>
    </Card>

    <Description
        v-if="entity"
        :module-config="entity.modules['description']"
        :description="entity.description"
        @save-description="saveDescription"
        @config-change="moduleConfigChange"
    />
    <RelationModule
        v-if="entity"
        :entity="entity"
        :module-config="entity.modules['relations']"
        @config-change="moduleConfigChange"
    />
    <div v-if="entity && isUsedModule('tagList')">
      Tags
    </div>
  </DashboardLayout>
</template>

<style scoped>

</style>