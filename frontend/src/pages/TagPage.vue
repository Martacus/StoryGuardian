<script setup lang="ts">
import {useRoute, useRouter} from "vue-router";
import {ArrowLeft} from "lucide-vue-next";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import PageHeaderCard from "@/components/shared/PageHeaderCard.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import ModuleCard from "@/components/shared/card/ModuleCard.vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {useItemFilter} from "@/composables/useItemFilter";
import {onMounted, ref} from "vue";
import {GetEntitiesByTag} from "../../bindings/storyguardian/src/project/entitymanager";
import BasicListItem from "@/components/story/entity-list/BasicListItem.vue";
import {ScrollArea} from "@/components/ui/scroll-area";
import {useNavigation} from "@/composables/useNavigation";
import {Entity} from "../../bindings/storyguardian/src/project";
import EntityTitle from "@/components/shared/EntityTitle.vue";
import {useToast} from "@/components/ui/toast";
import {EditTagName} from "../../bindings/storyguardian/src/project/tagmanager";

const route = useRoute();
const router = useRouter();
const {navigateToEntity} = useNavigation();
const {toast} = useToast();
const tag = ref<string>(route.params['id'] as string);
const entityList = ref<Entity[]>([]);

const {searchInput, searchResult} = useItemFilter(entityList, (entity, filter) => {
  return entity.name.toLowerCase().includes(filter.toLowerCase());
});

onMounted(() => {
  GetEntitiesByTag(tag.value).then((entities) => {
    entityList.value = entities;
    searchResult.value = entities;
  }).catch((error) => {
    console.error(error);
  });
});

async function editTagName(newTag: string) {
  EditTagName(tag.value, newTag).then(() => {
    tag.value = newTag;
  }).catch((error: string) => {
    toast({
      title: 'Uh oh! Something went wrong updating the tag.',
      description: error,
    });
  });
}
</script>

<template>
  <DashboardLayout>
    <PageHeaderCard>
      <IconButton @click="router.back()">
        <ArrowLeft/>
      </IconButton>

      <EntityTitle :title="tag" @save-title="editTagName" class="flex flex-1 justify-center"/>

      <div class="flex flex-row gap-2">
        <!--        <Dialog v-model:open="addModuleDialogOpened">-->
        <!--          <DialogTrigger>-->
        <!--            <TextToolTip text="Add a module">-->
        <!--              <IconButton @click="refreshUnusedEntityModules">-->
        <!--                <Plus/>-->
        <!--              </IconButton>-->
        <!--            </TextToolTip>-->
        <!--          </DialogTrigger>-->
        <!--          <DialogContent>-->
        <!--            <DialogHeader>-->
        <!--              <DialogTitle>Select a module</DialogTitle>-->
        <!--            </DialogHeader>-->
        <!--            <DialogDescription>Choose a module to add to your entity, you can always remove them.</DialogDescription>-->
        <!--            <div class="flex flex-row gap-2">-->
        <!--              <ModuleSelectItem v-if="isUnused('tagList')" @click="addEntityModule('tagList')">-->
        <!--                <p>Tags</p>-->
        <!--              </ModuleSelectItem>-->
        <!--              <ModuleSelectItem v-if="isUnused('relations')" @click="addEntityModule('relations')">-->
        <!--                <p>Relations</p>-->
        <!--              </ModuleSelectItem>-->
        <!--            </div>-->
        <!--          </DialogContent>-->
        <!--        </Dialog>-->
        <!--        <TextToolTip text="Story settings">-->
        <!--          <IconButton>-->
        <!--            <Settings/>-->
        <!--          </IconButton>-->
        <!--        </TextToolTip>-->
      </div>
    </PageHeaderCard>
    <ModuleCard title="Entities" class="col-span-4">

      <template #center-space>
        <ItemSearch v-model:search-input="searchInput" placeholder="Search tags..."/>
      </template>

      <template #card-content>
        <ScrollArea class="w-full">
          <div id="single-entity-list" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2" ref="contentRef">
            <BasicListItem v-for="entity in searchResult" :text="entity.name" @click="navigateToEntity(entity.id)">
            </BasicListItem>
          </div>
          <p v-if="searchResult.length <= 0">
            No Tags have been found.
          </p>
        </ScrollArea>
      </template>
    </ModuleCard>
  </DashboardLayout>
</template>

<style scoped>

</style>