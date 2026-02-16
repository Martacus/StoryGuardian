import { createClient, SupabaseClient } from '@supabase/supabase-js'

let supabaseInstance: SupabaseClient | null = null

export const useSupabase = () => {
  const config = useRuntimeConfig()
  
  // Check if Supabase credentials are configured
  const supabaseUrl = config.public.supabaseUrl as string
  const supabaseKey = config.public.supabaseKey as string
  
  if (!supabaseUrl || supabaseUrl === 'your_supabase_url_here' || !supabaseKey || supabaseKey === 'your_supabase_anon_key_here') {
    console.warn('Supabase credentials not configured. Please set SUPABASE_URL and SUPABASE_ANON_KEY in your .env file')
    
    // Return a mock client for development that shows helpful errors
    return {
      supabase: null as any
    }
  }
  
  // Create singleton instance
  if (!supabaseInstance) {
    supabaseInstance = createClient(supabaseUrl, supabaseKey)
  }
  
  return {
    supabase: supabaseInstance
  }
}

export const useSupabaseAuth = () => {
  const { supabase } = useSupabase()
  const user = useState('user', () => null)
  
  const checkSupabase = () => {
    if (!supabase) {
      throw new Error('Supabase not configured. Please add your Supabase credentials to .env file')
    }
  }
  
  const signUp = async (email: string, password: string) => {
    checkSupabase()
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
    })
    if (error) throw error
    return data
  }
  
  const signIn = async (email: string, password: string) => {
    checkSupabase()
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    if (error) throw error
    user.value = data.user
    return data
  }
  
  const signOut = async () => {
    checkSupabase()
    const { error } = await supabase.auth.signOut()
    if (error) throw error
    user.value = null
  }
  
  const getUser = async () => {
    checkSupabase()
    const { data: { user: currentUser } } = await supabase.auth.getUser()
    user.value = currentUser
    return currentUser
  }
  
  return {
    user,
    signUp,
    signIn,
    signOut,
    getUser,
  }
}
