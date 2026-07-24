<template>
  <div>
    <PageHeader :title="t('databases.postgres')" />
    <a-card>
      <a-table :dataSource="databases" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'size'">{{ formatFileSize(record.size) }}</template>
          <template v-else-if="column.key === 'actions'">
            <a-popconfirm :title="t('databases.deleteDB') + '?'" @confirm="handleDelete(record)">
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
import { listPostgres } from '@/api/database'
import type { PostgresDB } from '@/api/database'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const { t } = useI18n()
const loading = ref(false)
const databases = ref<PostgresDB[]>([])

const columns = [
  { title: t('databases.dbName'), dataIndex: 'name', key: 'name' },
  { title: 'Owner', dataIndex: 'owner', key: 'owner', width: 120 },
  { title: t('common.size'), key: 'size', width: 120 },
  { title: t('common.created_at'), dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadDBs() { loading.value = true; try { const res = await listPostgres(); databases.value = res.data.list } catch {} finally { loading.value = false } }
async function handleDelete(db: PostgresDB) { try { await import('@/api/database').then(m => m.deleteMySQLDB(db.id)); message.success(t('common.success')); loadDBs() } catch {} }

onMounted(loadDBs)
</script>
