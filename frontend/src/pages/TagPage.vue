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

const route = useRoute();
const router = useRouter();
const {navigateToEntity} = useNavigation();

const tag: string = route.params['id'] as string
const entityList = ref<Entity[]>([]);

const {searchInput, searchResult} = useItemFilter(entityList, (entity, filter) => {
  return entity.name.toLowerCase().includes(filter.toLowerCase());
});

onMounted(() => {
  GetEntitiesByTag(tag).then((entities) => {
    entityList.value = entities;
    searchResult.value = entities;
  }).catch((error) => {
    console.error(error);
  })


})
</script>

<template>
  <DashboardLayout>
    <PageHeaderCard>
      <IconButton @click="router.back()">
        <ArrowLeft/>
      </IconButton>
      <div class="flex flex-row gap-2 items-center m-auto">
        <p class="text-2xl leading-loose ml-2">{{ tag }}</p>
      </div>
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