<script setup lang="ts">
import {onMounted, ref, watch} from 'vue';
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Relation, StoryModule} from "../../../bindings/storyguardian/src/project";
import GridSizeSelector from "@/components/shared/button/GridSizeSelector.vue";
import {useGridSize} from "@/composables/useGridSize";
import {useToggleBody} from "@/composables/useToggleBody";
import VerticalSeperator from "@/components/ui/separator/VerticalSeparator.vue";
import CardBodyToggler from "@/components/shared/button/CardBodyToggler.vue";
import {Label} from "@/components/ui/label";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Edit} from "lucide-vue-next";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Button} from "@/components/ui/button";
import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog";
import EntitySelector from "@/components/shared/EntitySelector.vue";
import {SetRelationEntities} from "../../../bindings/storyguardian/src/project/relationmanager";
import {GetEntity} from "../../../bindings/storyguardian/src/project/entitymanager";
import {useNavigation} from "@/composables/useNavigation";

const props = defineProps<{
  moduleConfig: StoryModule,
  relation: Relation
}>();

const emit = defineEmits(['configChange']);
const infoDialogOpen = ref(false);

const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig);
const {columnSize, changeGridSize} = useGridSize(props.moduleConfig);
const {navigateToEntity} = useNavigation();

const selectedEntity1 = ref('');
const selectedEntity2 = ref('');
const validationError = ref('');
const entityOneName = ref('');
const entityTwoName = ref('');

onMounted(() => {
  selectedEntity1.value = props.relation.entityOne;
  selectedEntity2.value = props.relation.entityTwo;

  if (selectedEntity1.value) {
    GetEntity(selectedEntity1.value).then((entity) => {
      if (entity) {
        entityOneName.value = entity.name;
      } else {
        entityOneName.value = 'NaN';
      }
    });
  }

  if (selectedEntity2.value) {
    GetEntity(selectedEntity2.value).then((entity) => {
      if (entity) {
        entityTwoName.value = entity.name;
      } else {
        entityTwoName.value = 'NaN';
      }
    });
  }
});

function validateEntities(entity1: string, entity2: string) {
  if (!entity1 || !entity2) {
    validationError.value = 'Both entities must be selected.';
  } else if (entity1 === entity2) {
    validationError.value = 'Entities must be different.';
  } else {
    validationError.value = '';
  }
}

function saveEntities() {
  validateEntities(selectedEntity1.value, selectedEntity2.value);
  if (!validationError.value) {
    infoDialogOpen.value = false;
  }
  SetRelationEntities(selectedEntity1.value, selectedEntity2.value, props.relation.id);
  refreshEntityNames();
}

function refreshEntityNames() {
  if (selectedEntity1.value) {
    GetEntity(selectedEntity1.value).then((entity) => {
      if (entity) {
        entityOneName.value = entity.name;
      } else {
        entityOneName.value = 'NaN';
      }
    });
  }

  if (selectedEntity2.value) {
    GetEntity(selectedEntity2.value).then((entity) => {
      if (entity) {
        entityTwoName.value = entity.name;
      } else {
        entityTwoName.value = 'NaN';
      }
    });
  }
}

watch([selectedEntity1, selectedEntity2], () => {
  validationError.value = '';
});
</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Information</CardTitle>
      <div class="flex flex-row space-x-2">
        <Dialog v-model:open="infoDialogOpen" v-if="showCardBody">
          <DialogTrigger>
            <TextTooltip text="Edit Relation">
              <IconButton @click="">
                <Edit/>
              </IconButton>
            </TextTooltip>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Relation</DialogTitle>
            </DialogHeader>
            <div id="wg-dialog-content" class="w-full flex flex-col justify-center items-center">
              <EntitySelector v-model="selectedEntity1"></EntitySelector>
              <p class="my-4">And</p>
              <EntitySelector v-model="selectedEntity2"></EntitySelector>
              <p v-if="validationError" class="text-red-500">{{ validationError }}</p>
            </div>
            <DialogFooter>
              <Button class="w-full" @click="saveEntities">
                Save
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        <VerticalSeperator/>
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']"
                          @update-grid-size="(newSize) => changeGridSize('relationInfo', newSize, emit)"/>
        <CardBodyToggler :show-card-body="showCardBody" @toggle="toggleCardBody('relationInfo', emit)"/>
      </div>
    </CardHeader>
    <CardContent v-if="showCardBody" class="flex flex-row justify-center items-center space-x-4">
      <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
           @click="navigateToEntity(selectedEntity1)">
        <p class="px-4 text-center">
          {{ entityOneName }}
        </p>
      </div>
      <Label>-</Label>
      <div class=" bg-muted/30 hover:bg-muted/40 rounded-lg py-2 hover:cursor-pointer"
           @click="navigateToEntity(selectedEntity2)">
        <p class="px-4 text-center">
          {{ entityTwoName }}
        </p>
      </div>
    </CardContent>
  </Card>
</template>

<style scoped>
.text-red-500 {
  color: red;
}
</style>