import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface BackupTask {
  id: number
  name: string
  type: string
  target: string
  schedule: string
  retention: number
  enabled: boolean
  last_run_at: string
  created_at: string
}

export interface BackupRecord {
  id: number
  task_id: number
  task_name: string
  file_path: string
  file_size: number
  status: 'success' | 'failed' | 'running'
  created_at: string
}

export function listBackupTasks(params?: PageParams) {
  return get<PageResult<BackupTask>>('/api/v2/backup/tasks', params)
}

export function createBackupTask(data: { name: string; type: string; target: string; schedule: string; retention: number }) {
  return post<BackupTask>('/api/v2/backup/tasks', data)
}

export function deleteBackupTask(id: number) {
  return del(`/api/v2/backup/tasks/${id}`)
}

export function listBackups(params?: PageParams & { task_id?: number }) {
  return get<PageResult<BackupRecord>>('/api/v2/backup/records', params)
}

export function restoreBackup(id: number) {
  return post(`/api/v2/backup/records/${id}/restore`)
}
