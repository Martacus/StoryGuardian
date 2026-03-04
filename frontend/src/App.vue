<script setup lang="ts">
import { useWorldStore } from '@/stores/worldStore'
import WelcomeScreen from './components/WelcomeScreen.vue'
import { Button } from '@/components/ui/button'

const worldStore = useWorldStore()
</script>

<template>
  <WelcomeScreen v-if="!worldStore.hasOpenWorld" />

  <!-- Placeholder until the main editor is built (Step 6+) -->
  <div
    v-else
    class="h-screen w-screen flex flex-col items-center justify-center gap-4 bg-background text-foreground"
  >
    <p class="text-lg font-medium">
      "{{ worldStore.currentWorld?.name }}" is open.
    </p>
    <Button variant="outline" @click="worldStore.closeWorld()">
      ← Close World
    </Button>
  </div>
</template>
