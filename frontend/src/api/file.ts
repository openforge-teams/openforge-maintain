import { get, post, del } from './request'

export interface FileItem {
  name: string
  path: string
  is_dir: boolean
  size: number
  mode: string
  mod_time: string
  owner: string
  group: string
}

export interface FileContent {
  path: string
  content: string
  encoding: string
  size: number
}

export interface FileInfo {
  name: string
  path: string
  size: number
  mode: string
  mod_time: string
  is_dir: boolean
  owner: string
  group: string
  mime_type: string
}

export function listFiles(path: string) {
  return get<FileItem[]>('/api/v2/files', { path })
}

export function getFileContent(path: string) {
  return get<FileContent>('/api/v2/files/content', { path })
}

export function saveFileContent(path: string, content: string) {
  return post('/api/v2/files/content', { path, content })
}

export function uploadFile(path: string, file: File, onProgress?: (percent: number) => void) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('path', path)
  return post('/api/v2/files/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress(e) {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
}

export function downloadFile(path: string) {
  return get<Blob>('/api/v2/files/download', { path }, { responseType: 'blob' })
}

export function deleteFile(path: string) {
  return del('/api/v2/files', { params: { path } })
}

export function renameFile(path: string, new_name: string) {
  return post('/api/v2/files/rename', { path, new_name })
}

export function mkdir(path: string) {
  return post('/api/v2/files/mkdir', { path })
}

export function getFileInfo(path: string) {
  return get<FileInfo>('/api/v2/files/info', { path })
}

export function compress(path: string, dest: string) {
  return post('/api/v2/files/compress', { path, dest })
}

export function extract(path: string, dest: string) {
  return post('/api/v2/files/extract', { path, dest })
}