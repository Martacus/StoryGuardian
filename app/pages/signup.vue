<script setup lang="ts">
import { useSupabaseAuth } from '@/composables/useSupabase' 

const { signUp } = useSupabaseAuth()
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref(false)
const loading = ref(false)

async function handleSubmit() {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  
  loading.value = true
  error.value = ''
  try {
    await signUp(email.value, password.value)
    success.value = true
  } catch (e: any) {
    error.value = e.message || 'Failed to sign up'
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
          <h1 class="text-2xl font-bold text-[#1e3a5f]">Create Your Account</h1>
          <p class="text-slate-600 mt-2">Start building your worlds today</p>
        </div>

        <div v-if="success" class="text-center">
          <div class="h-16 w-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-[#1e3a5f] mb-2">Check Your Email</h2>
          <p class="text-slate-600 mb-6">We've sent you a confirmation link. Please check your email to complete your registration.</p>
          <Button class="w-full bg-[#1e3a5f]" as-child>
            <NuxtLink to="/login">Go to Login</NuxtLink>
          </Button>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <Label for="email">Email</Label>
            <Input id="email" v-model="email" type="email" placeholder="you@example.com" required />
          </div>
          
          <div>
            <Label for="password">Password</Label>
            <Input id="password" v-model="password" type="password" placeholder="••••••••" required />
          </div>

          <div>
            <Label for="confirmPassword">Confirm Password</Label>
            <Input id="confirmPassword" v-model="confirmPassword" type="password" placeholder="••••••••" required />
          </div>

          <p v-if="error" class="text-sm text-red-500">{{ error }}</p>

          <Button type="submit" class="w-full bg-[#1e3a5f]" :disabled="loading">
            {{ loading ? 'Creating account...' : 'Create Account' }}
          </Button>
        </form>

        <p v-if="!success" class="text-center mt-6 text-sm text-slate-600">
          Already have an account?
          <NuxtLink to="/login" class="text-[#1e3a5f] font-medium hover:underline">Sign in</NuxtLink>
        </p>

        <div v-if="!success" class="mt-6 text-center">
          <NuxtLink to="/" class="text-sm text-slate-500 hover:text-[#1e3a5f]">
            ← Back to home
          </NuxtLink>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
