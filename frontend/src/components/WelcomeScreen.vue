<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWorldStore } from '@/stores/worldStore'
import { WorldService } from '../../bindings/litguardian'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { TooltipProvider } from '@/components/ui/tooltip'
import { BookOpen, PlusCircle, FolderOpen } from 'lucide-vue-next'
import WorldListItem from './WorldListItem.vue'
import CreateWorldDialog from './CreateWorldDialog.vue'

const store = useWorldStore()
const showCreateDialog = ref(false)

onMounted(() => store.loadRecentWorlds())

async function openExisting() {
  const path = await WorldService.SelectFolder()
  if (path) await store.openWorld(path)
}
</script>

<template>
  <TooltipProvider>
    <CreateWorldDialog v-model:open="showCreateDialog" />

    <div class="h-screen w-screen flex flex-col bg-background text-foreground select-none overflow-hidden">


      <!-- Body: two-column -->
      <div class="flex flex-1 overflow-hidden">

        <!-- ── Left panel: Recent Worlds ── -->
        <aside class="flex flex-col w-72 border-r border-border shrink-0">
          <div class="px-4 pt-5 pb-2 shrink-0">
            <p class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
              Recent Worlds
            </p>
          </div>

          <ScrollArea class="flex-1">
            <div class="px-2 pb-2">
              <template v-if="store.recentWorlds.length > 0">
                <WorldListItem
                  v-for="world in store.recentWorlds"
                  :key="world.path"
                  :world="world"
                />
              </template>

              <!-- Empty state -->
              <div
                v-else
                class="flex flex-col items-center justify-center gap-2 py-16 text-muted-foreground"
              >
                <BookOpen class="h-8 w-8 opacity-25" />
                <p class="text-sm">No recent worlds</p>
              </div>
            </div>
          </ScrollArea>
        </aside>

        <Separator orientation="vertical" />

        <!-- ── Right panel: Actions ── -->
        <main class="flex flex-col flex-1 items-center justify-center gap-3 p-10">
          <div class="w-full max-w-xs flex flex-col gap-3">
            <p class="text-xs font-semibold text-muted-foreground uppercase tracking-wider text-center mb-1">
              Get Started
            </p>

            <Button class="w-full justify-start gap-3 h-11" @click="showCreateDialog = true">
              <PlusCircle class="h-5 w-5 shrink-0" />
              Create New World
            </Button>

            <Button
              variant="outline"
              class="w-full justify-start gap-3 h-11"
              @click="openExisting"
            >
              <FolderOpen class="h-5 w-5 shrink-0" />
              Open Existing Folder
            </Button>
          </div>

        </main>
      </div>
    </div>
  </TooltipProvider>
</template>
