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
import {DialogFooter} from "@/components/ui/dialog";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast();
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
    <template #dialog-content>
      <!-- Custom Dialog content for TagList -->
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
    </template>
  </TagListBase>
</template>