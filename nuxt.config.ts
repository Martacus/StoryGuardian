// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',

  devtools: { enabled: true }, 
  shadcn: { 
    prefix: '', 
    componentDir: '@/components/ui'
  },
  components: [
    {
      path: '~/app/components',
      pathPrefix: false
    }
  ],

  modules: [
    '@nuxtjs/tailwindcss',
    'shadcn-nuxt',
    '@vueuse/nuxt'
  ], 

  tailwindcss: {
    cssPath: '~/assets/css/tailwind.css'
  },

  runtimeConfig: {
    public: {
      supabaseUrl: process.env.SUPABASE_URL,
      supabaseKey: process.env.SUPABASE_ANON_KEY
    }
  },

  imports: {
    dirs: ['composables/**']
  }
})
