<script setup lang="ts">

import { Pencil, Check } from 'lucide-vue-next';
import {ref} from "vue";
import {Button} from "@/components/ui/button";
import {Project} from "../../../bindings/storyguardian/project";
import {Input} from "@/components/ui/input";
import {SetProjectTitle} from "../../../bindings/storyguardian/project/projectmanager";


const editTitle = ref(false);

const props = defineProps({
  story: Project
})

function activateEdit() {
  editTitle.value = !editTitle.value;
  console.log(props.story?.name)
}
async function save() {
  activateEdit();
  if(props.story){
    SetProjectTitle(props.story?.id, props.story?.name)
  } else {
    //Alert
  }
}
</script>

<template>
  <div v-if="story">
    <div class="flex flex-row gap-2 items-center m-auto" v-if="!editTitle">
      <p class="text-2xl leading-loose ml-2">{{story.name}}</p>
      <Button class="btn btn-secondary" variant="outline" size="icon" @click="activateEdit()">
        <Pencil/>
      </Button>
    </div>
    <div class="flex flex-row gap-2 m-2 items-center w-full" v-if="editTitle">
      <Input type="text" :placeholder="story.name" v-model="story.name" class="max-w-64"/>
      <Button class="btn btn-secondary" variant="outline" size="icon" @click="save()">
        <Check />
      </Button>
    </div>
  </div>

</template>

<style scoped>

</style>