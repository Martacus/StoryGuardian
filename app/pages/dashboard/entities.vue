<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold text-[#1e3a5f] mb-6">Entities</h1>
    <p class="text-slate-600">Manage your characters, creatures, factions, and organizations.</p>
  </div>
</template>
