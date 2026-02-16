<script setup lang="ts">
import { BookOpen, Users, Map as MapIcon, Clock, TrendingUp } from 'lucide-vue-next'
import Button from '~/components/ui/button.vue'
import Card from '~/components/ui/card.vue'
import CardContent from '~/components/ui/card-content.vue'
import Badge from '~/components/ui/badge.vue'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const stats = [
  { label: 'Worlds', value: '3', icon: BookOpen, change: '+1 this week' },
  { label: 'Entities', value: '47', icon: Users, change: '+12 this week' },
  { label: 'Locations', value: '23', icon: MapIcon, change: '+5 this week' },
  { label: 'Timeline Events', value: '156', icon: Clock, change: '+28 this week' },
]

const recentActivity = [
  { action: 'Created', item: 'The Crystal Kingdom', type: 'World', time: '2 hours ago' },
  { action: 'Updated', item: 'King Aldric III', type: 'Character', time: '5 hours ago' },
  { action: 'Added', item: 'Northern Mountains', type: 'Location', time: '1 day ago' },
  { action: 'Modified', item: 'The Great War', type: 'Event', time: '2 days ago' },
]
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-3xl font-bold text-[#1e3a5f]">Dashboard</h1>
      <p class="text-slate-600 mt-1">Welcome back! Here is what is happening in your worlds.</p>
    </div>

    <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card v-for="stat in stats" :key="stat.label" class="border-slate-200">
        <CardContent class="p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-slate-500">{{ stat.label }}</p>
              <p class="text-3xl font-bold text-[#1e3a5f] mt-1">{{ stat.value }}</p>
              <p class="text-xs text-green-600 mt-1 flex items-center gap-1">
                <TrendingUp class="h-3 w-3" />
                {{ stat.change }}
              </p>
            </div>
            <div class="h-12 w-12 bg-[#1e3a5f]/10 rounded-xl flex items-center justify-center">
              <component :is="stat.icon" class="h-6 w-6 text-[#1e3a5f]" />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid lg:grid-cols-3 gap-8">
      <Card class="lg:col-span-2 border-slate-200">
        <CardContent class="p-6">
          <h2 class="text-xl font-bold text-[#1e3a5f] mb-4">Recent Activity</h2>
          <div class="space-y-4">
            <div v-for="(activity, index) in recentActivity" :key="index" class="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
              <div class="flex items-center gap-4">
                <Badge :variant="activity.action === 'Created' ? 'default' : 'secondary'">
                  {{ activity.action }}
                </Badge>
                <div>
                  <p class="font-medium text-[#1e3a5f]">{{ activity.item }}</p>
                  <p class="text-sm text-slate-500">{{ activity.type }}</p>
                </div>
              </div>
              <span class="text-sm text-slate-400">{{ activity.time }}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card class="border-slate-200">
        <CardContent class="p-6">
          <h2 class="text-xl font-bold text-[#1e3a5f] mb-4">Quick Actions</h2>
          <div class="space-y-3">
            <Button class="w-full justify-start" variant="outline">
              <BookOpen class="mr-2 h-4 w-4" />
              Create New World
            </Button>
            <Button class="w-full justify-start" variant="outline">
              <Users class="mr-2 h-4 w-4" />
              Add Entity
            </Button>
            <Button class="w-full justify-start" variant="outline">
              <MapIcon class="mr-2 h-4 w-4" />
              Create Location
            </Button>
            <Button class="w-full justify-start" variant="outline">
              <Clock class="mr-2 h-4 w-4" />
              Add Timeline Event
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
