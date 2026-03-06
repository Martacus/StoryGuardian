<script setup lang="ts">
import { ref } from 'vue'
import { useWorldStore } from '@/stores/worldStore'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import SidebarNavItem from './SidebarNavItem.vue'
import WorldOverview from './WorldOverview.vue'
import {
  BookOpen,
  Globe,
  Users,
  Link,
  Tags,
  FolderTree,
  Search,
  LogOut,
} from 'lucide-vue-next'
import type { Component } from 'vue'

type View = 'overview' | 'entities' | 'links' | 'tags' | 'categories'

const store = useWorldStore()
const currentView = ref<View>('overview')

const navItems: { id: View; icon: Component; label: string; disabled: boolean }[] = [
  { id: 'overview',    icon: Globe,      label: 'Overview',   disabled: false },
  { id: 'entities',   icon: Users,      label: 'Entities',   disabled: true  },
  { id: 'links',      icon: Link,       label: 'Links',      disabled: true  },
  { id: 'tags',       icon: Tags,       label: 'Tags',       disabled: true  },
  { id: 'categories', icon: FolderTree, label: 'Categories', disabled: true  },
]
</script>

<template>
  <div class="h-screen w-screen flex bg-background text-foreground overflow-hidden select-none">

    <!-- ── Sidebar (full height) ──────────────────────────────────────────── -->
    <aside class="flex flex-col w-52 shrink-0 border-r border-sidebar-border bg-sidebar text-sidebar-foreground">

      <!-- Logo -->
      <div class="flex items-center gap-2.5 px-4 h-12 shrink-0 border-b border-sidebar-border">
        <BookOpen class="h-5 w-5 text-sidebar-primary" />
        <span class="font-semibold text-sm tracking-tight">LitGuardian</span>
      </div>

      <!-- Nav items -->
      <nav class="flex flex-col gap-1 p-2 flex-1 overflow-y-auto">
        <SidebarNavItem
          v-for="item in navItems"
          :key="item.id"
          :icon="item.icon"
          :label="item.label"
          :active="currentView === item.id"
          :disabled="item.disabled"
          @select="currentView = item.id"
        />
      </nav>

      <!-- Close world -->
      <div class="p-2 border-t border-sidebar-border shrink-0">
        <Button
          variant="ghost"
          class="w-full justify-start gap-3 text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
          @click="store.closeWorld()"
        >
          <LogOut class="h-4 w-4" />
          Close World
        </Button>
      </div>
    </aside>

    <!-- ── Right column ───────────────────────────────────────────────────── -->
    <div class="flex flex-col flex-1 overflow-hidden">

      <!-- Top nav -->
      <header class="flex items-center gap-3 px-4 h-12 border-b border-border shrink-0">
        <span class="text-sm font-medium truncate max-w-[200px]">
          {{ store.currentWorld?.name }}
        </span>
        <div class="ml-auto flex items-center gap-2">
          <Search class="h-4 w-4 text-muted-foreground shrink-0" />
          <Input
            class="h-8 w-56 bg-transparent border-none shadow-none focus-visible:ring-0 placeholder:text-muted-foreground/50 px-0"
            placeholder="Search… (coming soon)"
            disabled
          />
        </div>
      </header>

      <!-- Content area -->
      <main class="flex-1 overflow-auto">
        <WorldOverview v-if="currentView === 'overview'" />

        <!-- Placeholder for future views -->
        <div
          v-else
          class="flex items-center justify-center h-full text-muted-foreground text-sm"
        >
          {{ navItems.find(n => n.id === currentView)?.label }} — coming soon
        </div>
      </main>
    </div>
  </div>
</template>
