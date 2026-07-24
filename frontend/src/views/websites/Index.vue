<template>
  <div>
    <PageHeader :title="t('menu.websites')">
      <template #extra>
        <a-button type="primary" @click="showCreateModal = true">
          <PlusOutlined /> {{ t('websites.createWebsite') }}
        </a-button>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="websites" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusBadge :status="record.status" :text="record.status" />
          </template>
          <template v-else-if="column.key === 'ssl'">
            <a-tag :color="record.ssl_enabled ? 'green' : 'default'">SSL</a-tag>
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button v-if="!record.ssl_enabled" type="link" size="small" @click="handleEnableSSL(record)">
                {{ t('websites.enableSSL') }}
              </a-button>
              <a-button v-else type="link" size="small" @click="handleDisableSSL(record)">
                {{ t('websites.disableSSL') }}
              </a-button>
              <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showCreateModal" :title="t('websites.createWebsite')" @ok="handleCreate" :confirmLoading="creating">
      <a-form :model="createForm" layout="vertical">
        <a-form-item :label="t('websites.domain')">
          <a-input v-model:value="createForm.domain" />
        </a-form-item>
        <a-form-item :label="t('websites.websiteType')">
          <a-select v-model:value="createForm.type">
            <a-select-option value="static">{{ t('websites.static') }}</a-select-option>
            <a-select-option value="reverse_proxy">{{ t('websites.reverseProxy') }}</a-select-option>
            <a-select-option value="php">{{ t('websites.php') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="createForm.type === 'reverse_proxy'" :label="t('websites.proxyTarget')">
          <a-input v-model:value="createForm.proxy_target" />
        </a-form-item>
        <a-form-item v-if="createForm.type !== 'reverse_proxy'" :label="t('websites.rootPath')">
          <a-input v-model:value="createForm.root_path" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listWebsites, createWebsite, deleteWebsite, enableSSL, disableSSL } from '@/api/website'
import type { Website } from '@/api/website'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const creating = ref(false)
const websites = ref<Website[]>([])
const showCreateModal = ref(false)
const createForm = reactive({ domain: '', type: 'static', proxy_target: '', root_path: '' })

const columns = [
  { title: t('websites.domain'), dataIndex: 'domain', key: 'domain' },
  { title: t('websites.websiteType'), dataIndex: 'type', key: 'type', width: 120 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('websites.ssl'), key: 'ssl', width: 80 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 200 },
]

async function loadWebsites() {
  loading.value = true
  try { const res = await listWebsites(); websites.value = res.data.list } catch {} finally { loading.value = false }
}
async function handleCreate() {
  creating.value = true
  try { await createWebsite(createForm); message.success(t('common.success')); showCreateModal.value = false; loadWebsites() } catch {} finally { creating.value = false }
}
async function handleDelete(w: Website) { try { await deleteWebsite(w.id); message.success(t('common.success')); loadWebsites() } catch {} }
async function handleEnableSSL(w: Website) { try { await enableSSL(w.id); message.success(t('common.success')); loadWebsites() } catch {} }
async function handleDisableSSL(w: Website) { try { await disableSSL(w.id); message.success(t('common.success')); loadWebsites() } catch {} }

onMounted(loadWebsites)
</script>
