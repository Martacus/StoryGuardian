<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold text-[#1e3a5f] mb-6">Settings</h1>
    <p class="text-slate-600">Manage your account and application preferences.</p>
  </div>
</template>
