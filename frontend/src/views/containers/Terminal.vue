<template>
  <div>
    <PageHeader :title="t('containers.terminal')">
      <template #extra>
        <a-button @click="$router.back()">{{ t('common.cancel') }}</a-button>
      </template>
    </PageHeader>
    <a-card :bodyStyle="{ padding: 0 }">
      <div ref="terminalRef" class="terminal-container" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'
import PageHeader from '@/components/PageHeader.vue'

const route = useRoute()
const { t } = useI18n()
const terminalRef = ref<HTMLElement>()
const containerId = route.query.id as string
let term: Terminal | null = null
let ws: WebSocket | null = null
let fitAddon: FitAddon | null = null
let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (!terminalRef.value) return
  term = new Terminal({ cursorBlink: true, fontSize: 14, fontFamily: 'Consolas, Monaco, monospace', theme: { background: '#1e1e1e' } })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())
  term.open(terminalRef.value)
  fitAddon.fit()

  const token = localStorage.getItem('token')
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${protocol}//${location.host}/ws/terminal?container_id=${containerId}&token=${token}`)
  ws.onmessage = (e) => term?.write(e.data)
  ws.onclose = () => term?.writeln('\r\n\x1b[31mDisconnected\x1b[0m')
  ws.onerror = () => term?.writeln('\r\n\x1b[31mConnection error\x1b[0m')
  term.onData((data) => ws?.send(data))

  resizeObserver = new ResizeObserver(() => fitAddon?.fit())
  resizeObserver.observe(terminalRef.value)
})

onUnmounted(() => {
  ws?.close()
  term?.dispose()
  resizeObserver?.disconnect()
})
</script>

<style scoped>
.terminal-container { height: 70vh; padding: 8px; background: #1e1e1e; }
</style>
