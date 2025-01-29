<script setup lang="ts">
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import {Field, useForm} from 'vee-validate';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Button} from "@/components/ui/button";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {useToast} from "@/components/ui/toast";
import {CreateTag} from "../../../bindings/storyguardian/src/project/storymanager";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Plus} from "lucide-vue-next";
import {onMounted, ref, watch} from "vue";
import ModuleBase from "@/components/shared/module/ModuleBase.vue";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import {useItemFilter} from "@/composables/useItemFilter";
import BasicListItem from "@/components/story/entity-list/BasicListItem.vue";
import BasicItemList from "@/components/story/entity-list/BasicItemList.vue";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange', 'refreshTags'])

const {toast} = useToast();

const listHeight = ref<string>('h-0');
const tagList = ref<string[]>([]);
const itemView = ref(props.moduleConfig.configuration['itemView']);

//Composables
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});

function calcListHeight() {
  if (itemView.value === 'list') {
    console.log('val;' + searchResult.value.length)
    if (searchResult.value.length > 8) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(searchResult.value.length, 1) * 12;
    }
  } else {
    if (searchResult.value.length > 24) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(searchResult.value.length, 3) / 3 * 12;
    }
  }
}

onMounted(() => {
  tagList.value = props.tags
  searchResult.value = props.tags
  calcListHeight();
})

watch(() => props.tags, (newTags) => {
  tagList.value = newTags;
  searchResult.value = newTags;
  calcListHeight();
});

function updateItemView(view: string){
  itemView.value = view;
  calcListHeight();
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
    emit('refreshTags')
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})

function initDelete(tag: string){
  tagToDelete.value = tag;
  tagDeleteDialogOpen.value = true;
}

function deleteEntity(){
  // DeleteEntity(entityToDelete.value!.id).then(() => {
  //   entities.value = entities.value.filter(entity => entity.id !== entityToDelete.value!.id);
  //   toast({
  //     title: 'Success',
  //     description: 'Entity successfully deleted.',
  //     icon: 'check',
  //   });
  //   deleteDialogOpen.value = false;
  //   searchResult.value = entities.value;
  //   refreshViewLength();
  // }).catch((error: string) => {
  //   toast({
  //     title: 'Uh oh! Something went wrong.',
  //     description: error,
  //   });
  // });
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
          <BasicListItem v-for="tag in searchResult" :text="tag">
            <template #item-action>
              <slot name="item-action"/>
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
        <Button @click="deleteEntity" class="w-full">
          Remove
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>