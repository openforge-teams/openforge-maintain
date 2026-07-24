<template>
  <div>
    <PageHeader :title="t('menu.files')">
      <template #extra>
        <a-space>
          <a-button @click="showUploadModal = true"><UploadOutlined /> {{ t('files.upload') }}</a-button>
          <a-button @click="showNewFolderModal = true"><FolderAddOutlined /> {{ t('files.newFolder') }}</a-button>
          <a-button @click="handleDelete" :disabled="!selectedKeys.length" danger><DeleteOutlined /> {{ t('common.delete') }}</a-button>
          <a-button @click="loadFiles"><ReloadOutlined /> {{ t('common.refresh') }}</a-button>
        </a-space>
      </template>
    </PageHeader>

    <a-breadcrumb style="margin-bottom: 16px">
      <a-breadcrumb-item v-for="item in breadcrumbItems" :key="item.path">
        <a @click="navigateTo(item.path)">{{ item.name }}</a>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <a-row :gutter="16">
      <a-col :span="6">
        <a-card size="small" :title="t('files.path')" style="max-height: 600px; overflow-y: auto">
          <a-tree :tree-data="dirTree" :selected-keys="[currentPath]" @select="onTreeSelect"
            :field-names="{ title: 'name', key: 'path', children: 'children' }" />
        </a-card>
      </a-col>
      <a-col :span="18">
        <a-card>
          <a-table :dataSource="fileList" :columns="columns"
            :row-selection="{ selectedRowKeys: selectedKeys, onChange: onSelectChange }"
            rowKey="name" size="small" :pagination="false" :loading="loading"
            @rowDblclick="onRowDblClick">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'name'">
                <span v-if="record.is_dir"><FolderOutlined style="color: #1890ff" /> {{ record.name }}</span>
                <a v-else @click="openEditor(record)"><FileOutlined /> {{ record.name }}</a>
              </template>
              <template v-else-if="column.key === 'size'">{{ record.is_dir ? '-' : formatFileSize(record.size) }}</template>
              <template v-else-if="column.key === 'mod_time'">{{ formatTime(record.mod_time) }}</template>
              <template v-else-if="column.key === 'actions'">
                <a-space>
                  <a-button v-if="!record.is_dir" type="link" size="small" @click="openEditor(record)">{{ t('files.editFile') }}</a-button>
                  <a-popconfirm :title="t('common.delete') + '?'" @confirm="handleDeleteSingle(record)">
                    <a-button type="link" size="small" danger>{{ t('common.delete') }}</a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <a-modal v-model:open="showUploadModal" :title="t('files.upload')" :footer="null">
      <a-upload-dragger name="file" :multiple="true" :customRequest="handleUpload">
        <p class="ant-upload-drag-icon"><InboxOutlined /></p>
        <p class="ant-upload-text">{{ t('files.dragUpload') }}</p>
      </a-upload-dragger>
    </a-modal>

    <a-modal v-model:open="showNewFolderModal" :title="t('files.newFolder')" @ok="handleCreateFolder">
      <a-input v-model:value="newFolderName" />
    </a-modal>

    <Editor v-if="showEditor" :visible="showEditor" :file-path="editingFile?.path || ''" @close="showEditor = false" @saved="loadFiles" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { listFiles, uploadFile, deleteFile, mkdir } from '@/api/file'
import type { FileItem } from '@/api/file'
import { formatFileSize, formatTime } from '@/utils/format'
import PageHeader from '@/components/PageHeader.vue'
import Editor from './Editor.vue'
import { UploadOutlined, FolderAddOutlined, DeleteOutlined, ReloadOutlined, FolderOutlined, FileOutlined, InboxOutlined } from '@ant-design/icons-vue'

const { t } = useI18n()
const loading = ref(false)
const currentPath = ref('/root')
const fileList = ref<FileItem[]>([])
const selectedKeys = ref<string[]>([])
const showUploadModal = ref(false)
const showNewFolderModal = ref(false)
const showEditor = ref(false)
const newFolderName = ref('')
const editingFile = ref<FileItem | null>(null)

const columns = [
  { title: t('files.fileName'), key: 'name', dataIndex: 'name' },
  { title: t('common.size'), key: 'size', dataIndex: 'size', width: 120 },
  { title: t('files.permissions_col'), key: 'mode', dataIndex: 'mode', width: 100 },
  { title: t('files.modifyTime'), key: 'mod_time', dataIndex: 'mod_time', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 200 },
]

const breadcrumbItems = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  const items: { name: string; path: string }[] = [{ name: '/', path: '/' }]
  let path = ''
  for (const part of parts) { path += '/' + part; items.push({ name: part, path }) }
  return items
})

const dirTree = computed(() => fileList.value.filter(f => f.is_dir).map(f => ({ name: f.name, path: f.path, isLeaf: false, children: [] })))

async function loadFiles() {
  loading.value = true
  try { const res = await listFiles(currentPath.value); fileList.value = res.data } catch { /* ignore */ } finally { loading.value = false }
}

function navigateTo(path: string) { currentPath.value = path; loadFiles() }
function onTreeSelect(keys: string[]) { if (keys[0]) { currentPath.value = keys[0]; loadFiles() } }
function onSelectChange(keys: string[]) { selectedKeys.value = keys }
function onRowDblClick(record: FileItem) { if (record.is_dir) { currentPath.value = record.path; loadFiles() } else { openEditor(record) } }
function openEditor(file: FileItem) { editingFile.value = file; showEditor.value = true }

async function handleUpload(options: any) {
  try { await uploadFile(currentPath.value, options.file); message.success(t('common.success')); loadFiles() } catch { /* handled */ }
  options.onSuccess({}, options.file)
}

async function handleCreateFolder() {
  if (!newFolderName.value) return
  try { await mkdir(currentPath.value + '/' + newFolderName.value); message.success(t('common.success')); showNewFolderModal.value = false; newFolderName.value = ''; loadFiles() } catch { /* handled */ }
}

async function handleDeleteSingle(record: FileItem) {
  try { await deleteFile(record.path); message.success(t('common.success')); loadFiles() } catch { /* handled */ }
}

async function handleDelete() {
  try { for (const key of selectedKeys.value) { await deleteFile(currentPath.value + '/' + key) }; message.success(t('common.success')); selectedKeys.value = []; loadFiles() } catch { /* handled */ }
}

onMounted(loadFiles)
</script>
