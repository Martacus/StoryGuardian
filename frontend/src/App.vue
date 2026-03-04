<script setup lang="ts">
import { useWorldStore } from '@/stores/worldStore'
import WelcomeScreen from './components/WelcomeScreen.vue'
import DashboardLayout from './components/DashboardLayout.vue'

const worldStore = useWorldStore()
</script>

<template>
  <WelcomeScreen v-if="!worldStore.hasOpenWorld" />
  <DashboardLayout v-else />
</template>
