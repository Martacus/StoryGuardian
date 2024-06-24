<script setup lang="ts">
import {ChevronDown, ChevronUp, LayoutGrid, Plus, StretchHorizontal} from 'lucide-vue-next';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Field, useForm} from 'vee-validate';
import {onMounted, ref} from "vue";
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
import GridSizeSelector from "@/components/shared/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeperator.vue";
type EntityListViewMode = 'grid' | 'list';

const props = defineProps<{
  story: Story
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast()

const entities = ref<Entity[]>([]);
const showEntities = ref<Entity[]>([]);
const dialogOpen = ref(false);
const listView = ref<EntityListViewMode>('list') 
const router = useRouter()

const listHeight = ref<string>('h-0')


const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig)

onMounted(() => {
  getEntities();
})

async function getEntities() {
  try {
    let data = await LoadEntities();
    entities.value = data;
    showEntities.value = data;
    calcListHeight();
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
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
    calcListHeight();
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

//View
function changeListView(view: EntityListViewMode) {
  listView.value = view;
  calcListHeight();
} 

function calcListHeight() {
  if (listView.value === 'list') {
    if (entities.value.length > 8) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + entities.value.length * 12;
    }
  } else {
    if (entities.value.length > 24) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(entities.value.length, 3) / 3 * 12;
    }
  }
}

async function navigateToEntity(id: string){
  await router.push('/entity/' + id)
}
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle> Entities</CardTitle>
      <div class="flex flex-row space-x-2">
        <Dialog v-model:open="dialogOpen" v-if="showCardBody">
          <DialogTrigger>
            <TextTooltip text="Add an entity">
              <Button size="icon" aria-label="Toggle italic" variant="outline" @click="">
                <Plus/>
              </Button>
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
        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('entityList', newSize, emit)"/>
        <TextTooltip text="Switch to grid" v-if="listView === 'list' && showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="changeListView('grid')">
            <StretchHorizontal/>
          </Button>
        </TextTooltip>

        <TextTooltip text="Switch to list" v-if="listView === 'grid' && showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="changeListView('list')">
            <LayoutGrid/>
          </Button>
        </TextTooltip>

        <TextTooltip text="Expand" v-if="!showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCardBody('entityList', emit)">
            <ChevronDown/>
          </Button>
        </TextTooltip>

        <TextTooltip text="Minimize" v-if="showCardBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCardBody('entityList', emit)">
            <ChevronUp/>
          </Button>
        </TextTooltip>

      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <ScrollArea class="w-full" :class="listHeight">
        <div id="single-entity-list" class="flex flex-col gap-2" v-if="listView === 'list'">
          <!-- Entities -->
          <div :to="'/entity/' + entity.id" class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToEntity(entity.id)"
               v-for="entity in showEntities">
            <p class="px-4 text-center">
              {{ entity.name }}
            </p>
          </div>
          <!-- No entities -->
          <p v-if="showEntities.length <= 0">
            No Entities have been found.
          </p>
        </div>
        <div id="single-entity-list" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2"
             v-if="listView === 'grid'">
          <!-- Entities -->
          <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToEntity(entity.id)"
               v-for="entity in showEntities">
            <p class="px-4 text-center">
              {{ entity.name }}
            </p>
          </div>
          <!-- No entities -->
          <p v-if="showEntities.length <= 0">
            No Entities have been found.
          </p>
        </div>
      </ScrollArea>

    </CardContent>
  </Card>
</template>

