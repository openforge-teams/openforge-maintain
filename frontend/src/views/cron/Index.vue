<template>
  <div>
    <PageHeader :title="t('menu.cron')">
      <template #extra>
        <a-button type="primary" @click="openEditModal()"><PlusOutlined /> {{ t('cron.createJob') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="jobs" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-switch :checked="record.status === 'enabled'" @change="handleToggle(record)" />
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleRun(record)">{{ t('cron.runNow') }}</a-button>
              <a-button type="link" size="small" @click="openEditModal(record)">{{ t('common.edit') }}</a-button>
              <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showModal" :title="editing ? t('cron.editJob') : t('cron.createJob')" @ok="handleSave" :confirmLoading="saving">
      <a-form :model="form" layout="vertical">
        <a-form-item :label="t('cron.jobName')"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item :label="t('cron.spec')"><a-input v-model:value="form.spec" placeholder="* * * * *" /></a-form-item>
        <a-form-item :label="t('cron.command')"><a-input v-model:value="form.command" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listCronJobs, createCronJob, updateCronJob, deleteCronJob, startCronJob, stopCronJob, runCronJob } from '@/api/cron'
import type { CronJob } from '@/api/cron'
import PageHeader from '@/components/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const jobs = ref<CronJob[]>([])
const showModal = ref(false)
const editing = ref<CronJob | null>(null)
const form = reactive({ name: '', spec: '', command: '' })

const columns = [
  { title: t('cron.jobName'), dataIndex: 'name', key: 'name' },
  { title: t('cron.spec'), dataIndex: 'spec', key: 'spec', width: 120 },
  { title: t('cron.command'), dataIndex: 'command', key: 'command', ellipsis: true },
  { title: t('common.status'), key: 'status', width: 80 },
  { title: t('cron.lastRun'), dataIndex: 'last_run_at', key: 'last_run_at', width: 180 },
  { title: t('cron.nextRun'), dataIndex: 'next_run_at', key: 'next_run_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 250 },
]

async function loadJobs() { loading.value = true; try { const res = await listCronJobs(); jobs.value = res.data.list } catch {} finally { loading.value = false } }
function openEditModal(job?: CronJob) {
  if (job) { editing.value = job; form.name = job.name; form.spec = job.spec; form.command = job.command }
  else { editing.value = null; form.name = ''; form.spec = ''; form.command = '' }
  showModal.value = true
}
async function handleSave() {
  saving.value = true
  try {
    if (editing.value) { await updateCronJob(editing.value.id, form) } else { await createCronJob(form) }
    message.success(t('common.success')); showModal.value = false; loadJobs()
  } catch {} finally { saving.value = false }
}
async function handleToggle(job: CronJob) {
  try { if (job.status === 'enabled') { await stopCronJob(job.id) } else { await startCronJob(job.id) }; message.success(t('common.success')); loadJobs() } catch {}
}
async function handleRun(job: CronJob) { try { await runCronJob(job.id); message.success(t('common.success')) } catch {} }
async function handleDelete(job: CronJob) { try { await deleteCronJob(job.id); message.success(t('common.success')); loadJobs() } catch {} }

onMounted(loadJobs)
</script>
