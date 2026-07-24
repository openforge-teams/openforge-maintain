<template>
  <div>
    <PageHeader :title="t('menu.backup')">
      <template #extra>
        <a-button type="primary" @click="showTaskModal = true"><PlusOutlined /> {{ t('backup.createTask') }}</a-button>
      </template>
    </PageHeader>
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="tasks" :tab="t('backup.taskName')">
        <a-table :dataSource="tasks" :columns="taskColumns" rowKey="id" :loading="loading">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'actions'">
              <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDeleteTask(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </template>
          </template>
        </a-table>
      </a-tab-pane>
      <a-tab-pane key="records" :tab="t('backup.records')">
        <a-table :dataSource="records" :columns="recordColumns" rowKey="id" :loading="recordsLoading">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="record.status === 'success' ? 'green' : record.status === 'failed' ? 'red' : 'blue'">{{ record.status }}</a-tag>
            </template>
            <template v-else-if="column.key === 'file_size'">{{ formatFileSize(record.file_size) }}</template>
            <template v-else-if="column.key === 'actions'">
              <a-popconfirm :title="t('backup.restoreConfirm')" @confirm="handleRestore(record)">
                <a-button type="link" size="small">{{ t('backup.restore') }}</a-button>
              </a-popconfirm>
            </template>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>
    <a-modal v-model:open="showTaskModal" :title="t('backup.createTask')" @ok="handleCreateTask" :confirmLoading="creating">
      <a-form :model="taskForm" layout="vertical">
        <a-form-item :label="t('backup.taskName')"><a-input v-model:value="taskForm.name" /></a-form-item>
        <a-form-item :label="t('backup.backupType')"><a-input v-model:value="taskForm.type" placeholder="full" /></a-form-item>
        <a-form-item :label="t('backup.target')"><a-input v-model:value="taskForm.target" /></a-form-item>
        <a-form-item :label="t('backup.schedule')"><a-input v-model:value="taskForm.schedule" placeholder="0 2 * * *" /></a-form-item>
        <a-form-item :label="t('backup.retention')"><a-input-number v-model:value="taskForm.retention" :min="1" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listBackupTasks, createBackupTask, deleteBackupTask, listBackups, restoreBackup } from '@/api/backup'
import type { BackupTask, BackupRecord } from '@/api/backup'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const recordsLoading = ref(false)
const creating = ref(false)
const activeTab = ref('tasks')
const tasks = ref<BackupTask[]>([])
const records = ref<BackupRecord[]>([])
const showTaskModal = ref(false)
const taskForm = reactive({ name: '', type: 'full', target: '', schedule: '0 2 * * *', retention: 7 })

const taskColumns = [
  { title: t('backup.taskName'), dataIndex: 'name', key: 'name' },
  { title: t('backup.backupType'), dataIndex: 'type', key: 'type', width: 100 },
  { title: t('backup.target'), dataIndex: 'target', key: 'target' },
  { title: t('backup.schedule'), dataIndex: 'schedule', key: 'schedule', width: 120 },
  { title: t('backup.retention'), dataIndex: 'retention', key: 'retention', width: 80 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

const recordColumns = [
  { title: t('backup.taskName'), dataIndex: 'task_name', key: 'task_name' },
  { title: t('common.status'), key: 'status', width: 80 },
  { title: t('backup.fileSize'), key: 'file_size', width: 120 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadTasks() { loading.value = true; try { const res = await listBackupTasks(); tasks.value = res.data.list } catch {} finally { loading.value = false } }
async function loadRecords() { recordsLoading.value = true; try { const res = await listBackups(); records.value = res.data.list } catch {} finally { recordsLoading.value = false } }
async function handleCreateTask() { creating.value = true; try { await createBackupTask(taskForm); message.success(t('common.success')); showTaskModal.value = false; loadTasks() } catch {} finally { creating.value = false } }
async function handleDeleteTask(task: BackupTask) { try { await deleteBackupTask(task.id); message.success(t('common.success')); loadTasks() } catch {} }
async function handleRestore(record: BackupRecord) { try { await restoreBackup(record.id); message.success(t('common.success')) } catch {} }

watch(activeTab, (v) => { if (v === 'tasks') loadTasks(); else loadRecords() })
onMounted(loadTasks)
</script>
