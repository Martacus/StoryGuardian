<script setup lang="ts">

import {Button} from "@/components/ui/button";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {ChevronUp, ChevronDown, Plus} from "lucide-vue-next";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {onMounted, ref} from "vue";
import {Entity, RelationInfo} from "../../../bindings/storyguardian/internal/project";
import {
  CreateRelation,
  LoadRelationInfo
} from "../../../bindings/storyguardian/internal/project/entitymanager";
import {Toast, useToast} from "@/components/ui/toast";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import {Field, useForm} from "vee-validate";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Textarea} from "@/components/ui/textarea";
import {Input} from "@/components/ui/input";
import {toTypedSchema} from "@vee-validate/zod";
import {z} from "zod";

const dialogOpen = ref(false);
const showBody = ref(true);
const {toast} = useToast()

function toggleCard() {
  showBody.value = !showBody.value;
}

const props = defineProps<{
  entity: Entity
}>();

const relations = ref<RelationInfo[]>([]);

onMounted(async () => {
  await loadRelations()
})

async function loadRelations(){
  try {
    relations.value = await LoadRelationInfo(props.entity.id)
  } catch(error: any){
    toast({
      title: 'error loading relations',
      description: error
    })
  }
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
    await CreateRelation(props.entity.id, {
      description: values.description,
      name: values.name,
      to: 'e01e7ce9-36e1-4943-af20-6da01f77038a'
    });

    toast({
      title: 'Success',
      description: 'Relation successfully created.',
      icon: 'check',
    } as Toast);

    await loadRelations();
    dialogOpen.value = false;
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
})
</script>

<template>
  <Card class="bg-muted/30 col-span-2">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Relations</CardTitle>
      <div class="flex flex-row space-x-2">
        <TextTooltip text="Add relation" v-if="showBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="dialogOpen = !dialogOpen">
            <Plus/>
          </Button>
        </TextTooltip>
        <TextTooltip text="Minimize" v-if="showBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
            <ChevronUp/>
          </Button>
        </TextTooltip>
        <TextTooltip text="Minimize" v-if="!showBody">
          <Button size="icon" aria-label="Toggle italic" variant="outline" @click="toggleCard()">
            <ChevronDown/>
          </Button>
        </TextTooltip>

      </div>
    </CardHeader>
    <CardContent v-if="showBody">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead >
              Name
            </TableHead>
            <TableHead >
              Relation
            </TableHead>
            <TableHead>Description</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(relation) in relations">
            <!-- Table Data -->
            <TableCell class="font-medium">
              {{ relation.name }}
            </TableCell>
            <TableCell>{{ relation.toName }}</TableCell>
            <TableCell>{{ relation.description }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </CardContent>

    <Dialog v-model:open="dialogOpen" v-if="showBody">
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
  </Card>
</template>

<style scoped>

</style>