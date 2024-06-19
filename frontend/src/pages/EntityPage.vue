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
import {GetEntity, SetEntityDescription} from "../../bindings/storyguardian/internal/project/entitymanager";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import {Entity} from "../../bindings/storyguardian/internal/project";
import Description from "@/components/shared/Description.vue";
import {useToast} from "@/components/ui/toast";
import RelationModule from "@/components/modules/RelationModule.vue";


const route = useRoute();
const router = useRouter()
const entity = ref<Entity>();
const addModuleDialogOpened = ref(false)
const {toast} = useToast()

onMounted(async () => {
  let projectId: string = route.params['id'] as string
  const retrievedEntity = await GetEntity(projectId)
  if (retrievedEntity) {
    entity.value = retrievedEntity;
  }
})

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

</script>

<template>
  <DashboardLayout>
    <Card class="bg-muted/30 col-span-4" v-if="entity">
      <CardHeader class="flex flex-row items-center justify-between py-4">
        <Button class="btn btn-secondary" variant="outline" size="icon" @click="router.back()">
          <ArrowLeft />
        </Button>
        <EntityTitle :title="entity.name" class="flex flex-1 justify-center"/>
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
              <DialogDescription>Choose a module to add to your entity, you can always remove them.</DialogDescription>
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

    <Description v-if="entity" :description="entity.description" @save-description="saveDescription"/>
    <RelationModule v-if="entity" :entity="entity"/>
  </DashboardLayout>
</template>

<style scoped>

</style>