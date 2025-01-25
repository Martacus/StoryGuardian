<!-- frontend/src/components/entity/EntityTagList.vue -->
<script setup lang="ts">
import TagListBase from "@/components/shared/TagListBase.vue";
import {StoryModule} from "../../../bindings/storyguardian/src/project";
import TextTooltip from "@/components/ui/tooltip/TextTooltip.vue";
import {Dialog, DialogContent, DialogTrigger} from "@/components/ui/dialog";
import IconButton from "@/components/ui/button/IconButton.vue";
import {Plus} from "lucide-vue-next";
import {ref} from "vue";

const props = defineProps<{
  tags: string[],
  moduleConfig: StoryModule
}>();
const emit = defineEmits(['configChange'])

const dialogOpen = ref(false);
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

          <p>EntityTagList specific content</p>
        </DialogContent>
      </Dialog>
    </template>
  </TagListBase>
</template>