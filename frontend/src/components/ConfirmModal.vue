<template>
  <a-modal
    :open="open"
    :title="title"
    :confirm-loading="loading"
    @ok="handleOk"
    @cancel="$emit('update:open', false)"
  >
    <p>{{ content }}</p>
  </a-modal>
</template>

<script setup lang="ts">
const props = defineProps<{
  open: boolean
  title?: string
  content?: string
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:open': [val: boolean]
  'ok': []
}>()

function handleOk() {
  emit('ok')
}
</script>
