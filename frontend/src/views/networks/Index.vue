<template>
  <div>
    <PageHeader :title="t('menu.networks')" />
    <a-card>
      <a-table :dataSource="networks" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'actions'">
            <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleRemove(record)">
              <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listNetworks } from '@/api/container'
import type { NetworkInfo } from '@/api/container'
import { formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const { t } = useI18n()
const loading = ref(false)
const networks = ref<NetworkInfo[]>([])

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 120, ellipsis: true },
  { title: t('common.name'), dataIndex: 'name', key: 'name' },
  { title: t('networks.driver'), dataIndex: 'driver', key: 'driver', width: 100 },
  { title: t('networks.scope'), dataIndex: 'scope', key: 'scope', width: 100 },
  { title: t('networks.containers'), dataIndex: 'containers', key: 'containers', width: 100 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadNetworks() {
  loading.value = true
  try { const res = await listNetworks(); networks.value = res.data.list } catch { /* ignore */ } finally { loading.value = false }
}

function handleRemove(_n: NetworkInfo) { message.info('Not implemented') }

onMounted(loadNetworks)
</script>
