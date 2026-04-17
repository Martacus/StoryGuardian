<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useEntityStore } from '@/stores/entityStore'
import { useNavigationStore } from '@/stores/navigationStore'
import ModuleGrid from '@/components/modules/ModuleGrid.vue'

const entityStore = useEntityStore()
const navStore = useNavigationStore()
const entityName = computed(() => entityStore.selectedEntity?.name ?? 'Entity')
</script>

<template>
  <ModuleGrid view-id="entity-detail">
    <template #title>
      <div class="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          class="h-7 w-7"
          @click="navStore.goBack()"
        >
          <ArrowLeft class="h-4 w-4" />
        </Button>
        {{ entityName }}
      </div>
    </template>
  </ModuleGrid>
</template>
