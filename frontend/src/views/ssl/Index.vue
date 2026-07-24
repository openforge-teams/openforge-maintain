<template>
  <div>
    <PageHeader :title="t('menu.ssl')">
      <template #extra>
        <a-button type="primary" @click="showModal = true"><PlusOutlined /> {{ t('ssl.requestCert') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="certs" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'active' ? 'green' : record.status === 'expired' ? 'red' : 'blue'">{{ t('ssl.' + record.status) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'auto_renew'">
            <a-switch :checked="record.auto_renew" disabled />
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button v-if="record.status === 'expired'" type="link" size="small" @click="handleRenew(record)">{{ t('ssl.renewCert') }}</a-button>
              <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showModal" :title="t('ssl.requestCert')" @ok="handleRequest" :confirmLoading="requesting">
      <a-form :model="form" layout="vertical">
        <a-form-item :label="t('ssl.domain')"><a-input v-model:value="form.domain" /></a-form-item>
        <a-form-item :label="t('ssl.provider')">
          <a-select v-model:value="form.provider"><a-select-option value="letsencrypt">Let's Encrypt</a-select-option></a-select>
        </a-form-item>
        <a-form-item><a-switch v-model:checked="form.auto_renew" /> {{ t('ssl.autoRenew') }}</a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listCerts, requestCert, renewCert, deleteCert } from '@/api/ssl'
import type { CertInfo } from '@/api/ssl'
import { formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const requesting = ref(false)
const certs = ref<CertInfo[]>([])
const showModal = ref(false)
const form = reactive({ domain: '', provider: 'letsencrypt', auto_renew: true })

const columns = [
  { title: t('ssl.domain'), dataIndex: 'domain', key: 'domain' },
  { title: t('ssl.provider'), dataIndex: 'provider', key: 'provider', width: 120 },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('ssl.expiresAt'), dataIndex: 'expires_at', key: 'expires_at', width: 180 },
  { title: t('ssl.autoRenew'), key: 'auto_renew', width: 80 },
  { title: t('common.actions'), key: 'actions', width: 180 },
]

async function loadCerts() { loading.value = true; try { const res = await listCerts(); certs.value = res.data.list } catch {} finally { loading.value = false } }
async function handleRequest() { requesting.value = true; try { await requestCert(form); message.success(t('common.success')); showModal.value = false; loadCerts() } catch {} finally { requesting.value = false } }
async function handleRenew(cert: CertInfo) { try { await renewCert(cert.id); message.success(t('common.success')); loadCerts() } catch {} }
async function handleDelete(cert: CertInfo) { try { await deleteCert(cert.id); message.success(t('common.success')); loadCerts() } catch {} }

onMounted(loadCerts)
</script>
