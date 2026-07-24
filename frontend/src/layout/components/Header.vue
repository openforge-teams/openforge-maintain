<template>
  <a-layout-header class="header">
    <div class="header-left">
      <MenuUnfoldOutlined v-if="appStore.sidebarCollapsed" class="trigger" @click="appStore.toggleSidebar()" />
      <MenuFoldOutlined v-else class="trigger" @click="appStore.toggleSidebar()" />
    </div>
    <div class="header-right">
      <a-dropdown>
        <span class="lang-switch">
          <GlobalOutlined /> {{ appStore.language === 'zh-CN' ? '中文' : 'EN' }}
        </span>
        <template #overlay>
          <a-menu @click="onLangChange">
            <a-menu-item key="zh-CN">中文</a-menu-item>
            <a-menu-item key="en-US">English</a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
      <a-dropdown>
        <span class="user-info">
          <a-avatar :size="28" style="background-color: #1890ff">
            {{ userStore.userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </a-avatar>
          <span class="username">{{ userStore.userInfo?.username || 'User' }}</span>
        </span>
        <template #overlay>
          <a-menu>
            <a-menu-item key="settings" @click="$router.push('/settings')">
              <SettingOutlined /> {{ t('menu.settings') }}
            </a-menu-item>
            <a-menu-divider />
            <a-menu-item key="logout" @click="handleLogout">
              <LogoutOutlined /> {{ t('common.logout') || 'Logout' }}
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  GlobalOutlined,
  SettingOutlined,
  LogoutOutlined,
} from '@ant-design/icons-vue'

const { t, locale } = useI18n()
const appStore = useAppStore()
const userStore = useUserStore()

onMounted(() => {
  userStore.getUserInfo()
})

function onLangChange({ key }: { key: string }) {
  appStore.setLanguage(key)
  locale.value = key
}

function handleLogout() {
  userStore.logout()
}
</script>

<style scoped>
.header {
  background: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  z-index: 9;
}
.header-left {
  display: flex;
  align-items: center;
}
.trigger {
  font-size: 18px;
  cursor: pointer;
  padding: 0 12px;
  transition: color 0.3s;
}
.trigger:hover {
  color: #1890ff;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
.lang-switch {
  cursor: pointer;
  font-size: 14px;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.username {
  font-size: 14px;
}
</style>
