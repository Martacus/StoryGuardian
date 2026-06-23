<script setup lang="ts">
import { provide } from 'vue'
import { useWorldStore } from '@/stores/worldStore'
import { useNavigationStore } from '@/stores/navigationStore'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import SidebarNavItem from './SidebarNavItem.vue'
import WorldOverview from '../world/WorldOverview.vue'
import EntitiesListView from '../entity/EntitiesListView.vue'
import EntityDetailView from '../entity/EntityDetailView.vue'
import RelationDetailView from '../entity/RelationDetailView.vue'
import RelationsListView from '../entity/RelationsListView.vue'
import CategoriesListView from '../category/CategoriesListView.vue'
import CategoryDetailView from '../category/CategoryDetailView.vue'
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
const navStore = useNavigationStore()

provide('navigate', (view: View) => {
  navStore.navigateTo({ view, entityId: null, linkId: null, categoryName: null })
})

const navItems: { id: View; icon: Component; label: string; disabled: boolean }[] = [
  { id: 'overview',    icon: Globe,      label: 'Overview',   disabled: false },
  { id: 'entities',   icon: Users,      label: 'Entities',   disabled: false },
  { id: 'links',      icon: Link,       label: 'Relations',  disabled: false },
  { id: 'tags',       icon: Tags,       label: 'Tags',       disabled: true  },
  { id: 'categories', icon: FolderTree, label: 'Categories', disabled: false },
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
          :active="navStore.currentView === item.id"
          :disabled="item.disabled"
          @select="navStore.navigateTo({ view: item.id, entityId: null, linkId: null, categoryName: null })"
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
        <WorldOverview v-if="navStore.currentView === 'overview'" />

        <template v-else-if="navStore.currentView === 'entities'">
          <RelationDetailView v-if="navStore.currentLinkId" />
          <EntityDetailView v-else-if="navStore.currentEntityId" />
          <EntitiesListView v-else />
        </template>

        <template v-else-if="navStore.currentView === 'links'">
          <RelationDetailView v-if="navStore.currentLinkId" />
          <RelationsListView v-else />
        </template>

        <template v-else-if="navStore.currentView === 'categories'">
          <CategoryDetailView v-if="navStore.currentCategoryName" />
          <CategoriesListView v-else />
        </template>

        <!-- Placeholder for future views -->
        <div
          v-else
          class="flex items-center justify-center h-full text-muted-foreground text-sm"
        >
          {{ navItems.find(n => n.id === navStore.currentView)?.label }} — coming soon
        </div>
      </main>
    </div>
  </div>
</template>
