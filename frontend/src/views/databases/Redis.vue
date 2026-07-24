<template>
  <div>
    <PageHeader :title="t('databases.redis')" />
    <a-row :gutter="[16, 16]">
      <a-col :span="6"><a-card><a-statistic :title="t('databases.version')" :value="redisInfo.version || '-'" /></a-card></a-col>
      <a-col :span="6"><a-card><a-statistic :title="t('databases.usedMemory')" :value="redisInfo.used_memory" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
      <a-col :span="6"><a-card><a-statistic :title="t('databases.maxMemory')" :value="redisInfo.max_memory" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
      <a-col :span="6"><a-card><a-statistic :title="t('databases.clients')" :value="redisInfo.connected_clients" /></a-card></a-col>
    </a-row>
    <a-card title="Databases" style="margin-top: 16px">
      <a-table :dataSource="redisInfo.databases || []" rowKey="index" :pagination="false" size="small">
        <a-table-column title="Index" dataIndex="index" key="index" />
        <a-table-column :title="t('databases.dbKeys')" dataIndex="keys" key="keys" />
        <a-table-column title="Expires" dataIndex="expires" key="expires" />
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getRedisInfo } from '@/api/database'
import type { RedisInfo } from '@/api/database'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const { t } = useI18n()
const redisInfo = ref<RedisInfo>({ version: '', used_memory: 0, max_memory: 0, connected_clients: 0, uptime: 0, databases: [] })

async function loadInfo() { try { const res = await getRedisInfo(); redisInfo.value = res.data } catch {} }

onMounted(loadInfo)
</script>
