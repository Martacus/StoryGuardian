export default defineNuxtRouteMiddleware(async (to, from) => {
  // Skip auth check on server
  if (process.server) return
  
  const config = useRuntimeConfig()
  const hasSupabase = config.public.supabaseUrl && 
                      config.public.supabaseUrl !== 'your_supabase_url_here' &&
                      config.public.supabaseKey &&
                      config.public.supabaseKey !== 'your_supabase_anon_key_here'
  
  // If Supabase is not configured, skip auth checks for development
  if (!hasSupabase) {
    console.log('Supabase not configured - running in demo mode')
    return
  }
  
  try {
    const { getUser } = useSupabaseAuth()
    const user = await getUser()
    
    // If user is not logged in and trying to access dashboard, redirect to landing
    if (!user && to.path.startsWith('/dashboard')) {
      return navigateTo('/')
    }
    
    // If user is logged in and trying to access landing page, redirect to dashboard
    if (user && to.path === '/') {
      return navigateTo('/dashboard')
    }
  } catch (error) {
    console.error('Auth error:', error)
    // If there's an auth error and trying to access dashboard, redirect to landing
    if (to.path.startsWith('/dashboard')) {
      return navigateTo('/')
    }
  }
})
