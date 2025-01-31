<script setup lang="ts">
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import {computed, onMounted, ref, watch} from "vue";
import {useItemFilter} from "@/composables/useItemFilter";
import {GetStoryTags} from "../../../bindings/storyguardian/src/project/storymanager";
import {useToast} from "@/components/ui/toast";
import {AddTagToEntity, RemoveTagFromEntity} from "../../../bindings/storyguardian/src/project/entitymanager";
import {Dialog, DialogContent, DialogTrigger} from "@/components/ui/dialog";
import BasicListItem from "@/components/story/entity-list/BasicListItem.vue";
import BasicItemList from "@/components/story/entity-list/BasicItemList.vue";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {Diff} from "lucide-vue-next";
import ModuleBase from "@/components/shared/module/ModuleBase.vue";
import {ScrollArea} from "@/components/ui/scroll-area";
import {useNavigation} from "@/composables/useNavigation";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule,
  entityId: string
}>();

const emit = defineEmits(['update:tags', 'configChange']);
const {navigateToTag} = useNavigation()

const dialogOpen = ref(false);
const listHeight = ref<string>('h-0');
const tagList = ref<string[]>([]);
const storyTagsList = ref<string[]>([]);
const itemView = ref(props.moduleConfig.configuration['itemView']);

//Composables
const {toast} = useToast();
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});

const {searchInput: storyTagSearchInput, searchResult: storyTagsSearchResult} = useItemFilter(storyTagsList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});

//computes
const isTagSelected = computed(() => {
  return (tag: string) => props.tags.includes(tag);
});

//Hooks
onMounted(() => {
  GetStoryTags().then((tags) => {
    storyTagsList.value = tags;
    storyTagsSearchResult.value = tags;
    calcListHeight();
  }).catch((error) => {
    toast({
      title: 'Failed to load story tags',
      description: error,
    });
    calcListHeight();
  });

  tagList.value = props.tags;
  searchResult.value = props.tags;
});

//Watch
watch(() => props.tags, () => {
  tagList.value = props.tags;
  calcListHeight();
});

//methods
function calcListHeight() {
  if (searchResult.value.length > 24) {
    listHeight.value = 'h-96';
  } else {
    listHeight.value = 'h-' + Math.max(searchResult.value.length, 3) / 3 * 12;
  }
}

async function addOrRemoveTags(tag: string) {
  const newTags = [...props.tags];
  if (newTags.includes(tag)) {
    newTags.splice(newTags.indexOf(tag), 1);
    try {
      await RemoveTagFromEntity(props.entityId, tag);
      toast({
        title: 'Tag removed',
        description: `Tag ${tag} has been removed from the entity.`,
      });
    } catch (error: any) {
      toast({
        title: 'Failed to remove tag',
        description: error,
      });
    }
  } else {
    newTags.push(tag);
    try {
      await AddTagToEntity(props.entityId, tag);
      toast({
        title: 'Tag added',
        description: `Tag ${tag} has been added to the entity.`,
      });
    } catch (error: any) {
      toast({
        title: 'Failed to add tag',
        description: error,
      });
    }
  }
  emit('update:tags', newTags);
}

function updateItemView(view: string){
  itemView.value = view;
}
</script>

<template>
  <ModuleBase title="Tags" :module-config="moduleConfig"
              @config-change="(payload) => emit('configChange', payload)"
              @update:item-view="updateItemView"
              :item-grid-layout="true">

    <template #side-buttons>
      <Dialog v-model:open="dialogOpen">
        <DialogTrigger>
          <TextTooltip text="Add a tag">
            <IconButton @click="">
              <Diff/>
            </IconButton>
          </TextTooltip>
        </DialogTrigger>
        <DialogContent class="max-w-xl">
          <ItemSearch v-model:search-input="storyTagSearchInput" placeholder="Search tags..." class="mx-4"/>
          <ScrollArea class="w-full">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              <div
                  :class="{'bg-muted/90 border border-white': isTagSelected(tag), 'bg-muted/30 hover:bg-muted/40': !isTagSelected(tag)}"
                  class="rounded-lg py-2 hover:cursor-pointer"
                  @click="addOrRemoveTags(tag)"
                  v-for="tag in storyTagsSearchResult"
                  :key="tag"
              >
                <p class="px-4 text-center">
                  {{ tag }}
                </p>
              </div>
            </div>
            <p v-if="storyTagsSearchResult.length <= 0">
              No Tags have been found.
            </p>
          </ScrollArea>
        </DialogContent>
      </Dialog>
    </template>

    <template #center-space>
      <ItemSearch v-model:search-input="searchInput" placeholder="Search tags..."/>
    </template>

    <template #card-content>
      <BasicItemList :list-height="listHeight" :item-view="itemView">

        <template #items>
          <BasicListItem v-for="tag in searchResult" :text="tag" @click="navigateToTag(tag)">
          </BasicListItem>
        </template>

        <template #no-items>
          <p v-if="searchResult.length <= 0">
            No Tags have been found.
          </p>
        </template>

      </BasicItemList>
    </template>
  </ModuleBase>
</template>