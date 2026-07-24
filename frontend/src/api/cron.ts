import { get, post, put, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface CronJob {
  id: number
  name: string
  spec: string
  command: string
  status: 'enabled' | 'disabled'
  last_run_at: string
  next_run_at: string
  created_at: string
}

export function listCronJobs(params?: PageParams) {
  return get<PageResult<CronJob>>('/api/v2/cron', params)
}

export function createCronJob(data: { name: string; spec: string; command: string }) {
  return post<CronJob>('/api/v2/cron', data)
}

export function updateCronJob(id: number, data: Partial<CronJob>) {
  return put<CronJob>(`/api/v2/cron/${id}`, data)
}

export function deleteCronJob(id: number) {
  return del(`/api/v2/cron/${id}`)
}

export function startCronJob(id: number) {
  return post(`/api/v2/cron/${id}/start`)
}

export function stopCronJob(id: number) {
  return post(`/api/v2/cron/${id}/stop`)
}

export function runCronJob(id: number) {
  return post(`/api/v2/cron/${id}/run`)
}