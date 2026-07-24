<template>
  <div>
    <PageHeader :title="t('menu.volumes')" />
    <a-card>
      <a-table :dataSource="volumes" :columns="columns" rowKey="name" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'size'">{{ formatFileSize(record.size) }}</template>
          <template v-else-if="column.key === 'actions'">
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
import { listVolumes, removeVolume } from '@/api/container'
import type { VolumeInfo } from '@/api/container'
import { formatFileSize, formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const { t } = useI18n()
const loading = ref(false)
const volumes = ref<VolumeInfo[]>([])

const columns = [
  { title: t('common.name'), dataIndex: 'name', key: 'name' },
  { title: t('volumes.driver'), dataIndex: 'driver', key: 'driver', width: 100 },
  { title: t('volumes.mountpoint'), dataIndex: 'mountpoint', key: 'mountpoint', ellipsis: true },
  { title: t('common.size'), key: 'size', width: 120 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadVolumes() {
  loading.value = true
  try { const res = await listVolumes(); volumes.value = res.data.list } catch { /* ignore */ } finally { loading.value = false }
}

async function handleRemove(v: VolumeInfo) {
  try { await removeVolume(v.name); message.success(t('common.success')); loadVolumes() } catch { /* handled */ }
}

onMounted(loadVolumes)
</script>
