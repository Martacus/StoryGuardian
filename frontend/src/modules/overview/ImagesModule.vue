<script setup lang="ts">
import { ref, computed, watch, inject, h, onMounted, onBeforeUnmount, type ComputedRef, type ShallowRef, type Component } from 'vue'
import { Upload, LayoutGrid, List, Trash2, Check, ImageOff } from 'lucide-vue-next'
import { useImageStore } from '@/stores/imageStore'
import { useLayoutStore } from '@/stores/layoutStore'
import { useWorldStore } from '@/stores/worldStore'
import { ImageService } from '../../../bindings/litguardian/internal'
import { Button } from '@/components/ui/button'

const moduleActions = inject<ShallowRef<Component | null>>('moduleActions')!
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { timeAgo } from '@/lib/timeago'

const ctx = inject<ComputedRef<{ viewId: string; moduleId: string }>>('moduleContext')!

const imageStore = useImageStore()
const layoutStore = useLayoutStore()
const worldStore = useWorldStore()

// ── View mode (grid / list) ──────────────────────────────────────────────────

const viewMode = computed(() =>
  layoutStore.getModuleConfig(ctx.value.viewId, ctx.value.moduleId, 'viewMode', 'grid') as 'grid' | 'list'
)

async function setViewMode(mode: 'grid' | 'list') {
  layoutStore.setModuleConfig(ctx.value.viewId, ctx.value.moduleId, 'viewMode', mode)
  await layoutStore.saveLayout()
}

// ── Data URLs cache ──────────────────────────────────────────────────────────

const dataUrls = ref<Map<string, string>>(new Map())

async function refreshDataUrls() {
  if (!worldStore.currentWorld) return
  const path = worldStore.currentWorld.path
  const newMap = new Map<string, string>()
  await Promise.all(
    imageStore.images.map(async (img) => {
      try {
        const url = await ImageService.GetImageDataURL(path, img.fileName)
        if (url) newMap.set(img.fileName, url)
      } catch {
        // skip broken images
      }
    })
  )
  dataUrls.value = newMap
}

watch(() => imageStore.images, () => { refreshDataUrls() }, { immediate: true, deep: true })
onMounted(() => { refreshDataUrls() })

moduleActions.value = () => [
  h(Button, {
    variant: 'ghost',
    size: 'icon-sm',
    disabled: imageStore.loading,
    title: 'Import images',
    onClick: () => imageStore.importImages(),
  }, () => h(Upload, { class: 'h-3.5 w-3.5' })),
  h(Button, {
    variant: viewMode.value === 'grid' ? 'default' : 'ghost',
    size: 'icon-sm',
    title: 'Grid view',
    onClick: () => setViewMode('grid'),
  }, () => h(LayoutGrid, { class: 'h-3.5 w-3.5' })),
  h(Button, {
    variant: viewMode.value === 'list' ? 'default' : 'ghost',
    size: 'icon-sm',
    title: 'List view',
    onClick: () => setViewMode('list'),
  }, () => h(List, { class: 'h-3.5 w-3.5' })),
]
onBeforeUnmount(() => { moduleActions.value = null })

function getDataUrl(fileName: string): string {
  return dataUrls.value.get(fileName) ?? ''
}

// ── Delete confirmation ──────────────────────────────────────────────────────

const confirmingDelete = ref<string | null>(null)
let confirmTimer: ReturnType<typeof setTimeout> | null = null

function requestDelete(fileName: string) {
  if (confirmingDelete.value === fileName) {
    // Second click — confirmed
    if (confirmTimer) { clearTimeout(confirmTimer); confirmTimer = null }
    confirmingDelete.value = null
    imageStore.deleteImage(fileName)
  } else {
    // First click — arm the confirmation
    if (confirmTimer) clearTimeout(confirmTimer)
    confirmingDelete.value = fileName
    confirmTimer = setTimeout(() => {
      confirmingDelete.value = null
      confirmTimer = null
    }, 2000)
  }
}

// ── Formatting helpers ───────────────────────────────────────────────────────

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<template>
  <div class="flex flex-col h-full">

    <!-- ── Content ────────────────────────────────────────────────────────── -->
    <div class="flex-1 overflow-y-auto px-3 py-3">

      <!-- Empty state -->
      <div
        v-if="imageStore.images.length === 0 && !imageStore.loading"
        class="flex flex-col items-center justify-center h-full gap-2 text-muted-foreground"
      >
        <ImageOff class="h-8 w-8 opacity-40" />
        <p class="text-sm">No images yet. Click Import to add reference images.</p>
      </div>

      <!-- Grid mode -->
      <div
        v-else-if="viewMode === 'grid'"
        class="grid gap-3"
        style="grid-template-columns: repeat(auto-fill, minmax(120px, 1fr))"
      >
        <div
          v-for="img in imageStore.images"
          :key="img.fileName"
          class="group relative flex flex-col gap-1"
        >
          <!-- Thumbnail -->
          <div class="relative aspect-square overflow-hidden rounded-md bg-muted transition-transform group-hover:scale-[1.02] group-hover:shadow-md">
            <img
              v-if="getDataUrl(img.fileName)"
              :src="getDataUrl(img.fileName)"
              :alt="img.fileName"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full items-center justify-center">
              <ImageOff class="h-6 w-6 text-muted-foreground opacity-40" />
            </div>

            <!-- Delete overlay -->
            <button
              class="no-drag absolute top-1 right-1 flex h-6 w-6 items-center justify-center rounded bg-background/80 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-destructive hover:text-destructive-foreground"
              :title="confirmingDelete === img.fileName ? 'Click again to confirm' : 'Delete'"
              @click.stop="requestDelete(img.fileName)"
            >
              <Check v-if="confirmingDelete === img.fileName" class="h-3 w-3" />
              <Trash2 v-else class="h-3 w-3" />
            </button>
          </div>

          <!-- File name -->
          <p class="truncate text-xs text-muted-foreground leading-tight" :title="img.fileName">
            {{ img.fileName }}
          </p>
        </div>
      </div>

      <!-- List mode -->
      <div v-else class="flex flex-col gap-1">
        <TooltipProvider :delay-duration="400">
          <div
            v-for="img in imageStore.images"
            :key="img.fileName"
            class="group flex items-center gap-2 rounded-md px-2 py-1.5 hover:bg-muted/50 transition-colors"
          >
            <!-- Small thumbnail -->
            <Tooltip>
              <TooltipTrigger as-child>
                <div class="shrink-0 h-9 w-9 overflow-hidden rounded bg-muted">
                  <img
                    v-if="getDataUrl(img.fileName)"
                    :src="getDataUrl(img.fileName)"
                    :alt="img.fileName"
                    class="h-full w-full object-cover"
                    loading="lazy"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <ImageOff class="h-4 w-4 text-muted-foreground opacity-40" />
                  </div>
                </div>
              </TooltipTrigger>
              <TooltipContent side="right" class="p-1">
                <img
                  v-if="getDataUrl(img.fileName)"
                  :src="getDataUrl(img.fileName)"
                  :alt="img.fileName"
                  class="max-w-[300px] max-h-[300px] rounded object-contain"
                />
              </TooltipContent>
            </Tooltip>

            <!-- Name + meta -->
            <div class="flex-1 min-w-0">
              <p class="truncate text-sm leading-tight" :title="img.fileName">{{ img.fileName }}</p>
              <p class="text-xs text-muted-foreground">
                {{ formatBytes(img.sizeBytes) }} · {{ timeAgo(img.modTime) }}
              </p>
            </div>

            <!-- Delete button -->
            <button
              class="no-drag shrink-0 flex h-6 w-6 items-center justify-center rounded opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive/10 hover:text-destructive"
              :title="confirmingDelete === img.fileName ? 'Click again to confirm' : 'Delete'"
              @click="requestDelete(img.fileName)"
            >
              <Check v-if="confirmingDelete === img.fileName" class="h-3.5 w-3.5" />
              <Trash2 v-else class="h-3.5 w-3.5" />
            </button>
          </div>
        </TooltipProvider>
      </div>

    </div>
  </div>
</template>
