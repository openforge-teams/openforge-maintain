<template>
  <div>
    <PageHeader :title="t('menu.host')" />
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="cpu" :tab="t('host.cpu')">
        <a-row :gutter="[16, 16]">
          <a-col :span="6"><a-card><a-statistic :title="t('host.cores')" :value="cpuInfo.cores" /></a-card></a-col>
          <a-col :span="6"><a-card><a-statistic :title="t('host.model')" :value="cpuInfo.model || '-'" /></a-card></a-col>
          <a-col :span="12"><a-card><a-statistic title="Usage" :value="cpuInfo.usage" suffix="%" :precision="1" /></a-card></a-col>
        </a-row>
        <a-card style="margin-top: 16px"><div ref="cpuChartRef" style="height: 350px"></div></a-card>
      </a-tab-pane>
      <a-tab-pane key="memory" :tab="t('host.memory')">
        <a-row :gutter="[16, 16]">
          <a-col :span="6"><a-card><a-statistic :title="t('host.total')" :value="memInfo.total" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
          <a-col :span="6"><a-card><a-statistic :title="t('host.used')" :value="memInfo.used" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
          <a-col :span="6"><a-card><a-statistic :title="t('host.cached')" :value="memInfo.cached" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
          <a-col :span="6"><a-card><a-statistic :title="t('host.free')" :value="memInfo.free" :formatter="(v: number) => formatFileSize(v)" /></a-card></a-col>
        </a-row>
        <a-card style="margin-top: 16px"><div ref="memChartRef" style="height: 350px"></div></a-card>
      </a-tab-pane>
      <a-tab-pane key="disk" :tab="t('host.disk')">
        <a-row :gutter="[16, 16]">
          <a-col :span="24" :lg="12"><a-card title="Usage"><div ref="diskPieRef" style="height: 300px"></div></a-card></a-col>
          <a-col :span="24" :lg="12"><a-card title="Partitions"><a-table :dataSource="diskInfo.partitions" :columns="diskColumns" size="small" :pagination="false" rowKey="mount" /></a-card></a-col>
        </a-row>
      </a-tab-pane>
      <a-tab-pane key="network" :tab="t('host.network')">
        <a-card><div ref="netChartRef" style="height: 350px"></div></a-card>
      </a-tab-pane>
      <a-tab-pane key="processes" :tab="t('host.processes')">
        <a-card>
          <template #extra><a-input-search v-model:value="searchProcess" :placeholder="t('common.search')" style="width: 200px" /></template>
          <a-table :dataSource="filteredProcesses" :columns="processColumns" size="small" :pagination="{ pageSize: 20 }" rowKey="pid" />
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import { getCPU, getMemory, getDisk, getNetwork, getProcesses } from '@/api/dashboard'
import type { CPUInfo, MemoryInfo, DiskInfo, NetworkInfo, ProcessInfo } from '@/api/dashboard'
import { formatFileSize } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'

const { t } = useI18n()
const activeTab = ref('cpu')
const searchProcess = ref('')

const cpuInfo = ref<CPUInfo>({ usage: 0, cores: 0, model: '', history: [] })
const memInfo = ref<MemoryInfo>({ total: 0, used: 0, free: 0, cached: 0, buffers: 0, history: [] })
const diskInfo = ref<DiskInfo>({ total: 0, used: 0, free: 0, partitions: [] })
const networkInfo = ref<NetworkInfo>({ interfaces: [], history: [] })
const processes = ref<ProcessInfo[]>([])

const cpuChartRef = ref<HTMLElement>()
const memChartRef = ref<HTMLElement>()
const diskPieRef = ref<HTMLElement>()
const netChartRef = ref<HTMLElement>()
let charts: echarts.ECharts[] = []
let timer: ReturnType<typeof setInterval> | null = null

const diskColumns = [
  { title: t('host.partition'), dataIndex: 'device', key: 'device' },
  { title: t('host.mount'), dataIndex: 'mount', key: 'mount' },
  { title: t('host.total'), dataIndex: 'total', key: 'total', customRender: ({ text }: any) => formatFileSize(text) },
  { title: t('host.used'), dataIndex: 'used', key: 'used', customRender: ({ text }: any) => formatFileSize(text) },
  { title: 'Usage', dataIndex: 'usage', key: 'usage', customRender: ({ text }: any) => text.toFixed(1) + '%' },
]

const processColumns = [
  { title: 'PID', dataIndex: 'pid', key: 'pid', sorter: (a: ProcessInfo, b: ProcessInfo) => a.pid - b.pid },
  { title: t('host.processName'), dataIndex: 'name', key: 'name' },
  { title: t('host.user'), dataIndex: 'user', key: 'user' },
  { title: 'CPU %', dataIndex: 'cpu', key: 'cpu', sorter: (a: ProcessInfo, b: ProcessInfo) => a.cpu - b.cpu },
  { title: 'MEM %', dataIndex: 'memory', key: 'memory', sorter: (a: ProcessInfo, b: ProcessInfo) => a.memory - b.memory },
  { title: t('common.status'), dataIndex: 'status', key: 'status' },
]

const filteredProcesses = computed(() => {
  if (!searchProcess.value) return processes.value
  const q = searchProcess.value.toLowerCase()
  return processes.value.filter(p => p.name.toLowerCase().includes(q) || String(p.pid).includes(q))
})

function initChart(el: HTMLElement) { const c = echarts.init(el); charts.push(c); return c }

async function fetchData() {
  try {
    const [cpu, mem, disk, net, procs] = await Promise.all([getCPU(), getMemory(), getDisk(), getNetwork(), getProcesses()])
    cpuInfo.value = cpu.data; memInfo.value = mem.data; diskInfo.value = disk.data; networkInfo.value = net.data; processes.value = procs.data
    updateCharts()
  } catch { /* ignore */ }
}

function updateCharts() {
  if (cpuChartRef.value && !charts[0]) { const c = initChart(cpuChartRef.value); c.setOption({ tooltip: { trigger: 'axis' }, xAxis: { type: 'category', data: [] }, yAxis: { type: 'value', max: 100 }, series: [{ type: 'line', smooth: true, data: [], areaStyle: {} }] }) }
  if (charts[0] && cpuInfo.value.history?.length) { charts[0].setOption({ xAxis: { data: cpuInfo.value.history.map(h => h.time) }, series: [{ data: cpuInfo.value.history.map(h => h.value) }] }) }
  if (memChartRef.value && !charts[1]) { const c = initChart(memChartRef.value); c.setOption({ tooltip: { trigger: 'axis' }, xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' }, series: [{ type: 'line', smooth: true, data: [], areaStyle: {} }] }) }
  if (charts[1] && memInfo.value.history?.length) { charts[1].setOption({ xAxis: { data: memInfo.value.history.map(h => h.time) }, series: [{ data: memInfo.value.history.map(h => h.value) }] }) }
  if (diskPieRef.value && !charts[2]) { const c = initChart(diskPieRef.value); c.setOption({ tooltip: { trigger: 'item' }, series: [{ type: 'pie', radius: ['40%', '70%'], data: [] }] }) }
  if (charts[2] && diskInfo.value.partitions?.length) { charts[2].setOption({ series: [{ data: diskInfo.value.partitions.map(p => ({ name: p.mount, value: p.used })) }] }) }
  if (netChartRef.value && !charts[3]) { const c = initChart(netChartRef.value); c.setOption({ tooltip: { trigger: 'axis' }, legend: { data: ['RX', 'TX'] }, xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' }, series: [{ name: 'RX', type: 'line', smooth: true, data: [] }, { name: 'TX', type: 'line', smooth: true, data: [] }] }) }
  if (charts[3] && networkInfo.value.history?.length) { charts[3].setOption({ xAxis: { data: networkInfo.value.history.map(h => h.time) }, series: [{ data: networkInfo.value.history.map(h => h.rx) }, { data: networkInfo.value.history.map(h => h.tx) }] }) }
}

onMounted(() => { fetchData(); timer = setInterval(fetchData, 5000) })
onUnmounted(() => { timer && clearInterval(timer); charts.forEach(c => c.dispose()) })
</script>