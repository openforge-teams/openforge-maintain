<template>
  <div>
    <PageHeader :title="app.name || t('appstore.appDetail')">
      <template #extra><a-button @click="$router.back()">Back</a-button></template>
    </PageHeader>
    <a-row :gutter="[16, 16]">
      <a-col :span="16">
        <a-card :title="t('appstore.installConfig')">
          <a-form :model="params" layout="vertical">
            <a-form-item v-for="p in app.install_params" :key="p.name" :label="p.label">
              <a-input v-if="p.type !== 'select'" v-model:value="params[p.name]" :placeholder="p.default_value" />
              <a-select v-else v-model:value="params[p.name]">
                <a-select-option v-for="opt in (p.options || [])" :key="opt" :value="opt">{{ opt }}</a-select-option>
              </a-select>
            </a-form-item>
            <a-button type="primary" :loading="installing" @click="handleInstall">{{ t('appstore.install') }}</a-button>
          </a-form>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <p><strong>{{ t('common.description') }}:</strong> {{ app.description }}</p>
          <p><strong>Version:</strong> {{ app.version }}</p>
          <p><strong>Author:</strong> {{ app.author }}</p>
          <p><strong>Category:</strong> {{ app.category }}</p>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { getAppDetail, installApp } from '@/api/appstore'
import type { AppInfo } from '@/api/appstore'
import PageHeader from '@/components/PageHeader.vue'

const route = useRoute()
const { t } = useI18n()
const app = ref<AppInfo>({ id: 0, name: '', icon: '', description: '', category: '', version: '', author: '', install_params: [] })
const params = reactive<Record<string, string>>({})
const installing = ref(false)

async function loadApp() {
  try {
    const res = await getAppDetail(Number(route.params.id))
    app.value = res.data
    for (const p of res.data.install_params) { params[p.name] = p.default_value }
  } catch {}
}

async function handleInstall() {
  installing.value = true
  try { await installApp(app.value.id, params); message.success(t('appstore.installSuccess')) } catch {} finally { installing.value = false }
}

onMounted(loadApp)
</script>
