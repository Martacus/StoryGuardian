<script setup lang="ts">

import {Button} from "@/components/ui/button";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Plus, Trash2} from "lucide-vue-next";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {onMounted, ref} from "vue";
import {Entity, RelationInfo, StoryModule} from "../../../bindings/storyguardian/src/project";
import {useToast} from "@/components/ui/toast";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog";
import {Field, useForm} from "vee-validate";
import {FormControl, FormItem, FormLabel, FormMessage} from "@/components/ui/form";
import {Textarea} from "@/components/ui/textarea";
import {Input} from "@/components/ui/input";
import {toTypedSchema} from "@vee-validate/zod";
import {z} from "zod";
import {useRouter} from "vue-router";
import {useToggleBody} from "@/composables/useToggleBody";
import {useGridSize} from "@/composables/useGridSize";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import IconButton from "@/components/ui/button/IconButton.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";
import {
  CreateRelation,
  DeleteRelation,
  LoadRelationInfo
} from "../../../bindings/storyguardian/src/project/relationmanager";

const props = defineProps<{
  entity: Entity,
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const {toast} = useToast();
const router = useRouter();
const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig)

const dialogOpen = ref(false);
const deleteRelationConfirmationOpen = ref(false);
const relToDelete = ref<RelationInfo>();

const relations = ref<RelationInfo[]>([]);

onMounted(async () => {
  await loadRelations()
})

async function loadRelations() {
  try {
    relations.value = await LoadRelationInfo(props.entity.id, 0, 20)
  } catch (error: any) {
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

const onSubmit = handleSubmit(async () => {

});

async function createRelation(){
  try {
    const relationId = await CreateRelation(props.entity.id);
    await router.push("/relation/" + relationId);
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: error.message,
    });
  }
}

function openRelation(relationId: string){
 router.push("/relation/" + relationId);
}

function deleteRelation(){
  if(!relToDelete.value) return;
  try{
    DeleteRelation(relToDelete.value?.id).then(() => {
      loadRelations();
      deleteRelationConfirmationOpen.value = false;
    })
  } catch (error: any) {
    toast({
      title: 'Uh oh! Something went wrong removing the relation',
      description: error.message,
    });
  }
}
function openDeleteConfirmation(relationId: string){
  relToDelete.value = relations.value.find(relation => relation.id === relationId);
  deleteRelationConfirmationOpen.value = true;
}
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Relations</CardTitle>
      <div class="flex flex-row space-x-2">
        <TextTooltip text="Add relation" v-if="showCardBody">
          <IconButton @click="createRelation()">
            <Plus/>
          </IconButton>
        </TextTooltip>
        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('relations', newSize, emit)"/>
        <CardBodyToggler :show-card-body="showCardBody"  @toggle="toggleCardBody('relations', emit)"/>

      </div>
    </CardHeader>
    <CardContent v-if="showCardBody">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Relation</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="relation in relations" @click="openRelation(relation.id)" class="hover:cursor-pointer">
            <!-- Table Data -->
            <TableCell class="font-medium">
              {{ relation.toName }}
            </TableCell>
            <TableCell>{{ relation.name }}</TableCell>
            <TableCell class="text-right">
              <TextTooltip text="Delete">
                <IconButton @click.stop="openDeleteConfirmation(relation.id)">
                  <Trash2/>
                </IconButton>
              </TextTooltip>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </CardContent>

    <Dialog v-model:open="dialogOpen">
      <DialogContent>
        <form class="space-y-6" @submit="onSubmit">
          <DialogHeader>
            <DialogTitle>Create a relation</DialogTitle>
          </DialogHeader>

          <!-- Form -->
          <Field :validate-on-blur="false" v-slot="{ componentField }" name="name">
            <FormItem>
              <FormLabel>Relation Name</FormLabel>
              <FormControl>
                <Input type="text" placeholder="The first Guardian" v-bind="componentField" autocomplete="off"/>
              </FormControl>
              <FormMessage/>
            </FormItem>
          </Field>
          <Field :validate-on-blur="false" v-slot="{ componentField }" name="description">
            <FormItem>
              <FormLabel>Relation Description</FormLabel>
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

  <Dialog v-model:open="deleteRelationConfirmationOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Are you sure you want to remove this relation?</DialogTitle>
      </DialogHeader>
      <p v-if="relToDelete">{{ relToDelete.name }}</p>
      <DialogFooter>
        <Button @click="deleteRelation" class="w-full">
          Remove
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<style scoped>

</style>