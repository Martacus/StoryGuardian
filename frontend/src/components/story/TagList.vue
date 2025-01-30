<script setup lang="ts">
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import {Field, useForm} from 'vee-validate';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Button} from "@/components/ui/button";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {useToast} from "@/components/ui/toast";
import {CreateTag, RemoveTagFromStory} from "../../../bindings/storyguardian/src/project/storymanager";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Plus, Trash2} from "lucide-vue-next";
import {onMounted, ref, watch} from "vue";
import ModuleBase from "@/components/shared/module/ModuleBase.vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {useItemFilter} from "@/composables/useItemFilter";
import BasicListItem from "@/components/story/entity-list/BasicListItem.vue";
import BasicItemList from "@/components/story/entity-list/BasicItemList.vue";
import {useBasicListHeight} from "@/composables/list/useBasicListHeight";
import {useNavigation} from "@/composables/useNavigation";
import {RemoveTagFromEntities} from "../../../bindings/storyguardian/src/project/entitymanager";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange', 'update:tags'])

const {toast} = useToast();
const {navigateToTag} = useNavigation()

const tagList = ref<string[]>([]);
const itemView = ref(props.moduleConfig.configuration['itemView']);

//Composables
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});
const {listHeight} = useBasicListHeight(itemView, searchResult);

onMounted(() => {
  tagList.value = props.tags
  searchResult.value = props.tags
})

watch(() => props.tags, (newTags) => {
  tagList.value = newTags;
  searchResult.value = newTags;
  console.log('update')
});

function updateItemView(view: string) {
  itemView.value = view;
}

//Add dialog
const dialogOpen = ref(false);

//Delete dialog
const tagDeleteDialogOpen = ref(false);
const tagToDelete = ref<string>('');

const {handleSubmit} = useForm({
  validationSchema: toTypedSchema(z.object({
    tag: z.string(),
  })),
})

const onSubmit = handleSubmit(async (values) => {
  try {
    await CreateTag(values.tag)
    toast({
      title: 'Success',
      description: 'Tag successfully created.',
      icon: 'check',
    });
    dialogOpen.value = false;
    const newTagList = [...tagList.value, values.tag];
    emit('update:tags', newTagList);
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})

function initDelete(tag: string) {
  tagToDelete.value = tag;
  tagDeleteDialogOpen.value = true;
}

function deleteTagFromStory() {
  RemoveTagFromStory(tagToDelete.value).then(() => {
    const newTagList = tagList.value.filter(tag => tag !== tagToDelete.value);
    emit('update:tags', newTagList);
  }).catch((error: string) => {
    toast({
      title: 'Uh oh! Something went wrong removing the tag from the story.',
      description: error,
    });
  });

  RemoveTagFromEntities(tagToDelete.value).catch((error: string) => {
    toast({
      title: 'Uh oh! Something went wrong removing the tag from entities.',
      description: error,
    });
  });

  toast({
    title: 'Success',
    description: 'Tag successfully removed.',
    icon: 'check',
  });

  tagDeleteDialogOpen.value = false;
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
              <Plus/>
            </IconButton>
          </TextTooltip>
        </DialogTrigger>
        <DialogContent>
          <form class="space-y-6" @submit="onSubmit">
            <Field :validate-on-blur="false" v-slot="{ componentField }" name="tag">
              <FormItem>
                <FormLabel>Tag Name</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="Warrior" v-bind="componentField" autocomplete="off"/>
                </FormControl>
                <FormMessage/>
              </FormItem>
            </Field>
            <DialogFooter>
              <Button type="submit" class="w-full">
                Create
              </Button>
            </DialogFooter>
          </form>
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
            <template #item-action>
              <IconButton @click.stop="initDelete(tag)" class="absolute right-0 flex-shrink-0 hidden group-hover:flex">
                <Trash2/>
              </IconButton>
            </template>
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


  <Dialog v-model:open="tagDeleteDialogOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Are you sure you want to remove this entity?</DialogTitle>
      </DialogHeader>
      <p v-if="tagToDelete">{{ tagToDelete }}</p>
      <DialogFooter>
        <Button @click="deleteTagFromStory" class="w-full">
          Remove
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>