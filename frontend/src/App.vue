<script setup lang="ts">
import { useWorldStore } from '@/stores/worldStore'
import { Toaster } from '@/components/ui/sonner'
import WelcomeScreen from './components/layout/WelcomeScreen.vue'
import DashboardLayout from './components/layout/DashboardLayout.vue'

const worldStore = useWorldStore()
</script>

<template>
  <!-- Global toast provider — renders on both WelcomeScreen and DashboardLayout.
       richColors  → red bg for errors, green for success (instantly scannable)
       closeButton → user can dismiss early
       duration    → 8 s auto-dismiss (enough time to read + copy the detail) -->
  <Toaster position="bottom-right" :duration="8000" richColors closeButton />
  <WelcomeScreen v-if="!worldStore.hasOpenWorld" />
  <DashboardLayout v-else />
</template>
