<script setup lang="ts">
import { Pencil, Check } from 'lucide-vue-next';
import { ref, watch } from 'vue';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const editTitle = ref(false);
const emit = defineEmits(['saveTitle']);

const props = defineProps({
  title: String
});

const localTitle = ref(props.title);

watch(() => props.title, (newTitle) => {
  localTitle.value = newTitle;
});

function activateEdit() {
  editTitle.value = !editTitle.value;
}

async function save() {
  activateEdit();
  emit('saveTitle', localTitle.value);
}
</script>

<template>
  <div v-if="localTitle">
    <div class="flex flex-row gap-2 items-center m-auto" v-if="!editTitle">
      <p class="text-2xl leading-loose ml-2">{{ localTitle }}</p>
      <Button class="btn btn-secondary" variant="outline" size="icon" @click="activateEdit">
        <Pencil/>
      </Button>
    </div>
    <div class="flex flex-row gap-2 m-2 items-center" v-if="editTitle">
      <Input type="text" :placeholder="props.title" v-model="localTitle" class="min-w-64"/>
      <Button class="btn btn-secondary" variant="outline" size="icon" @click="save">
        <Check />
      </Button>
    </div>
  </div>
</template>