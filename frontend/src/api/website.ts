import { get, post, put, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface Website {
  id: number
  domain: string
  type: 'static' | 'reverse_proxy' | 'php'
  status: 'running' | 'stopped' | 'error'
  ssl_enabled: boolean
  ssl_expire: string
  proxy_target: string
  root_path: string
  created_at: string
}

export interface CreateWebsiteParams {
  domain: string
  type: string
  proxy_target?: string
  root_path?: string
  php_version?: string
}

export function listWebsites(params?: PageParams) {
  return get<PageResult<Website>>('/api/v2/websites', params)
}

export function createWebsite(data: CreateWebsiteParams) {
  return post<Website>('/api/v2/websites', data)
}

export function getWebsite(id: number) {
  return get<Website>(`/api/v2/websites/${id}`)
}

export function updateWebsite(id: number, data: Partial<CreateWebsiteParams>) {
  return put<Website>(`/api/v2/websites/${id}`, data)
}

export function deleteWebsite(id: number) {
  return del(`/api/v2/websites/${id}`)
}

export function enableSSL(id: number, data?: { domain: string; auto_renew?: boolean }) {
  return post(`/api/v2/websites/${id}/ssl/enable`, data)
}

export function disableSSL(id: number) {
  return post(`/api/v2/websites/${id}/ssl/disable`)
}