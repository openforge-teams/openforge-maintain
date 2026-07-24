<template>
  <div>
    <PageHeader :title="t('menu.firewall')">
      <template #extra>
        <a-button type="primary" @click="showAddModal = true"><PlusOutlined /> {{ t('firewall.addRule') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <template #title>
        <a-space>
          <span>{{ t('firewall.status') }}:</span>
          <a-switch :checked="fwEnabled" @change="handleToggleFW" :loading="fwLoading" />
          <a-tag :color="fwEnabled ? 'green' : 'red'">{{ fwEnabled ? t('firewall.enabled') : t('firewall.disabled') }}</a-tag>
        </a-space>
      </template>
      <a-table :dataSource="rules" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'enabled'">
            <a-switch :checked="record.enabled" @change="handleToggleRule(record)" />
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDelete(record)">
              <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showAddModal" :title="t('firewall.addRule')" @ok="handleAdd" :confirmLoading="adding">
      <a-form :model="ruleForm" layout="vertical">
        <a-form-item :label="t('firewall.protocol')">
          <a-select v-model:value="ruleForm.protocol">
            <a-select-option value="tcp">TCP</a-select-option>
            <a-select-option value="udp">UDP</a-select-option>
            <a-select-option value="icmp">ICMP</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('firewall.port')"><a-input v-model:value="ruleForm.port" placeholder="80, 443" /></a-form-item>
        <a-form-item :label="t('firewall.source')"><a-input v-model:value="ruleForm.source" placeholder="0.0.0.0/0" /></a-form-item>
        <a-form-item :label="t('firewall.ruleAction')">
          <a-select v-model:value="ruleForm.action">
            <a-select-option value="allow">{{ t('firewall.allow') }}</a-select-option>
            <a-select-option value="deny">{{ t('firewall.deny') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('firewall.comment')"><a-input v-model:value="ruleForm.comment" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listRules, addRule, deleteRule, enableRule, disableRule, getFirewallStatus, enableFirewall, disableFirewall } from '@/api/firewall'
import type { FirewallRule } from '@/api/firewall'
import PageHeader from '@/components/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const fwLoading = ref(false)
const adding = ref(false)
const fwEnabled = ref(false)
const rules = ref<FirewallRule[]>([])
const showAddModal = ref(false)
const ruleForm = reactive({ protocol: 'tcp', port: '', source: '0.0.0.0/0', action: 'allow', comment: '' })

const columns = [
  { title: t('firewall.protocol'), dataIndex: 'protocol', key: 'protocol', width: 80 },
  { title: t('firewall.port'), dataIndex: 'port', key: 'port', width: 120 },
  { title: t('firewall.source'), dataIndex: 'source', key: 'source' },
  { title: t('firewall.ruleAction'), dataIndex: 'action', key: 'action', width: 80 },
  { title: t('common.status'), key: 'enabled', width: 80 },
  { title: t('firewall.comment'), dataIndex: 'comment', key: 'comment' },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadData() {
  loading.value = true
  try {
    const [rulesRes, statusRes] = await Promise.all([listRules(), getFirewallStatus()])
    rules.value = rulesRes.data.list
    fwEnabled.value = statusRes.data.enabled
  } catch {} finally { loading.value = false }
}

async function handleToggleFW() {
  fwLoading.value = true
  try { if (fwEnabled.value) { await disableFirewall() } else { await enableFirewall() }; fwEnabled.value = !fwEnabled.value; message.success(t('common.success')) } catch {} finally { fwLoading.value = false }
}

async function handleToggleRule(rule: FirewallRule) {
  try { if (rule.enabled) { await disableRule(rule.id) } else { await enableRule(rule.id) }; rule.enabled = !rule.enabled; message.success(t('common.success')) } catch {}
}

async function handleAdd() {
  adding.value = true
  try { await addRule(ruleForm); message.success(t('common.success')); showAddModal.value = false; loadData() } catch {} finally { adding.value = false }
}

async function handleDelete(rule: FirewallRule) {
  try { await deleteRule(rule.id); message.success(t('common.success')); loadData() } catch {}
}

onMounted(loadData)
</script>