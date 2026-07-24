<template>
  <div>
    <PageHeader :title="t('menu.images')">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="showPullModal = true"><CloudDownloadOutlined /> {{ t('images.pullImage') }}</a-button>
          <a-button @click="loadImages"><ReloadOutlined /> {{ t('common.refresh') }}</a-button>
        </a-space>
      </template>
    </PageHeader>
    <a-card>
      <a-table :dataSource="images" :columns="columns" rowKey="id" :loading="loading" :pagination="{ pageSize: 20 }">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'repo_tags'">
            <span v-for="tag in record.repo_tags" :key="tag">{{ tag }}<br /></span>
          </template>
          <template v-else-if="column.key === 'size'">{{ formatFileSize(record.size) }}</template>
          <template v-else-if="column.key === 'created_at'">{{ formatTime(record.created_at) }}</template>
          <template v-else-if="column.key === 'actions'">
            <a-popconfirm title="Delete image?" @confirm="handleRemove(record)">
              <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="showPullModal" :title="t('images.pullImage')" @ok="handlePull" :confirmLoading="pulling">
      <a-form layout="vertical">
        <a-form-item :label="t('images.imageName')"><a-input v-model:value="pullImage" placeholder="nginx" /></a-form-item>
        <a-form-item :label="t('images.imageTag')"><a-input v-model:value="pullTag" placeholder="latest" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listImages, pullImage as pullImageApi, removeImage } from '@/api/container'
import type { ImageInfo } from '@/api/container'
import { formatFileSize, formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import { CloudDownloadOutlined, ReloadOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const images = ref<ImageInfo[]>([])
const showPullModal = ref(false)
const pulling = ref(false)
const pullImage = ref('')
const pullTag = ref('latest')

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 120, ellipsis: true },
  { title: t('images.repo'), key: 'repo_tags' },
  { title: t('common.size'), key: 'size', width: 120 },
  { title: t('common.created_at'), key: 'created_at', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 100 },
]

async function loadImages() {
  loading.value = true
  try { const res = await listImages(); images.value = res.data.list } catch { /* ignore */ } finally { loading.value = false }
}

async function handlePull() {
  pulling.value = true
  try { await pullImageApi(pullImage.value, pullTag.value || 'latest'); message.success(t('images.pullSuccess')); showPullModal.value = false; loadImages() } catch { /* handled */ } finally { pulling.value = false }
}

async function handleRemove(img: ImageInfo) {
  try { await removeImage(img.id); message.success(t('common.success')); loadImages() } catch { /* handled */ }
}

onMounted(loadImages)
</script>
