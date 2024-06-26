<script setup lang="ts">
import {Plus} from 'lucide-vue-next';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Field, useForm} from 'vee-validate';
import {onMounted, onUnmounted, onUpdated, ref, watch} from "vue";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Button} from "@/components/ui/button";
import {Toast, useToast} from "@/components/ui/toast";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {ScrollArea} from "@/components/ui/scroll-area";
import {Textarea} from "@/components/ui/textarea";
import {v4} from "uuid";
import {useRouter} from "vue-router";
import {Entity, Story, StoryModule} from "../../../bindings/storyguardian/src/project";
import {CreateEntity, LoadEntities} from "../../../bindings/storyguardian/src/project/entitymanager";
import {useToggleBody} from "@/composables/useToggleBody";
import {useGridSize} from "@/composables/useGridSize";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import {useItemGridLayout} from "@/composables/useItemGridLayout";
import {useItemFilter} from "@/composables/useItemFilter";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";
import ItemViewSelector from "@/components/shared/button/ItemViewSelector.vue";


const props = defineProps<{
  story: Story
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast()
const router = useRouter()

const entities = ref<Entity[]>([]);
const dialogOpen = ref(false);
const scrollAreaRef = ref<any>(null);
const contentRef = ref<any>(null);
const isScrollable = ref(false);
const listHeight = ref<string>('h-0')

const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize} = useGridSize(props.moduleConfig)
const {itemView, changeItemView} = useItemGridLayout(props.moduleConfig);
const {searchInput, searchResult} = useItemFilter(entities, (entity, filter) => {
  return entity.name.toLowerCase().includes(filter.toLowerCase());
});

onMounted(async () => {
  await getEntities();

  checkScrollable();
  window.addEventListener('resize', checkScrollable);
})

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

function refreshViewLength(){
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
    } as Entity);

    entities.value.push(entity);
    toast({
      title: 'Success',
      description: 'Entity successfully created.',
      icon: 'check',
    } as Toast);
    dialogOpen.value = false;
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})

async function navigateToEntity(id: string) {
  await router.push('/entity/' + id)
}

const checkScrollable = () => {
  if (scrollAreaRef.value && contentRef.value) {
    const scrollAreaEl = scrollAreaRef.value.$el; // Get the DOM element of the ScrollArea component
    isScrollable.value = contentRef.value.scrollHeight > scrollAreaEl.clientHeight;
  }
};

onUpdated(() => {
  checkScrollable();
});

onUnmounted(() => {
  window.removeEventListener('resize', checkScrollable);
});

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
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle> Entities</CardTitle>
      <ItemSearch v-if="showCardBody" v-model:search-input="searchInput" placeholder="Search entities..."/>
      <div class="flex flex-row space-x-2">

        <!--    Add Entity    -->
        <Dialog v-model:open="dialogOpen" v-if="showCardBody">
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

        <VerticalSeperator/>

        <GridSizeSelector
            v-if="moduleConfig"
            :column-size="moduleConfig.configuration['columnSize']"
            @update-grid-size="(newSize) => changeGridSize('entityList', newSize, emit)"
        />
        <ItemViewSelector
            :item-view="itemView"
            :show-card-body="showCardBody"
            @toggle="(payload) => changeItemView('entityList', payload, emit)"
        />
        <CardBodyToggler
            :show-card-body="showCardBody"
            @toggle="toggleCardBody('entityList', emit)"
        />

      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <ScrollArea :class="listHeight" type="auto" ref="scrollAreaRef">
        <div id="single-entity-list" class="flex flex-col gap-2" :class="{'mr-4': isScrollable}"
             v-if="itemView === 'list'" ref="contentRef">
          <!-- Entities -->
          <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToEntity(entity.id)"
               v-for="entity in searchResult">
            <p class="px-4 text-center">
              {{ entity.name }}
            </p>
          </div>
          <!-- No entities -->
          <p v-if="searchResult.length <= 0">
            No Entities have been found.
          </p>
        </div>
        <div id="single-entity-list" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2"
             v-if="itemView === 'grid'" ref="contentRef">
          <!-- Entities -->
          <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToEntity(entity.id)"
               v-for="entity in searchResult">
            <p class="px-4 text-center">
              {{ entity.name }}
            </p>
          </div>
          <!-- No entities -->
          <p v-if="searchResult.length <= 0">
            No Entities have been found.
          </p>
        </div>
      </ScrollArea>

    </CardContent>
  </Card>
</template>

