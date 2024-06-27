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
import {
  GetRelation,
  SetRelationDescription,
  SetRelationName
} from "../../bindings/storyguardian/src/project/entitymanager";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import IconButton from "@/components/ui/button/IconButton.vue";

const router = useRouter();
const route = useRoute();
const {toast} = useToast();

const relation = ref<Relation>();
const addModuleDialogOpened = ref(false)

onMounted(async () => {
  let relationid: string = route.params['id'] as string
  relation.value = await GetRelation(relationid)
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
        :module-config="null"
    />
  </DashboardLayout>
</template>

<style scoped>

</style>