<template>
  <a-modal :open="visible" :title="filePath" width="90vw" :footer="null" @cancel="$emit('close')" :destroyOnClose="true">
    <div style="display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 12px">
      <a-button type="primary" @click="handleSave" :loading="saving">{{ t('common.save') }}</a-button>
      <a-button @click="$emit('close')">{{ t('common.cancel') }}</a-button>
    </div>
    <vue-monaco-editor v-model:value="content" language="plaintext" :height="'70vh'" theme="vs-dark"
      :options="{ minimap: { enabled: false }, fontSize: 14, wordWrap: 'on' }" />
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { getFileContent, saveFileContent } from '@/api/file'

const { t } = useI18n()
const props = defineProps<{ visible: boolean; filePath: string }>()
const emit = defineEmits<{ close: []; saved: [] }>()

const content = ref('')
const saving = ref(false)

watch(() => props.visible, async (val) => {
  if (val && props.filePath) {
    try { const res = await getFileContent(props.filePath); content.value = res.data.content } catch { /* handled */ }
  }
})

async function handleSave() {
  saving.value = true
  try { await saveFileContent(props.filePath, content.value); message.success(t('files.saveSuccess')); emit('saved'); emit('close') } catch { /* handled */ } finally { saving.value = false }
}
</script>
