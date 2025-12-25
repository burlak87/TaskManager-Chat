import { defineNuxtPlugin } from '#app'
import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(async (nuxtApp) => {
  const authStore = useAuthStore()

  if (authStore.accessToken && !authStore.user) {
    try {
      await authStore.fetchProfile()
    } catch (error) {
      console.error('Error initializing auth:', error)
      authStore.logout()
    }
  }
})
