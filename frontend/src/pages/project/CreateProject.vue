<script setup lang="ts">

import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {toTypedSchema} from "@vee-validate/zod";
import {Field, useForm} from "vee-validate";
import {z} from "zod";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";

const formSchema = toTypedSchema(z.object({
  storyname: z.string().min(2).max(50),
  description: z.string()
}))


const { handleSubmit } = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async () => {
  // const { data, error } = await supabase
  //     .from('stories')
  //     .insert([
  //       { name: values.storyname, description: values.description },
  //     ])
  //     .select();
  //
  // if (error) {
  //   console.error(error);
  // }
  //
  // if(data)
  // {
  //   storyStore.setStoryId(data[0].id);
  //   navigateTo('/dashboard');
  // }
})
</script>

<template>
  <main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
    <Card class="m-auto max-w-lg lg:min-w-96">
      <CardHeader>
        <CardTitle class="text-2xl">
          Create a new story
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form class="space-y-6" @submit="onSubmit">
          <Field v-slot="{ componentField }" name="storyname">
            <FormItem>
              <FormLabel>Story Name</FormLabel>
              <FormControl>
                <Input type="text" placeholder="The first Guardian" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </Field>
          <Field v-slot="{ componentField }" name="description">
            <FormItem>
              <FormLabel>Story Description</FormLabel>
              <FormControl>
                <Textarea type="text" placeholder="Story description" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </Field>
          <Button type="submit" class="w-full">
            Create
          </Button>
        </form>
      </CardContent>
    </Card>
  </main>
</template>

<style scoped>

</style>