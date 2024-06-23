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
import {useRouter} from "vue-router";
import {CreateTag} from "../../../bindings/storyguardian/internal/project/storymanager";

type ListViewMode = 'grid' | 'list';

const {toast} = useToast()

const openFilter = ref(false);
const dialogOpen = ref(false);
const listView = ref<ListViewMode>('list')
const showBody = ref(true);
const router = useRouter()

const listHeight = ref<string>('h-0')

const props = defineProps<{
  tags: string[]
}>();

//Form
const formSchema = toTypedSchema(z.object({
  tag: z.string(),
}))

const {handleSubmit} = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  try {
    await CreateTag(values.tag)
    props.tags.push(values.tag)
    //Push tag to story object

    calcListHeight();
    toast({
      title: 'Success',
      description: 'Tag successfully created.',
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
function changeListView(view: ListViewMode) {
  listView.value = view;
  calcListHeight();
}

function toggleCard() {
  showBody.value = !showBody.value;
}

function calcListHeight() {
  if (listView.value === 'list') {
    if (props.tags.length > 8) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(props.tags.length, 1) * 12;
    }
  } else {
    if (props.tags.length > 24) {
      listHeight.value = 'h-96';
    } else {
      listHeight.value = 'h-' + Math.max(props.tags.length, 3) / 3 * 12;
    }
  }
}

async function navigateToTag(id: string){

}

onMounted(() => {
  calcListHeight();
})
</script>

<template>
  <Card class="bg-muted/30 col-span-4">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Tags</CardTitle>
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
                <DialogTitle>Create a Tag</DialogTitle>
              </DialogHeader>
              <!-- Form -->
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
        <div class="flex flex-col gap-2" v-if="listView === 'list'">
          <div class="bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToTag(tag)"
               v-for="tag in tags">
            <p class="px-4 text-center">
              {{ tag }}
            </p>
          </div>
          <p v-if="tags.length <= 0">
            No Tags have been found.
          </p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2"
             v-if="listView === 'grid'">
          <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click="navigateToTag(tag)"
               v-for="tag in tags">
            <p class="px-4 text-center">
              {{ tag }}
            </p>
          </div>
          <p v-if="tags.length <= 0">
            No Tags have been found.
          </p>
        </div>
      </ScrollArea>

    </CardContent>
  </Card>
</template>

