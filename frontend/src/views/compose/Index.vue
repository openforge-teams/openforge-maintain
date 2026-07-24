<template>
  <div>
    <PageHeader :title="t('menu.compose')" />
    <a-card>
      <a-upload-dragger name="file" :beforeUpload="handleUpload" :showUploadList="false" accept=".yml,.yaml">
        <p class="ant-upload-drag-icon"><InboxOutlined /></p>
        <p class="ant-upload-text">{{ t('compose.uploadFile') }}</p>
      </a-upload-dragger>
      <a-divider />
      <a-table :dataSource="projects" :columns="columns" rowKey="name" :loading="loading">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleUp(record)">{{ t('compose.start') }}</a-button>
              <a-popconfirm :title="Confirm?" @confirm="handleDown(record)">
                <a-button type="link" size="small" danger>{{ t('compose.down') }}</a-button>
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
import { composeUp, composeDown } from '@/api/container'
import PageHeader from '@/components/PageHeader.vue'
import { InboxOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const projects = ref<any[]>([])

const columns = [
  { title: t('compose.projectName'), dataIndex: 'name', key: 'name' },
  { title: t('common.status'), dataIndex: 'status', key: 'status' },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at' },
  { title: t('common.actions'), key: 'actions', width: 200 },
]

async function loadProjects() { loading.value = true; setTimeout(() => { loading.value = false }, 500) }
async function handleUp(p: any) { try { await composeUp(p.name, ''); message.success(t('common.success')); loadProjects() } catch { /* handled */ } }
async function handleDown(p: any) { try { await composeDown(p.name); message.success(t('common.success')); loadProjects() } catch { /* handled */ } }
function handleUpload(file: File) { message.info('Upload: ' + file.name); return false }

onMounted(loadProjects)
</script>
