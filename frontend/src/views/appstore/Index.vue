<template>
  <div>
    <PageHeader :title="t('appstore.title')" />
    <a-tabs v-model:activeKey="activeTab" @change="loadApps">
      <a-tab-pane key="store" :tab="t('appstore.title')">
        <a-radio-group v-model:value="category" button-style="solid" style="margin-bottom: 16px" @change="loadApps">
          <a-radio-button value="">{{ t('appstore.all') }}</a-radio-button>
          <a-radio-button value="website">{{ t('appstore.website') }}</a-radio-button>
          <a-radio-button value="database">{{ t('appstore.database') }}</a-radio-button>
          <a-radio-button value="tools">{{ t('appstore.tools') }}</a-radio-button>
          <a-radio-button value="dev">{{ t('appstore.dev') }}</a-radio-button>
          <a-radio-button value="ai">{{ t('appstore.ai') }}</a-radio-button>
        </a-radio-group>
        <a-row :gutter="[16, 16]">
          <a-col v-for="app in apps" :key="app.id" :xs="24" :sm="12" :md="8" :lg="6">
            <a-card hoverable @click="$router.push('/appstore/' + app.id)">
              <a-card-meta :title="app.name" :description="app.description">
                <template #avatar><AppstoreOutlined style="font-size: 32px; color: #1890ff" /></template>
              </a-card-meta>
              <div style="margin-top: 12px"><a-tag>{{ app.category }}</a-tag><a-tag>{{ app.version }}</a-tag></div>
            </a-card>
          </a-col>
        </a-row>
      </a-tab-pane>
      <a-tab-pane key="installed" :tab="t('appstore.installed')">
        <a-table :dataSource="installedApps" :columns="installedColumns" rowKey="id" :loading="installedLoading">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'"><StatusBadge :status="record.status" :text="record.status" /></template>
            <template v-else-if="column.key === 'actions'"><a-space>
              <a-button type="link" size="small" @click="handleUpgrade(record)">{{ t('appstore.upgrade') }}</a-button>
              <a-popconfirm :title="t('appstore.uninstall') + '?'" @confirm="handleUninstall(record)"><a-button type="link" size="small" danger>{{ t('appstore.uninstall') }}</a-button></a-popconfirm>
            </a-space></template>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { getAppList, getInstalledApps, uninstallApp, upgradeApp } from '@/api/appstore'
import type { AppInfo, InstalledApp } from '@/api/appstore'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { AppstoreOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const activeTab = ref('store')
const category = ref('')
const apps = ref<AppInfo[]>([])
const installedApps = ref<InstalledApp[]>([])
const installedLoading = ref(false)

const installedColumns = [
  { title: t('common.name'), dataIndex: 'name', key: 'name' },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: 'Version', dataIndex: 'version', key: 'version', width: 100 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 180 },
]

async function loadApps() {
  if (activeTab.value === 'store') {
    try { const res = await getAppList({ category: category.value || undefined }); apps.value = res.data.list } catch {}
  } else {
    installedLoading.value = true
    try { const res = await getInstalledApps(); installedApps.value = res.data.list } catch {} finally { installedLoading.value = false }
  }
}
async function handleUpgrade(app: InstalledApp) { try { await upgradeApp(app.id); message.success(t('common.success')); loadApps() } catch {} }
async function handleUninstall(app: InstalledApp) { try { await uninstallApp(app.id); message.success(t('common.success')); loadApps() } catch {} }

onMounted(loadApps)
</script>
