import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true')
  const theme = ref<'light' | 'dark'>(localStorage.getItem('theme') as 'light' | 'dark' || 'light')
  const language = ref(localStorage.getItem('language') || 'zh-CN')

  watch(sidebarCollapsed, (val) => {
    localStorage.setItem('sidebarCollapsed', String(val))
  })

  watch(theme, (val) => {
    localStorage.setItem('theme', val)
  })

  watch(language, (val) => {
    localStorage.setItem('language', val)
  })

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setTheme(val: 'light' | 'dark') {
    theme.value = val
  }

  function setLanguage(val: string) {
    language.value = val
  }

  return { sidebarCollapsed, theme, language, toggleSidebar, setTheme, setLanguage }
})
