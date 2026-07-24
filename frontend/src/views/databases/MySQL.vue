<template>
  <div>
    <PageHeader :title="t('databases.mysql')">
      <template #extra>
        <a-button type="primary" @click="showCreateModal = true"><PlusOutlined /> {{ t('databases.createDB') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="databases" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'size'">{{ formatFileSize(record.size) }}</template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="handleBackup(record)">{{ t('databases.backup') }}</a-button>
              <a-popconfirm :title="t('databases.deleteDB') + '?'" @confirm="handleDelete(record)">
                <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showCreateModal" :title="t('databases.createDB')" @ok="handleCreate" :confirmLoading="creating">
      <a-form :model="createForm" layout="vertical">
        <a-form-item :label="t('databases.dbName')"><a-input v-model:value="createForm.name" /></a-form-item>
        <a-form-item :label="t('databases.characterSet')"><a-input v-model:value="createForm.character_set" placeholder="utf8mb4" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listMySQL, createMySQLDB, deleteMySQLDB, backupMySQL } from '@/api/database'
import type { MySQLDB } from '@/api/database'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const creating = ref(false)
const databases = ref<MySQLDB[]>([])
const showCreateModal = ref(false)
const createForm = reactive({ name: '', character_set: 'utf8mb4' })

const columns = [
  { title: t('databases.dbName'), dataIndex: 'name', key: 'name' },
  { title: t('databases.characterSet'), dataIndex: 'character_set', key: 'character_set', width: 120 },
  { title: t('common.size'), key: 'size', width: 120 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 200 },
]

async function loadDBs() { loading.value = true; try { const res = await listMySQL(); databases.value = res.data.list } catch {} finally { loading.value = false } }
async function handleCreate() { creating.value = true; try { await createMySQLDB(createForm); message.success(t('common.success')); showCreateModal.value = false; loadDBs() } catch {} finally { creating.value = false } }
async function handleDelete(db: MySQLDB) { try { await deleteMySQLDB(db.id); message.success(t('common.success')); loadDBs() } catch {} }
async function handleBackup(db: MySQLDB) { try { await backupMySQL(db.id); message.success(t('common.success')) } catch {} }

onMounted(loadDBs)
</script>
