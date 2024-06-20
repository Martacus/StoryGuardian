<script setup lang="ts">

import {Card, CardHeader} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import {ArrowLeft, Plus, Settings} from "lucide-vue-next";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import {Relation} from "../../bindings/storyguardian/internal/project";
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
import {
  GetEntity,
  SetEntityDescription,
  SetEntityName
} from "../../bindings/storyguardian/internal/project/entitymanager";
import {useToast} from "@/components/ui/toast";

const router = useRouter();
const route = useRoute();
const {toast} = useToast();

const relation = ref<Relation>();
const addModuleDialogOpened = ref(false)

onMounted(async () => {
  let relationid: string = route.params['id'] as string
  const getRelation = await GetRelation(relationid)
  if (getRelation) {
    relation.value = getRelation;
  }
});

async function saveRelationTitle(title: string) {
  // if (!relation.value) return;
  // try {
  //   await SetRelationName(relation.value.id, title)
  // } catch (error: any) {
  //   toast({
  //     title: 'Failed to save relation title',
  //     description: error,
  //   });
  // }
}

async function saveDescription(descriptionValue: string) {
  // if (!relation.value) return;
  //
  // try {
  //   relation.value.description = await SetRelationDescription(relation.value.id, descriptionValue);
  // } catch (error: any) {
  //   toast({
  //     title: 'Failed to save relation description',
  //     description: error,
  //   });
  // }
}
</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 col-span-4" v-if="relation">
      <CardHeader class="flex flex-row items-center justify-between py-4">
        <Button class="btn btn-secondary" variant="outline" size="icon" @click="router.back()">
          <ArrowLeft />
        </Button>
        <EntityTitle :title="relation.name" @save-title="saveRelationTitle" class="flex flex-1 justify-center"/>
        <div class="flex flex-row gap-2">
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
              <DialogDescription>Choose a module to add to your relation, you can always remove them.</DialogDescription>
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

    <Description v-if="relation" :description="relation.description" @save-description="saveDescription"/>
  </DashboardLayout>
</template>

<style scoped>

</style>