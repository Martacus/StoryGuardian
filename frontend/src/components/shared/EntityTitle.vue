<script setup lang="ts">
import { Pencil, Check } from 'lucide-vue-next';
import { ref, watch } from 'vue';
import { Input } from "@/components/ui/input";
import IconButton from "@/components/ui/button/IconButton.vue";

const editTitle = ref(false);
const emit = defineEmits(['saveTitle']);

const props = defineProps<{
  title: String
}>();

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
      <IconButton @click="activateEdit">
        <Pencil/>
      </IconButton>
    </div>
    <div class="flex flex-row gap-2 m-2 items-center" v-if="editTitle">
      <Input type="text" :placeholder="props.title" v-model="localTitle" class="min-w-64"/>
      <IconButton @click="save">
        <Check />
      </IconButton>
    </div>
  </div>
</template>