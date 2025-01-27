<script setup lang="ts">
import TagListBase from "@/components/shared/TagListBase.vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Dialog, DialogContent, DialogTrigger} from "@/components/ui/dialog";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Diff} from "lucide-vue-next";
import {computed, onMounted, ref, watch} from "vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {ScrollArea} from "@/components/ui/scroll-area";
import {useItemFilter} from "@/composables/useItemFilter";
import {GetStoryTags} from "../../../bindings/storyguardian/src/project/storymanager";
import {useToast} from "@/components/ui/toast";
import {AddTagToEntity, RemoveTagFromEntity} from "../../../bindings/storyguardian/src/project/entitymanager";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule,
  entityId: string
}>();

const emit = defineEmits(['update:tags']);

const dialogOpen = ref(false);
const listHeight = ref<string>('h-0');
const tagList = ref<string[]>([]);

const {toast} = useToast();
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});

const isTagSelected = computed(() => {
  return (tag: string) => props.tags.includes(tag);
});

function calcListHeight() {
  if (searchResult.value.length > 24) {
    listHeight.value = 'h-96';
  } else {
    listHeight.value = 'h-' + Math.max(searchResult.value.length, 3) / 3 * 12;
  }
}

onMounted(() => {
  GetStoryTags().then((tags) => {
    tagList.value = tags;
    searchResult.value = tags;
    calcListHeight();
  }).catch((error) => {
    toast({
      title: 'Failed to load story tags',
      description: error,
    });
    tagList.value = props.tags;
    searchResult.value = props.tags;
    calcListHeight();
  });
});

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
    } catch (error: string) {
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
    } catch (error: string) {
      toast({
        title: 'Failed to add tag',
        description: error,
      });
    }
  }
  emit('update:tags', newTags);
}

watch(() => props.tags, () => {
  calcListHeight();
});
</script>

<template>
  <TagListBase :tags="props.tags" :moduleConfig="props.moduleConfig">
    <template #add-dialog>
      <Dialog v-model:open="dialogOpen">
        <DialogTrigger>
          <TextTooltip text="Add a tag">
            <IconButton @click="">
              <Diff/>
            </IconButton>
          </TextTooltip>
        </DialogTrigger>
        <DialogContent class="max-w-xl">
          <ItemSearch v-model:search-input="searchInput" placeholder="Search tags..." class="mx-4"/>
          <ScrollArea class="w-full" :class="listHeight">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              <div
                  :class="{'bg-muted/90 border border-white': isTagSelected(tag), 'bg-muted/30 hover:bg-muted/40': !isTagSelected(tag)}"
                  class="rounded-lg py-2 hover:cursor-pointer"
                  @click="addOrRemoveTags(tag)"
                  v-for="tag in searchResult"
                  :key="tag"
              >
                <p class="px-4 text-center">
                  {{ tag }}
                </p>
              </div>
            </div>
            <p v-if="searchResult.length <= 0">
              No Tags have been found.
            </p>
          </ScrollArea>
        </DialogContent>
      </Dialog>
    </template>
  </TagListBase>
</template>