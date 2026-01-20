import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

const tokenKey = 'access_token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(tokenKey) || '')
  const isAuthenticated = computed(() => token.value !== '')

  function setToken(value) {
    token.value = value
    localStorage.setItem(tokenKey, value)
  }

  function clearToken() {
    token.value = ''
    localStorage.removeItem(tokenKey)
    document.cookie = `${tokenKey}=; Max-Age=0; path=/`
  }

  return { token, isAuthenticated, setToken, clearToken }
})
