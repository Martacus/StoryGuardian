<!-- frontend/src/components/entity/EntityTagList.vue -->
<script setup lang="ts">
import TagListBase from "@/components/shared/TagListBase.vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Dialog, DialogContent, DialogTrigger} from "@/components/ui/dialog";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Plus} from "lucide-vue-next";
import {onMounted, ref} from "vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {ScrollArea} from "@/components/ui/scroll-area";
import {useItemFilter} from "@/composables/useItemFilter";
import {GetStoryTags} from "../../../bindings/storyguardian/src/project/storymanager";
import {useToast} from "@/components/ui/toast";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();

const dialogOpen = ref(false);
const listHeight = ref<string>('h-0');
const tagList = ref<string[]>([]);

const {toast} = useToast()
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
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
    tagList.value = props.tags
    searchResult.value = props.tags
    calcListHeight();
  });
})
</script>

<template>
  <TagListBase :tags="props.tags" :moduleConfig="props.moduleConfig">
    <template #add-dialog>
      <Dialog v-model:open="dialogOpen">
        <DialogTrigger>
          <TextTooltip text="Add a tag">
            <IconButton @click="">
              <Plus/>
            </IconButton>
          </TextTooltip>
        </DialogTrigger>
        <DialogContent>
          <ItemSearch v-model:search-input="searchInput" placeholder="Search tags..." class="mx-4"/>
          <ScrollArea class="w-full" :class="listHeight" >
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
              <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
                   @click=""
                   v-for="tag in searchResult">
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