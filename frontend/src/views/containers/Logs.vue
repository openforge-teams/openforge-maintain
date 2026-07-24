<template>
  <div>
    <PageHeader :title="t('containers.logs')">
      <template #extra>
        <a-button @click="$router.back()">{{ t('common.cancel') }}</a-button>
      </template>
    </PageHeader>
    <a-card>
      <div class="log-container" ref="logContainerRef">
        <pre v-for="(line, idx) in logs" :key="idx">{{ line }}</pre>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getContainerLogs } from '@/api/container'
import PageHeader from '@/components/PageHeader.vue'

const route = useRoute()
const { t } = useI18n()
const logs = ref<string[]>([])
const logContainerRef = ref<HTMLElement>()
let timer: ReturnType<typeof setInterval> | null = null
const containerId = route.query.id as string

async function fetchLogs() {
  if (!containerId) return
  try {
    const res = await getContainerLogs(containerId, { tail: 500 })
    logs.value = res.data.split('\n')
    if (logContainerRef.value) logContainerRef.value.scrollTop = logContainerRef.value.scrollHeight
  } catch { /* handled */ }
}

onMounted(() => { fetchLogs(); timer = setInterval(fetchLogs, 3000) })
onUnmounted(() => { timer && clearInterval(timer) })
</script>

<style scoped>
.log-container {
  height: 70vh; overflow-y: auto; background: #1e1e1e; color: #d4d4d4;
  padding: 16px; border-radius: 4px; font-family: 'Consolas', 'Monaco', monospace; font-size: 13px;
}
pre { margin: 0; white-space: pre-wrap; word-break: break-all; }
</style>
