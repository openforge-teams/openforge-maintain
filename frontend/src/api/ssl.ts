import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface CertInfo {
  id: number
  domain: string
  provider: string
  issued_at: string
  expires_at: string
  auto_renew: boolean
  status: 'active' | 'expired' | 'pending'
}

export function listCerts(params?: PageParams) {
  return get<PageResult<CertInfo>>('/api/v2/ssl', params)
}

export function requestCert(data: { domain: string; provider: string; auto_renew?: boolean }) {
  return post<CertInfo>('/api/v2/ssl/request', data)
}

export function renewCert(id: number) {
  return post(`/api/v2/ssl/${id}/renew`)
}

export function deleteCert(id: number) {
  return del(`/api/v2/ssl/${id}`)
}
