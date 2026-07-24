<template>
  <a-layout-sider
    v-model:collapsed="appStore.sidebarCollapsed"
    :trigger="null"
    collapsible
    :width="220"
    :collapsed-width="64"
    class="sidebar"
  >
    <div class="logo">
      <CloudServerOutlined style="font-size: 24px; color: #1890ff" />
      <span v-show="!appStore.sidebarCollapsed" class="logo-text">openforge-maintain</span>
    </div>
    <a-menu
      v-model:selectedKeys="selectedKeys"
      v-model:openKeys="openKeys"
      theme="dark"
      mode="inline"
      @click="onMenuClick"
    >
      <a-menu-item key="/dashboard">
        <DashboardOutlined />
        <span>{{ t('menu.dashboard') }}</span>
      </a-menu-item>
      <a-menu-item key="/host">
        <MonitorOutlined />
        <span>{{ t('menu.host') }}</span>
      </a-menu-item>
      <a-menu-item key="/terminal">
        <CodeOutlined />
        <span>{{ t('menu.terminal') }}</span>
      </a-menu-item>
      <a-menu-item key="/files">
        <FolderOutlined />
        <span>{{ t('menu.files') }}</span>
      </a-menu-item>
      <a-sub-menu key="docker">
        <template #title>
          <ContainerOutlined />
          <span>Docker</span>
        </template>
        <a-menu-item key="/containers">{{ t('menu.containers') }}</a-menu-item>
        <a-menu-item key="/images">{{ t('menu.images') }}</a-menu-item>
        <a-menu-item key="/volumes">{{ t('menu.volumes') }}</a-menu-item>
        <a-menu-item key="/networks">{{ t('menu.networks') }}</a-menu-item>
        <a-menu-item key="/compose">{{ t('menu.compose') }}</a-menu-item>
      </a-sub-menu>
      <a-menu-item key="/websites">
        <GlobalOutlined />
        <span>{{ t('menu.websites') }}</span>
      </a-menu-item>
      <a-sub-menu key="databases">
        <template #title>
          <DatabaseOutlined />
          <span>{{ t('menu.mysql') }}</span>
        </template>
        <a-menu-item key="/databases/mysql">MySQL</a-menu-item>
        <a-menu-item key="/databases/postgres">PostgreSQL</a-menu-item>
        <a-menu-item key="/databases/redis">Redis</a-menu-item>
      </a-sub-menu>
      <a-menu-item key="/appstore">
        <AppstoreOutlined />
        <span>{{ t('menu.appstore') }}</span>
      </a-menu-item>
      <a-menu-item key="/cron">
        <ClockCircleOutlined />
        <span>{{ t('menu.cron') }}</span>
      </a-menu-item>
      <a-menu-item key="/ssl">
        <SafetyCertificateOutlined />
        <span>{{ t('menu.ssl') }}</span>
      </a-menu-item>
      <a-menu-item key="/backup">
        <CloudUploadOutlined />
        <span>{{ t('menu.backup') }}</span>
      </a-menu-item>
      <a-menu-item key="/firewall">
        <FireOutlined />
        <span>{{ t('menu.firewall') }}</span>
      </a-menu-item>
      <a-menu-item key="/settings">
        <SettingOutlined />
        <span>{{ t('menu.settings') }}</span>
      </a-menu-item>
    </a-menu>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/store/app'
import {
  DashboardOutlined,
  MonitorOutlined,
  FolderOutlined,
  ContainerOutlined,
  GlobalOutlined,
  DatabaseOutlined,
  AppstoreOutlined,
  ClockCircleOutlined,
  SafetyCertificateOutlined,
  CloudUploadOutlined,
  FireOutlined,
  SettingOutlined,
  CodeOutlined,
  CloudServerOutlined,
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const appStore = useAppStore()

const selectedKeys = computed(() => [route.path])
const openKeys = computed(() => {
  const path = route.path
  const keys: string[] = []
  if (path.startsWith('/container') || path.startsWith('/image') || path.startsWith('/volume') || path.startsWith('/network') || path === '/compose') {
    keys.push('docker')
  }
  if (path.startsWith('/database')) {
    keys.push('databases')
  }
  return keys
})

function onMenuClick({ key }: { key: string }) {
  router.push(key)
}
</script>

<style scoped>
.sidebar {
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  z-index: 10;
}
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}
.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
}
</style>
