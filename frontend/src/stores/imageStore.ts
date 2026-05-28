import { ref } from 'vue'
import { defineStore } from 'pinia'
import { ImageService } from '../../bindings/litguardian/internal'
import type { ImageInfo } from '../../bindings/litguardian/internal'
import { useWorldStore } from './worldStore'
import { useAppToast } from '@/composables/useAppToast'

export const useImageStore = defineStore('images', () => {
  const { errorToast } = useAppToast()

  // ── State ─────────────────────────────────────────────────────────────────
  const images = ref<ImageInfo[]>([])
  const loading = ref(false)

  // ── Actions ────────────────────────────────────────────────────────────────

  /** Load image list from disk for the open world. */
  async function loadImages() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      const result = await ImageService.ListImages(worldStore.currentWorld.path)
      images.value = result ?? []
    } catch (e) {
      errorToast('Failed to load images', String(e))
    } finally {
      loading.value = false
    }
  }

  /**
   * Opens the native file picker, then copies selected files into the world's
   * images/ directory. Reloads the image list afterwards.
   */
  async function importImages() {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    loading.value = true
    try {
      const paths = await ImageService.SelectImages()
      if (!paths || paths.length === 0) return
      await ImageService.ImportImages(worldStore.currentWorld.path, paths)
      await loadImages()
    } catch (e) {
      errorToast('Failed to import images', String(e))
    } finally {
      loading.value = false
    }
  }

  /** Delete a single image from disk and reload the list. */
  async function deleteImage(fileName: string) {
    const worldStore = useWorldStore()
    if (!worldStore.currentWorld) return
    try {
      await ImageService.DeleteImage(worldStore.currentWorld.path, fileName)
      await loadImages()
    } catch (e) {
      errorToast('Failed to delete image', String(e))
    }
  }

  /** Reset to blank state. Called when a world is closed. */
  function reset() {
    images.value = []
    loading.value = false
  }

  return {
    images,
    loading,
    loadImages,
    importImages,
    deleteImage,
    reset,
  }
})
