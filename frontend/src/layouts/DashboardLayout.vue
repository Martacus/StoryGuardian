<script setup lang="ts">
import {FolderOpenDot,ShieldHalf, Tag, Map, History, LayoutDashboard   } from 'lucide-vue-next'
import ThemePicker from "@/layouts/components/ThemePicker.vue";
import MainMenuButton from "@/components/menu/MainMenuButton.vue";
import {useRouter} from "vue-router";

const router = useRouter();
</script>

<template>
  <div class="min-h-screen max-h-screen flex flex-col">
    <header class="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 lg:h-[60px] lg:px-6">
      <div class="w-full flex flex-row">
      </div>
      <ThemePicker/>
    </header>
    <div class="flex flex-1 overflow-hidden">
      <div class="grid md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr] w-full">
        <!--  Menu  -->
        <div class="border-r bg-muted/40 flex flex-col">
          <div class="flex-1 overflow-auto">
            <nav class="grid items-start px-2 text-sm font-medium lg:px-4 pt-8 gap-1">
              <MainMenuButton @click="router.push('/')">
                <FolderOpenDot class="h-4 w-4"/>
                Stories
              </MainMenuButton>
              <MainMenuButton @click="router.push('/dashboard')">
                <LayoutDashboard  class="h-4 w-4"/>
                Dashboard
              </MainMenuButton>
              <MainMenuButton  @click="router.push('/entities')">
                <ShieldHalf  class="h-4 w-4"/>
                Entities
              </MainMenuButton>
              <MainMenuButton  @click="router.push('/')">
                <Tag class="h-4 w-4"/>
                Tags
              </MainMenuButton>
              <MainMenuButton  @click="router.push('/')">
                <Map class="h-4 w-4"/>
                Maps
              </MainMenuButton>
              <MainMenuButton  @click="router.push('/')">
                <History  class="h-4 w-4"/>
                Timelines
              </MainMenuButton>
            </nav>
          </div>
        </div>
        <main class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 auto-rows-min gap-x-2 md:gap-x-4 lg:gap-x-6 gap-y-2 md:gap-y-4 lg:gap-y-6 p-2 md:p-4 lg:p-6 overflow-auto">
          <slot/>
        </main>
      </div>
    </div>
  </div>
</template>
<style scoped>
</style>