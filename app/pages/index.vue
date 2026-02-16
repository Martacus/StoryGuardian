<script setup lang="ts">
import { BookOpen, Users, Map as MapIcon, Share2, Sparkles, ArrowRight, Check } from 'lucide-vue-next'


const features = [
  {
    icon: BookOpen,
    title: 'Lore Management',
    description: 'Organize your entire universe in one place. Track histories, cultures, magic systems, and every detail that makes your world unique.',
    color: 'bg-blue-100 text-blue-600',
  },
  {
    icon: Users,
    title: 'Entity Management',
    description: 'Create and manage characters, creatures, factions, and organizations. Build deep relationships and track their stories across your narratives.',
    color: 'bg-orange-100 text-orange-600',
  },
  {
    icon: MapIcon,
    title: 'World Building',
    description: 'Design landscapes, cities, and territories. Visualize how different locations connect and influence each other in your story.',
    color: 'bg-emerald-100 text-emerald-600',
  },
  {
    icon: Share2,
    title: 'Collaborative Writing',
    description: 'Invite co-creators to build worlds together. Share specific sections or keep some lore private while collaborating on the big picture.',
    color: 'bg-purple-100 text-purple-600',
  },
]

const benefits = [
  'No more scattered notes and lost ideas',
  'Visual timelines and relationship maps',
  'Cross-reference everything automatically',
  'Export your world bible anytime',
]

definePageMeta({
  middleware: 'auth'
})
</script>

<template>
  <div>
    <nav class="flex items-center justify-between px-6 py-4 lg:px-12">
      <div class="flex items-center gap-2">
        <div class="h-10 w-10 bg-[#1e3a5f] rounded-xl flex items-center justify-center">
          <Sparkles class="h-6 w-6 text-white" />
        </div>
        <span class="text-xl font-bold text-[#1e3a5f]">StoryGuardian</span>
      </div>
      <div class="flex items-center gap-4">
        <Button variant="ghost" class="text-[#1e3a5f]" as-child>
          <NuxtLink to="/login">Sign in</NuxtLink>
        </Button>
        <Button class="bg-[#1e3a5f] hover:bg-[#1e3a5f]/90" as-child>
          <NuxtLink to="/signup">Get Started</NuxtLink>
        </Button>
      </div>
    </nav>

    <section class="px-6 py-20 lg:px-12 lg:py-32">
      <div class="max-w-6xl mx-auto">
        <div class="grid lg:grid-cols-2 gap-12 items-center">
          <div class="space-y-8">
            <div class="inline-flex items-center gap-2 px-4 py-2 bg-orange-100 text-orange-700 rounded-full text-sm font-medium">
              <Sparkles class="h-4 w-4" />
              <span>For storytellers, by storytellers</span>
            </div>
            <h1 class="text-5xl lg:text-7xl font-bold text-[#1e3a5f] leading-tight">
              Build Worlds That 
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-orange-500 to-amber-500">Live & Breathe</span>
            </h1>
            <p class="text-xl text-slate-600 leading-relaxed">
              The all-in-one worldbuilding platform for authors, game masters, and creative minds.
            </p>
            <div class="flex flex-wrap gap-4">
              <Button size="lg" class="bg-[#1e3a5f] hover:bg-[#1e3a5f]/90 text-lg px-8" as-child>
                <NuxtLink to="/signup">Start Creating Free<ArrowRight class="ml-2 h-5 w-5" /></NuxtLink>
              </Button>
              <Button size="lg" variant="outline" class="text-lg px-8 border-[#1e3a5f] text-[#1e3a5f]" as-child>
                <NuxtLink to="#features">See Features</NuxtLink>
              </Button>
            </div>
          </div>
          <div class="relative">
            <div class="absolute inset-0 bg-gradient-to-r from-orange-200 to-amber-200 rounded-3xl transform rotate-3 opacity-50"></div>
            <div class="relative bg-white rounded-3xl shadow-2xl p-8 border border-slate-100">
              <div class="space-y-4">
                <div class="flex items-center gap-3 p-4 bg-slate-50 rounded-xl">
                  <div class="h-10 w-10 bg-[#1e3a5f] rounded-lg flex items-center justify-center">
                    <BookOpen class="h-5 w-5 text-white" />
                  </div>
                  <div>
                    <p class="font-semibold text-[#1e3a5f]">The Crystal Kingdom</p>
                    <p class="text-sm text-slate-500">Fantasy • 245 entries</p>
                  </div>
                </div>
                <div class="flex items-center gap-3 p-4 bg-slate-50 rounded-xl">
                  <div class="h-10 w-10 bg-orange-500 rounded-lg flex items-center justify-center">
                    <Users class="h-5 w-5 text-white" />
                  </div>
                  <div>
                    <p class="font-semibold text-[#1e3a5f]">King Aldric III</p>
                    <p class="text-sm text-slate-500">Character • 12 connections</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="features" class="px-6 py-20 lg:px-12 bg-white">
      <div class="max-w-6xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-[#1e3a5f] mb-4">Everything You Need to Build</h2>
          <p class="text-xl text-slate-600 max-w-2xl mx-auto">Powerful tools designed specifically for worldbuilders.</p>
        </div>
        <div class="grid md:grid-cols-2 gap-8">
          <Card v-for="feature in features" :key="feature.title" class="group hover:shadow-xl transition-shadow border-slate-100">
            <CardContent class="p-8">
              <div :class="`h-14 w-14 ${feature.color} rounded-2xl flex items-center justify-center mb-6`">
                <component :is="feature.icon" class="h-7 w-7" />
              </div>
              <h3 class="text-2xl font-bold text-[#1e3a5f] mb-3">{{ feature.title }}</h3>
              <p class="text-slate-600">{{ feature.description }}</p>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>

    <footer class="px-6 py-12 lg:px-12 border-t border-slate-200 bg-white">
      <div class="max-w-6xl mx-auto text-center">
        <p class="text-slate-500">© 2025 StoryGuardian. Built for storytellers.</p>
      </div>
    </footer>
  </div>
</template>
