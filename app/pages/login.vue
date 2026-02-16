<script setup lang="ts">
import { useSupabaseAuth } from '@/composables/useSupabase' 

const { signIn } = useSupabaseAuth()
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  loading.value = true
  error.value = ''
  try {
    await signIn(email.value, password.value)
    navigateTo('/dashboard')
  } catch (e: any) {
    error.value = e.message || 'Failed to sign in'
  } finally {
    loading.value = false
  }
}

definePageMeta({
  layout: 'default'
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <Card class="w-full max-w-md">
      <CardContent class="p-8">
        <div class="text-center mb-8">
          <h1 class="text-2xl font-bold text-[#1e3a5f]">Welcome Back</h1>
          <p class="text-slate-600 mt-2">Sign in to continue building your worlds</p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <Label for="email">Email</Label>
            <Input id="email" v-model="email" type="email" placeholder="you@example.com" required />
          </div>
          
          <div>
            <Label for="password">Password</Label>
            <Input id="password" v-model="password" type="password" placeholder="••••••••" required />
          </div>

          <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

          <Button type="submit" class="w-full bg-[#1e3a5f]" :disabled="loading">
            {{ loading ? 'Signing in...' : 'Sign In' }}
          </Button>
        </form>

        <p class="text-center mt-6 text-sm text-slate-600">
          Don't have an account?
          <NuxtLink to="/signup" class="text-[#1e3a5f] font-medium hover:underline">Sign up</NuxtLink>
        </p>

        <div class="mt-6 text-center">
          <NuxtLink to="/" class="text-sm text-slate-500 hover:text-[#1e3a5f]">
            ← Back to home
          </NuxtLink>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
