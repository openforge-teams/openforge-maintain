<template>
  <a-badge :status="statusType" :text="text" />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  status?: string
  text?: string
}>(), {
  status: 'stopped',
  text: '',
})

const statusType = computed(() => {
  const map: Record<string, string> = {
    running: 'success',
    stopped: 'error',
    created: 'processing',
    paused: 'warning',
    restarting: 'processing',
    active: 'success',
    expired: 'error',
    pending: 'processing',
    enabled: 'success',
    disabled: 'default',
  }
  return (map[props.status] || 'default') as any
})
</script>
