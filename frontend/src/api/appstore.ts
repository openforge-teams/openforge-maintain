import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface AppInfo {
  id: number
  name: string
  icon: string
  description: string
  category: string
  version: string
  author: string
  install_params: AppParam[]
}

export interface AppParam {
  name: string
  label: string
  type: 'port' | 'password' | 'path' | 'select' | 'text'
  default_value: string
  required: boolean
  options?: string[]
}

export interface InstalledApp {
  id: number
  app_id: number
  name: string
  status: 'running' | 'stopped' | 'installing' | 'error'
  version: string
  ports: { app_port: number; host_port: number }[]
  created_at: string
}

export function getAppList(params?: PageParams & { category?: string }) {
  return get<PageResult<AppInfo>>('/api/v2/appstore', params)
}

export function getAppDetail(id: number) {
  return get<AppInfo>(`/api/v2/appstore/${id}`)
}

export function installApp(id: number, params: Record<string, string>) {
  return post<InstalledApp>(`/api/v2/appstore/${id}/install`, params)
}

export function uninstallApp(id: number) {
  return del(`/api/v2/appstore/apps/${id}`)
}

export function upgradeApp(id: number) {
  return post(`/api/v2/appstore/apps/${id}/upgrade`)
}

export function getInstalledApps(params?: PageParams) {
  return get<PageResult<InstalledApp>>('/api/v2/appstore/installed', params)
}