<script setup lang="ts">
import TagListBase from "@/components/shared/TagListBase.vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import {Field, useForm} from 'vee-validate';
import {toTypedSchema} from '@vee-validate/zod';
import {z} from 'zod';
import {Button} from "@/components/ui/button";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Input} from "@/components/ui/input";
import {useToast} from "@/components/ui/toast";
import {CreateTag} from "../../../bindings/storyguardian/src/project/storymanager";
import {Dialog, DialogContent, DialogFooter, DialogTrigger} from "@/components/ui/dialog";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Plus, Trash2} from "lucide-vue-next";
import {ref} from "vue";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange', 'refreshTags'])

const {toast} = useToast();

const dialogOpen = ref(false);

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


</script>

<template>
  <TagListBase :tags="props.tags" :moduleConfig="props.moduleConfig" @configChange="emit('configChange', $event)">
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
    <template #item-action>
      <IconButton class="absolute right-0 flex-shrink-0 hidden group-hover:flex">
        <Trash2/>
      </IconButton>
    </template>
  </TagListBase>
</template>