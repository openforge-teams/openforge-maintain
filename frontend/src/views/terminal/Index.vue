<template>
  <div style="height: calc(100vh - 140px)">
    <div class="terminal-bar">
      <span>{{ connected ? t('terminal.connected') : t('terminal.disconnected') }}</span>
      <a-button v-if="!connected" size="small" @click="connect">{{ t('terminal.reconnect') }}</a-button>
    </div>
    <div ref="termRef" class="terminal" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'

const { t } = useI18n()
const termRef = ref<HTMLElement>()
const connected = ref(false)
let term: Terminal | null = null
let ws: WebSocket | null = null
let fitAddon: FitAddon | null = null
let resizeObserver: ResizeObserver | null = null

function connect() {
  if (!termRef.value || ws) return
  term = new Terminal({ cursorBlink: true, fontSize: 14, fontFamily: 'Consolas, Monaco, monospace', theme: { background: '#1e1e1e' } })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())
  term.open(termRef.value)
  fitAddon.fit()

  const token = localStorage.getItem('token')
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${protocol}//${location.host}/ws/terminal?token=${token}`)
  ws.onopen = () => { connected.value = true }
  ws.onmessage = (e) => term?.write(e.data)
  ws.onclose = () => { connected.value = false; ws = null; term?.writeln('\r\n\x1b[31mDisconnected\x1b[0m') }
  ws.onerror = () => { connected.value = false; ws = null }
  term.onData((data) => ws?.send(data))

  resizeObserver = new ResizeObserver(() => fitAddon?.fit())
  resizeObserver.observe(termRef.value)
}

onMounted(connect)

onUnmounted(() => {
  ws?.close()
  term?.dispose()
  resizeObserver?.disconnect()
})
</script>

<style scoped>
.terminal-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 8px 16px; background: #001529; color: #fff; border-radius: 4px 4px 0 0;
}
.terminal { height: calc(100% - 40px); padding: 8px; background: #1e1e1e; border-radius: 0 0 4px 4px; }
</style>
