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
import {Entity, Story} from "../../../bindings/storyguardian/internal/project";
import {CreateEntity, LoadEntities} from "../../../bindings/storyguardian/internal/project/entitymanager";

type EntityListViewMode = 'grid' | 'list';

const {toast} = useToast()

const entities = ref<Entity[]>([]);
const showEntities = ref<Entity[]>([]);
const openFilter = ref(false);
const dialogOpen = ref(false);
const listView = ref<EntityListViewMode>('list')
const showBody = ref(true);
const router = useRouter()

const listHeight = ref<string>('h-0')

const props = defineProps<{
  story: Story
}>();

onMounted(() => {
  getEntities();
})

async function getEntities() {
  try {
    let data = await LoadEntities(props.story.id);
    entities.value = data;
    showEntities.value = data;
    calcListHeight();
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
  // listView.value = configStore.entityListView;
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
    } as Entity);

    entities.value.push(entity);
    //showEntities.value.push(data[0]); //Reapply filtering
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

function toggleCard() {
  showBody.value = !showBody.value;
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
  <div>
    <Card class="bg-muted/30">
      <CardHeader class="flex flex-row justify-between items-center">
        <CardTitle> Entities</CardTitle>
        <div class="flex flex-row space-x-2">
          <Dialog v-model:open="dialogOpen" v-if="showBody">
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

          <TextTooltip text="Switch to grid" v-if="listView === 'list' && showBody">
            <Button size="icon" aria-label="Toggle italic" variant="outline" @click="changeListView('grid')">
              <StretchHorizontal/>
            </Button>
          </TextTooltip>

          <TextTooltip text="Switch to list" v-if="listView === 'grid' && showBody">
            <Button size="icon" aria-label="Toggle italic" variant="outline" @click="changeListView('list')">
              <LayoutGrid/>
            </Button>
          </TextTooltip>

          <TextTooltip text="Expand" v-if="!showBody">
            <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
              <ChevronDown/>
            </Button>
          </TextTooltip>

          <TextTooltip text="Minimize" v-if="showBody">
            <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
              <ChevronUp/>
            </Button>
          </TextTooltip>

        </div>
      </CardHeader>
      <CardContent v-if="showBody">
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
  </div>
</template>

