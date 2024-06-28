<script setup lang="ts">
import {Plus} from 'lucide-vue-next';
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
import {CreateTag} from "../../../bindings/storyguardian/src/project/storymanager";
import {useToggleBody} from "@/composables/useToggleBody";
import {useGridSize} from "@/composables/useGridSize";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import {useItemGridLayout} from "@/composables/useItemGridLayout";
import {useItemFilter} from "@/composables/useItemFilter";
import ItemSearch from "@/components/shared/ItemSearch.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";
import ItemViewSelector from "@/components/shared/button/ItemViewSelector.vue";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const dialogOpen = ref(false);
const router = useRouter();
const listHeight = ref<string>('h-0');
const tagList = ref<string[]>([]);

const {toast} = useToast();
const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig);
const {columnSize, changeGridSize} = useGridSize(props.moduleConfig);
const {itemView, changeItemView} = useItemGridLayout(props.moduleConfig);
const {searchInput, searchResult} = useItemFilter(tagList, (tag, filter) => {
  return tag.toLowerCase().includes(filter.toLowerCase());
});


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
    tagList.value.push(values.tag)
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

function calcListHeight() {
  if (itemView.value === 'list') {
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
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Tags</CardTitle>
      <ItemSearch v-if="showCardBody" v-model:search-input="searchInput" placeholder="Search tags..."/>
      <div class="flex flex-row space-x-2">
        <Dialog v-model:open="dialogOpen" v-if="showCardBody">
          <DialogTrigger>
            <TextTooltip text="Add a tag">
              <IconButton @click="">
                <Plus/>
              </IconButton>
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
        <VerticalSeperator/>
        <GridSizeSelector
            v-if="moduleConfig"
            :column-size="moduleConfig.configuration['columnSize']"
            @update-grid-size="(newSize) => changeGridSize('tagList', newSize, emit)"
        />
        <ItemViewSelector
            :item-view="itemView"
            :show-card-body="showCardBody"
            @toggle="(payload) => changeItemView('tagList', payload, emit)"
        />
        <CardBodyToggler
            :show-card-body="showCardBody"
            @toggle="toggleCardBody('tagList', emit)"
        />

      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <ScrollArea class="w-full" :class="listHeight">
        <div class="flex flex-col gap-2" v-if="itemView === 'list'">
          <div class="bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click=""
               v-for="tag in searchResult">
            <p class="px-4 text-center">
              {{ tag }}
            </p>
          </div>
          <p v-if="searchResult.length <= 0">
            No Tags have been found.
          </p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2"
             v-if="itemView === 'grid'">
          <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
               @click=""
               v-for="tag in searchResult">
            <p class="px-4 text-center">
              {{ tag }}
            </p>
          </div>
          <p v-if="searchResult.length <= 0">
            No Tags have been found.
          </p>
        </div>
      </ScrollArea>

    </CardContent>
  </Card>
</template>

