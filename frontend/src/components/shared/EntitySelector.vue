<script setup lang="ts">
import { Button } from '@/components/ui/button'

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import {onMounted, ref} from 'vue'
import {GetEntities} from "../../../bindings/storyguardian/src/project/entitymanager";

const entities = ref<{ value: string; label: string }[]>([]);

const open = ref(false)
const selectedEntity = ref('')

onMounted(async () => {
  const retrievedEntities = await GetEntities();
  retrievedEntities.forEach(e => {
    if(e !== null){
      entities.value.push({value: e.id.toString(), label: e.name});
    }
  });
});

async function setEntity(entityId: string){
  selectedEntity.value = entityId;
  open.value = false;
}
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button
          variant="outline"
          role="combobox"
          :aria-expanded="open"
          class="w-[200px] justify-between"
      >
        {{ selectedEntity ? entities.find((entity) => entity.value === selectedEntity)?.label : 'Select entity...' }}

        <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-[200px] p-0">
      <Command>
        <CommandInput placeholder="Search entity..." />
        <CommandEmpty>No entities found.</CommandEmpty>
        <CommandList>
          <CommandGroup>
            <CommandItem
                v-for="entity in entities"
                :key="entity.value"
                :value="entity.label"
                @select="setEntity(entity.value)"
            >
              <Check
                  :class="cn(
                  'mr-2 h-4 w-4',
                  selectedEntity === entity.value ? 'opacity-100' : 'opacity-0',
                )"
              />
              {{ entity.label }}
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>

<style scoped>

</style>