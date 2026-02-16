<script setup lang="ts">
import { ref } from 'vue'
import { useSupabaseAuth } from '@/composables/useSupabase'
import { 
  LayoutDashboard, 
  Users, 
  Map as MapIcon, 
  Clock, 
  Settings, 
  LogOut,
  Menu,
  ChevronDown
} from 'lucide-vue-next'
import Button from '~/components/ui/button.vue'
import DropdownMenu from '~/components/ui/dropdown-menu.vue'
import DropdownMenuTrigger from '~/components/ui/dropdown-menu-trigger.vue'
import DropdownMenuContent from '~/components/ui/dropdown-menu-content.vue'
import DropdownMenuItem from '~/components/ui/dropdown-menu-item.vue'
import Avatar from '~/components/ui/avatar.vue'
import AvatarFallback from '~/components/ui/avatar-fallback.vue'

const { user, signOut } = useSupabaseAuth()
const isSidebarOpen = ref(false)
const isDropdownOpen = ref(false)

const navItems = [
  { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { name: 'Entities', path: '/dashboard/entities', icon: Users },
  { name: 'Maps', path: '/dashboard/maps', icon: MapIcon },
  { name: 'Timeline', path: '/dashboard/timeline', icon: Clock },
  { name: 'Settings', path: '/dashboard/settings', icon: Settings },
]

async function handleSignOut() {
  await signOut()
  navigateTo('/')
}

const route = useRoute()
</script>

<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Mobile Header -->
    <div class="lg:hidden bg-[#1e3a5f] text-white p-4 flex items-center justify-between">
      <span class="text-xl font-bold">StoryGuardian</span>
      <Button variant="ghost" size="icon" @click="isSidebarOpen = true">
        <Menu class="h-6 w-6" />
      </Button>
    </div>
    
    <div class="flex">
      <!-- Sidebar -->
      <aside 
        class="fixed inset-y-0 left-0 z-50 w-64 bg-[#1e3a5f] text-white transform transition-transform duration-300 lg:relative lg:translate-x-0"
        :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
      >
        <div class="p-6">
          <h1 class="text-2xl font-bold mb-8">StoryGuardian</h1>
          
          <nav class="space-y-2">
            <NuxtLink
              v-for="item in navItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-3 px-4 py-3 rounded-lg transition-colors"
              :class="route.path === item.path ? 'bg-white/10' : 'hover:bg-white/5'"
              @click="isSidebarOpen = false"
            >
              <component :is="item.icon" class="h-5 w-5" />
              <span>{{ item.name }}</span>
            </NuxtLink>
          </nav>
        </div>
        
        <!-- User Menu -->
        <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-white/10">
          <DropdownMenu v-model:open="isDropdownOpen">
            <DropdownMenuTrigger class="w-full flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-white/5">
              <Avatar class="h-8 w-8 bg-orange-500">
                <AvatarFallback>{{ user?.email?.charAt(0).toUpperCase() || 'U' }}</AvatarFallback>
              </Avatar>
              <div class="flex-1 text-left">
                <p class="text-sm font-medium truncate">{{ user?.email || 'User' }}</p>
              </div>
              <ChevronDown class="h-4 w-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" class="w-56">
              <DropdownMenuItem @click="handleSignOut">
                <LogOut class="mr-2 h-4 w-4" />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </aside>
      
      <!-- Mobile Overlay -->
      <div 
        v-if="isSidebarOpen" 
        class="fixed inset-0 bg-black/50 z-40 lg:hidden"
        @click="isSidebarOpen = false"
      />
      
      <!-- Main Content -->
      <main class="flex-1 p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
