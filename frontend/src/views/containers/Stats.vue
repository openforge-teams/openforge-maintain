<template>
  <div>
    <PageHeader :title="t('containers.stats')">
      <template #extra><a-button @click="$router.back()">{{ t('common.cancel') }}</a-button></template>
    </PageHeader>
    <a-row :gutter="[16, 16]">
      <a-col :span="24" :lg="12"><a-card :title="t('containers.cpuUsage')"><div ref="cpuChartRef" style="height: 300px"></div></a-card></a-col>
      <a-col :span="24" :lg="12"><a-card :title="t('containers.memUsage')"><div ref="memChartRef" style="height: 300px"></div></a-card></a-col>
      <a-col :span="24"><a-card :title="t('containers.netIO')"><div ref="netChartRef" style="height: 300px"></div></a-card></a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import { getContainerStats } from '@/api/container'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const route = useRoute()
const { t } = useI18n()
const containerId = route.query.id as string
const cpuChartRef = ref<HTMLElement>()
const memChartRef = ref<HTMLElement>()
const netChartRef = ref<HTMLElement>()
let charts: echarts.ECharts[] = []
let timer: ReturnType<typeof setInterval> | null = null

function mk(el: HTMLElement, series: any[]) {
  const c = echarts.init(el)
  c.setOption({ tooltip: { trigger: 'axis' }, xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' }, series })
  charts.push(c)
  return c
}

async function fetchData() {
  if (!containerId) return
  try {
    const res = await getContainerStats(containerId)
    const h = res.data.history || []
    const times = h.map(x => x.time)
    if (charts[0]) charts[0].setOption({ xAxis: { data: times }, series: [{ data: h.map(x => x.cpu) }] })
    if (charts[1]) charts[1].setOption({ xAxis: { data: times }, series: [{ data: h.map(x => x.memory) }] })
    if (charts[2]) charts[2].setOption({ xAxis: { data: times }, series: [{ data: h.map(x => x.rx) }, { data: h.map(x => x.tx) }] })
  } catch { /* ignore */ }
}

onMounted(() => {
  if (cpuChartRef.value) mk(cpuChartRef.value, [{ name: 'CPU', type: 'line', smooth: true, data: [] }])
  if (memChartRef.value) mk(memChartRef.value, [{ name: 'MEM', type: 'line', smooth: true, data: [] }])
  if (netChartRef.value) mk(netChartRef.value, [{ name: 'RX', type: 'line', smooth: true, data: [] }, { name: 'TX', type: 'line', smooth: true, data: [] }])
  fetchData()
  timer = setInterval(fetchData, 3000)
})

onUnmounted(() => { timer && clearInterval(timer); charts.forEach(c => c.dispose()) })
</script>