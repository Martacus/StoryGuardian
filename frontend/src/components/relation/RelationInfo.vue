<script setup lang="ts">

import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
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
import {ref} from "vue";
import EntitySelector from "@/components/shared/EntitySelector.vue";

const props = defineProps<{
  moduleConfig: StoryModule
}>()

const emit = defineEmits(['configChange'])
const infoDialogOpen = ref(false);

const {showCardBody, toggleCardBody} = useToggleBody(props.moduleConfig)
const {columnSize, changeGridSize } = useGridSize(props.moduleConfig)

</script>

<template>
  <Card class="bg-muted/30 min-w-[22rem]" :class="columnSize">
    <CardHeader class="flex flex-row justify-between items-center">
      <CardTitle>Information</CardTitle>
      <div class="flex flex-row space-x-2">

        <!--    Edit Relation Info Modal    -->
        <Dialog v-model:open="infoDialogOpen" v-if="showCardBody">
          <DialogTrigger>
            <TextTooltip text="Edit Relation">
              <IconButton @click="">
                <Edit />
              </IconButton>
            </TextTooltip>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Edit Relation</DialogTitle>
            </DialogHeader>

            <EntitySelector></EntitySelector>

            <DialogFooter>
              <Button class="w-full">
                Create
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <VerticalSeperator />
        <GridSizeSelector v-if="moduleConfig" :column-size="moduleConfig.configuration['columnSize']" @update-grid-size="(newSize) => changeGridSize('relationInfo', newSize, emit)"/>
        <CardBodyToggler :show-card-body="showCardBody"  @toggle="toggleCardBody('relationInfo', emit)"/>
      </div>

    </CardHeader>
    <CardContent v-if="showCardBody">
      <Label>
        Zac - Agorass
      </Label>
    </CardContent>
  </Card>
</template>

<style scoped>

</style>