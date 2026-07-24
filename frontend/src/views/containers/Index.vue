<template>
  <div>
    <PageHeader :title="t('menu.containers')">
      <template #extra>
        <a-radio-group v-model:value="filter" button-style="solid" @change="loadContainers">
          <a-radio-button value="">{{ t('containers.filterAll') }}</a-radio-button>
          <a-radio-button value="running">{{ t('containers.filterRunning') }}</a-radio-button>
          <a-radio-button value="exited">{{ t('containers.filterStopped') }}</a-radio-button>
        </a-radio-group>
        <a-button style="margin-left: 16px" type="primary" @click="showCreateModal = true"><PlusOutlined /> {{ t('common.create') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="containers" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <router-link :to="{ path: '/containers', query: { id: record.id } }">{{ record.name }}</router-link>
          </template>
          <template v-else-if="column.key === 'status'"><StatusBadge :status="record.state" :text="record.status" /></template>
          <template v-else-if="column.key === 'ports'">
            <span v-for="p in record.ports" :key="p.private_port">{{ p.public_port }}:{{ p.private_port }}<br /></span>
          </template>
          <template v-else-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button v-if="record.state !== 'running'" type="link" size="small" @click="handleStart(record)"><PlayCircleOutlined /> {{ t('common.start') }}</a-button>
              <a-button v-if="record.state === 'running'" type="link" size="small" @click="handleStop(record)"><PauseCircleOutlined /> {{ t('common.stop') }}</a-button>
              <a-button type="link" size="small" @click="handleRestart(record)"><ReloadOutlined /> {{ t('common.restart') }}</a-button>
              <a-button type="link" size="small" @click="$router.push('/containers/logs?id=' + record.id)">Logs</a-button>
              <a-button type="link" size="small" @click="$router.push('/containers/stats?id=' + record.id)">Stats</a-button>
              <a-popconfirm :title="t('containers.removeConfirm')" @confirm="handleRemove(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
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
import { listContainers, startContainer, stopContainer, restartContainer, removeContainer } from '@/api/container'
import type { ContainerInfo } from '@/api/container'
import { formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { PlusOutlined, PlayCircleOutlined, PauseCircleOutlined, ReloadOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const containers = ref<ContainerInfo[]>([])
const filter = ref('')
const showCreateModal = ref(false)

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 120, ellipsis: true },
  { title: t('containers.id'), dataIndex: 'name', key: 'name' },
  { title: t('containers.image'), dataIndex: 'image', key: 'image', ellipsis: true },
  { title: t('containers.state'), dataIndex: 'status', key: 'status', width: 120 },
  { title: t('containers.ports'), key: 'ports', width: 150 },
  { title: t('common.created_at'), key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 320 },
]

async function loadContainers() {
  loading.value = true
  try { const res = await listContainers({ status: filter.value || undefined }); containers.value = res.data.list } catch { /* ignore */ } finally { loading.value = false }
}

async function handleStart(c: ContainerInfo) { try { await startContainer(c.id); message.success(t('common.success')); loadContainers() } catch { /* handled */ } }
async function handleStop(c: ContainerInfo) { try { await stopContainer(c.id); message.success(t('common.success')); loadContainers() } catch { /* handled */ } }
async function handleRestart(c: ContainerInfo) { try { await restartContainer(c.id); message.success(t('common.success')); loadContainers() } catch { /* handled */ } }
async function handleRemove(c: ContainerInfo) { try { await removeContainer(c.id); message.success(t('common.success')); loadContainers() } catch { /* handled */ } }

onMounted(loadContainers)
</script>