<script setup lang="ts">
import { Bold, Italic, Underline } from 'lucide-vue-next';
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import {onMounted, ref, watch} from "vue";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group";

const editor = ref<Editor>();
const selectedStyles = ref([]);

const props = defineProps<{
  content: String
}>()

onMounted(() => {
  editor.value = new Editor({
    content: props.content,
    extensions: [
      StarterKit,
      Document,
      Paragraph,
      Text,
    ],
    editorProps: {
      attributes: {
        class: 'prose dark:prose-invert prose-sm sm:prose-base lg:prose-lg xl:prose-2xl m-5 focus:outline-none',
      },
    },
  })
})

defineExpose({
  getHTML: () => editor.value?.getHTML(),
})

watch(selectedStyles, (newValue, oldValue) => {
  console.log(newValue, oldValue);
});
</script>

<template>
  <ToggleGroup type="multiple" size="sm" variant="outline" v-if="editor">
    <ToggleGroupItem value="bold" aria-label="Toggle bold" @click="editor?.chain().focus().toggleBold().run()">
      <Bold class="h-4 w-4" />
    </ToggleGroupItem>
    <ToggleGroupItem value="italic" aria-label="Toggle italic" @click="editor?.chain().focus().toggleItalic().run()">
      <Italic class="h-4 w-4" />
    </ToggleGroupItem>
    <ToggleGroupItem value="underline" aria-label="Toggle underline"
                     @click="editor?.chain().focus().toggleStrike().run()">
      <Underline class="h-4 w-4" />
    </ToggleGroupItem>
  </ToggleGroup>

  <EditorContent :editor="editor" class="w-full min-h-20 h-full border-2 border-ga-border-primary mt-1" v-if="editor">


  </EditorContent>
</template>

