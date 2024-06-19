<script setup lang="ts">
import {FolderOpenDot } from 'lucide-vue-next'
import ThemePicker from "@/layouts/components/ThemePicker.vue";

</script>

<template>
  <div class="min-h-screen flex flex-col">
    <header class="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
      <div class="w-full flex flex-row">
      </div>
      <ThemePicker/>
    </header>
    <div class="flex flex-1 h-full w-full">
      <div class="grid md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr] w-full">
        <!--  Menu  -->
        <div class="border-r bg-muted/40 flex flex-col">
          <div class="flex-1">
            <nav class="grid items-start px-2 text-sm font-medium lg:px-4 pt-8">
              <a href="/dashboard">
                <div
                    class="flex items-center gap-3 rounded-lg px-3 py-2 text-primary transition-all hover:text-primary bg-muted">
                  <FolderOpenDot  class="h-4 w-4"/>
                  Stories
                </div>
              </a>
            </nav>
          </div>
        </div>
        <main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6 w-full h-full">
          <slot/>
        </main>
      </div>
    </div>
  </div>
</template>
<style scoped>

</style>