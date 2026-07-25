<template>
  <div>
    <PageHeader :title="t('dashboard.title')" />
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card><a-statistic :title="t('dashboard.cpuUsage')" :value="overview.cpu_usage" suffix="%" :precision="1">
          <template #prefix><DashboardOutlined /></template></a-statistic>
          <a-progress :percent="overview.cpu_usage" size="small" style="margin-top: 12px" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card><a-statistic :title="t('dashboard.memUsage')" :value="overview.memory_usage" suffix="%" :precision="1">
          <template #prefix><CloudOutlined /></template></a-statistic>
          <a-progress :percent="overview.memory_usage" size="small" style="margin-top: 12px" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card><a-statistic :title="t('dashboard.diskUsage')" :value="overview.disk_usage" suffix="%" :precision="1">
          <template #prefix><DatabaseOutlined /></template></a-statistic>
          <a-progress :percent="overview.disk_usage" size="small" style="margin-top: 12px" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card><a-statistic :title="t('dashboard.containerStats')" :value="overview.container_running">
          <template #prefix><ContainerOutlined /></template>
          <template #suffix>/ {{ overview.container_count }}</template></a-statistic>
        </a-card>
      </a-col>
    </a-row>
    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <a-col :xs="24" :lg="12"><a-card :title="t('dashboard.cpuRealtime')"><div ref="cpuChartRef" style="height: 300px"></div></a-card></a-col>
      <a-col :xs="24" :lg="12"><a-card :title="t('dashboard.memRealtime')"><div ref="memChartRef" style="height: 300px"></div></a-card></a-col>
    </a-row>
    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <a-col :xs="24" :lg="12">
        <a-card :title="t('dashboard.quickActions')">
          <a-row :gutter="[16, 16]">
            <a-col :span="8"><a-button type="primary" block @click="$router.push('/terminal')"><CodeOutlined /> Terminal</a-button></a-col>
            <a-col :span="8"><a-button block @click="$router.push('/files')"><FolderOutlined /> {{ t('menu.files') }}</a-button></a-col>
            <a-col :span="8"><a-button block @click="$router.push('/containers')"><ContainerOutlined /> {{ t('menu.containers') }}</a-button></a-col>
            <a-col :span="8"><a-button block @click="$router.push('/websites')"><GlobalOutlined /> {{ t('menu.websites') }}</a-button></a-col>
            <a-col :span="8"><a-button block @click="$router.push('/appstore')"><AppstoreOutlined /> {{ t('menu.appstore') }}</a-button></a-col>
            <a-col :span="8"><a-button block @click="$router.push('/backup')"><CloudUploadOutlined /> {{ t('menu.backup') }}</a-button></a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card :title="t('dashboard.recentLogs')">
          <a-table :dataSource="recentLogs" :columns="logColumns" size="small" :pagination="false" rowKey="time" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import { getOverview, getCPU, getMemory } from '@/api/dashboard'
import PageHeader from '@/components/PageHeader.vue'
import { DashboardOutlined, CloudOutlined, DatabaseOutlined, ContainerOutlined, CodeOutlined, FolderOutlined, GlobalOutlined, AppstoreOutlined, CloudUploadOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const cpuChartRef = ref<HTMLElement>()
const memChartRef = ref<HTMLElement>()
let cpuChart: echarts.ECharts | null = null
let memChart: echarts.ECharts | null = null
let timer: ReturnType<typeof setInterval> | null = null

const overview = ref({ cpu_usage: 0, memory_usage: 0, memory_total: 0, memory_used: 0, disk_usage: 0, disk_total: 0, disk_used: 0, network_in: 0, network_out: 0, container_count: 0, container_running: 0, container_stopped: 0, uptime: 0, hostname: '', os: '' })
const recentLogs = ref<any[]>([])
const logColumns = [{ title: 'Time', dataIndex: 'time', key: 'time' }, { title: 'Action', dataIndex: 'action', key: 'action' }, { title: 'User', dataIndex: 'user', key: 'user' }]

function createLineChart(el: HTMLElement, title: string, color: string) {
  const chart = echarts.init(el)
  chart.setOption({
    tooltip: { trigger: 'axis' }, legend: { data: [title] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: [], boundaryGap: false },
    yAxis: { type: 'value', max: 100, axisLabel: { formatter: '{value}%' } },
    series: [{ name: title, type: 'line', smooth: true, data: [], areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: color + '80' }, { offset: 1, color: color + '10' }]) }, lineStyle: { color }, itemStyle: { color } }],
  })
  return chart
}

async function fetchData() {
  try {
    const [ov, cpu, mem] = await Promise.all([getOverview(), getCPU(), getMemory()])
    overview.value = ov.data
    if (cpuChart && cpu.data.history?.length) {
      cpuChart.setOption({ xAxis: { data: cpu.data.history.map(h => h.time) }, series: [{ data: cpu.data.history.map(h => h.value) }] })
    }
    if (memChart && mem.data.history?.length) {
      memChart.setOption({ xAxis: { data: mem.data.history.map(h => h.time) }, series: [{ data: mem.data.history.map(h => h.value) }] })
    }
  } catch { /* ignore */ }
}

onMounted(async () => {
  if (cpuChartRef.value) cpuChart = createLineChart(cpuChartRef.value, t('dashboard.cpuUsage'), '#1890ff')
  if (memChartRef.value) memChart = createLineChart(memChartRef.value, t('dashboard.memUsage'), '#52c41a')
  await fetchData()
  timer = setInterval(fetchData, 5000)
})

onUnmounted(() => {
  timer && clearInterval(timer)
  cpuChart?.dispose()
  memChart?.dispose()
})
</script>
