<script setup lang="ts">
import type { WorldInfo } from '../../../bindings/litguardian/internal'
import { useWorldStore } from '@/stores/worldStore'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Button } from '@/components/ui/button'
import { MoreHorizontal } from 'lucide-vue-next'
import { timeAgo } from '@/lib/timeago'

const props = defineProps<{ world: WorldInfo }>()
const store = useWorldStore()
</script>

<template>
  <div
    class="group flex items-center gap-2 rounded-lg px-3 py-3 hover:bg-accent cursor-pointer transition-colors"
    @click="store.openWorld(props.world.path)"
  >
    <!-- World info -->
    <div class="min-w-0 flex-1">
      <p class="font-medium text-sm leading-tight truncate">{{ props.world.name }}</p>
      <Tooltip>
        <TooltipTrigger as-child>
          <p class="text-xs text-muted-foreground truncate mt-0.5">{{ props.world.path }}</p>
        </TooltipTrigger>
        <TooltipContent side="bottom" align="start">
          <p>{{ props.world.path }}</p>
        </TooltipContent>
      </Tooltip>
      <p class="text-xs text-muted-foreground/60 mt-0.5">{{ timeAgo(props.world.lastOpened) }}</p>
    </div>

    <!-- ⋯ dropdown -->
    <DropdownMenu>
      <DropdownMenuTrigger as-child @click.stop>
        <Button
          variant="ghost"
          size="icon"
          class="h-7 w-7 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <MoreHorizontal class="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" @click.stop>
        <DropdownMenuItem
          class="text-destructive focus:text-destructive cursor-pointer"
          @click="store.removeRecentWorld(props.world.path)"
        >
          Remove from recents
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  </div>
</template>
