import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import { getProfile } from '@/api/user'
import type { UserInfo } from '@/env'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  async function login(username: string, password: string, totp_code?: string) {
    const res = await loginApi({ username, password, totp_code })
    token.value = res.data.access_token
    localStorage.setItem('token', res.data.access_token)
    if (res.data.refresh_token) {
      localStorage.setItem('refresh_token', res.data.refresh_token)
    }
  }

  async function logout() {
    try {
      await logoutApi()
    } catch { /* ignore */ }
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    router.push('/login')
  }

  async function getUserInfo() {
    const res = await getProfile()
    userInfo.value = res.data
  }

  return { token, userInfo, login, logout, getUserInfo }
})
