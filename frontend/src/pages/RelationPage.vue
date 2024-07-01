<script setup lang="ts">

import EntityTitle from "@/components/shared/EntityTitle.vue";
import {ArrowLeft, Plus, Settings} from "lucide-vue-next";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import {Relation} from "../../bindings/storyguardian/src/project";
import {useRoute, useRouter} from "vue-router";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from "@/components/ui/dialog";
import Description from "@/components/shared/Description.vue";
import TextToolTip from "@/components/ui/tooltip/TextTooltip.vue";
import {onMounted, ref} from "vue";
import {useToast} from "@/components/ui/toast";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import RelationInfo from "@/components/relation/RelationInfo.vue";
import {
  EditRelationModuleConfig,
  GetRelation,
  SetRelationDescription,
  SetRelationName
} from "../../bindings/storyguardian/src/project/relationmanager";

const router = useRouter();
const route = useRoute();
const {toast} = useToast();

const relation = ref<Relation>();
const relationId = ref('');
const addModuleDialogOpened = ref(false)

onMounted(async () => {
  let relationid: string = route.params['id'] as string
  let gotRelation = await GetRelation(relationid)
  if(gotRelation){
    relation.value = gotRelation
    relationId.value = gotRelation.id
  }
});

async function saveRelationTitle(title: string) {
  if (!relation.value) return;
  try {
    await SetRelationName(relation.value.id, title)
  } catch (error: any) {
    toast({
      title: 'Failed to save relation title',
      description: error,
    });
  }
}

async function saveDescription(descriptionValue: string) {
  if (!relation.value) return;

  try {
    relation.value.description = await SetRelationDescription(relation.value.id, descriptionValue);
  } catch (error: any) {
    toast({
      title: 'Failed to save relation description',
      description: error,
    });
  }
}

function moduleConfigChange(module: string, key: string, value: string) {
  EditRelationModuleConfig(relationId.value, module, key, value).catch((error: string) => {
    toast({
      title: 'Failed to save module config change',
      description: error,
    });
  });
}
</script>

<template>
  <DashboardLayout>
    <PageHeaderCard v-if="relation">
      <IconButton size="icon" @click="router.back()">
        <ArrowLeft/>
      </IconButton>
      <div class="flex flex-1 justify-center">
        <p class="text-2xl leading-loose ml-2">Relation - </p> <EntityTitle :title="relation.name" @save-title="saveRelationTitle" />
      </div>
      <div class="flex flex-row gap-2">
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
            <DialogDescription>Choose a module to add to your relation, you can always remove them.
            </DialogDescription>
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
        v-if="relation"
        :description="relation.description"
        @save-description="saveDescription"
        :module-config="relation.modules['description']"
        @config-change="moduleConfigChange"
    />
    <RelationInfo
        v-if="relation"
        :module-config="relation.modules['relationInfo']"
        @config-change="moduleConfigChange"
    />
  </DashboardLayout>
</template>

<style scoped>

</style>