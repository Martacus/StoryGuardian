<script setup lang="ts">
import {Plus, Trash2} from 'lucide-vue-next';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Field, useForm} from 'vee-validate';
import {onMounted, onUnmounted, onUpdated, ref, watch} from "vue";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Button} from "@/components/ui/button";
import {useToast} from "@/components/ui/toast";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {ScrollArea} from "@/components/ui/scroll-area";
import {Textarea} from "@/components/ui/textarea";
import {v4} from "uuid";
import {Entity, Story, StoryModule} from "../../../../bindings/storyguardian/src/project";
import {CreateEntity, DeleteEntity, LoadEntities} from "../../../../bindings/storyguardian/src/project/entitymanager";
import {useItemFilter} from "@/composables/useItemFilter";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import {useNavigation} from "@/composables/useNavigation";
import BasicListItem from "@/components/story/entity-list/BasicListItem.vue";
import ModuleBase from "@/components/shared/module/ModuleBase.vue";


const props = defineProps<{
  story: Story
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast()
const {navigateToEntity} = useNavigation()

const entities = ref<Entity[]>([]);
const dialogOpen = ref(false);
const scrollAreaRef = ref<any>(null);
const contentRef = ref<any>(null);
const isScrollable = ref(false);
const listHeight = ref<string>('h-0')

//Config values
const itemView = ref(props.moduleConfig.configuration['itemView']);

//Delete dialog
const deleteDialogOpen = ref(false);
const entityToDelete = ref<Entity>();

//Composables
const {searchInput, searchResult} = useItemFilter(entities, (entity, filter) => {
  return entity.name.toLowerCase().includes(filter.toLowerCase());
});

//Hooks
onMounted(async () => {
  await getEntities();

  checkScrollable();
  window.addEventListener('resize', checkScrollable);
})

onUpdated(() => {
  checkScrollable();
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScrollable);
});

//Functions
async function getEntities() {
  try {
    let data = await LoadEntities();
    entities.value = data;
    searchResult.value = data;
    refreshViewLength();
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
}

function refreshViewLength() {
  if (itemView.value === 'list') {
    if (searchResult.value.length > 8) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + entities.value.length * 12;
    }
  } else {
    if (searchResult.value.length > 24) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(entities.value.length, 3) / 3 * 12;
    }
  }
}

const checkScrollable = () => {
  if (scrollAreaRef.value && contentRef.value) {
    const scrollAreaEl = scrollAreaRef.value.$el;
    isScrollable.value = contentRef.value.scrollHeight > scrollAreaEl.clientHeight;
  }
};

function initDelete(entityId: string){
  entityToDelete.value = entities.value.find(entity => entity.id === entityId);
  deleteDialogOpen.value = true;
}

function deleteEntity(){
  DeleteEntity(entityToDelete.value!.id).then(() => {
    entities.value = entities.value.filter(entity => entity.id !== entityToDelete.value!.id);
    toast({
      title: 'Success',
      description: 'Entity successfully deleted.',
      icon: 'check',
    });
    deleteDialogOpen.value = false;
    searchResult.value = entities.value;
    refreshViewLength();
  }).catch((error: string) => {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error,
    });
  });
}

//Form
const formSchema = toTypedSchema(z.object({
  name: z.string(),
  description: z.string(),
}))

const {handleSubmit} = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  try {
    let entity = await CreateEntity({
      id: v4(),
      name: values.name,
      description: values.description,
      storyId: props.story.id,
      modules: {}
    } as Entity);

    entities.value.push(entity);
    toast({
      title: 'Success',
      description: 'Entity successfully created.',
      icon: 'check',
    });
    dialogOpen.value = false;
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})



//Watches
watch(
    () => [scrollAreaRef.value, contentRef.value],
    () => {
      checkScrollable();
    },
    {immediate: true}
);

watch(itemView, () => {
  refreshViewLength();
})
</script>

<template>
  <ModuleBase title="Entities"
              :module-config="moduleConfig"
              @config-change="(payload) => emit('configChange', payload)"
              @update:item-view="args => itemView = args"
              :item-grid-layout="true">

    <template #center-space>
      <ItemSearch  v-model:search-input="searchInput" placeholder="Search entities..."/>
    </template>

    <template #side-buttons>
      <Dialog v-model:open="dialogOpen">
        <DialogTrigger>
          <TextTooltip text="Add an entity">
            <IconButton @click="">
              <Plus/>
            </IconButton>
          </TextTooltip>
        </DialogTrigger>
        <DialogContent>
          <form class="space-y-6" @submit="onSubmit">
            <DialogHeader>
              <DialogTitle>Create an entity</DialogTitle>
            </DialogHeader>
            <!-- Form -->
            <Field :validate-on-blur="false" v-slot="{ componentField }" name="name">
              <FormItem>
                <FormLabel>Entity Name</FormLabel>
                <FormControl>
                  <Input type="text" placeholder="The first Guardian" v-bind="componentField" autocomplete="off"/>
                </FormControl>
                <FormMessage/>
              </FormItem>
            </Field>
            <Field :validate-on-blur="false" v-slot="{ componentField }" name="description">
              <FormItem>
                <FormLabel>Entity Description</FormLabel>
                <FormControl>
                      <Textarea type="text" placeholder="The first guardian of Xybal" v-bind="componentField"
                                autocomplete="off"/>
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

    <template #card-content>
      <ScrollArea :class="listHeight" type="auto" ref="scrollAreaRef">
        <div id="single-entity-list" :class="[itemView === 'list' ? 'flex flex-col gap-2' : 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2', {'mr-4': isScrollable}]" ref="contentRef">
          <BasicListItem v-for="entity in searchResult" :text="entity.name" @click="navigateToEntity(entity.id)">
            <template #item-action>
              <IconButton @click.stop="initDelete(entity.id)" class="absolute right-0 flex-shrink-0 hidden group-hover:flex">
                <Trash2/>
              </IconButton>
            </template>
          </BasicListItem>
        </div>
        <p v-if="searchResult.length <= 0">
          No Entities have been found.
        </p>
      </ScrollArea>
    </template>
  </ModuleBase>

  <Dialog v-model:open="deleteDialogOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Are you sure you want to remove this entity?</DialogTitle>
      </DialogHeader>
      <p v-if="entityToDelete">{{ entityToDelete.name }}</p>
      <DialogFooter>
        <Button @click="deleteEntity" class="w-full">
          Remove
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

