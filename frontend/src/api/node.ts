import { get, post, put, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface NodeInfo {
  id: number
  name: string
  host: string
  port: number
  status: string
  os: string
  arch: string
  created_at: string
}

export function listNodes(params?: PageParams) {
  return get<PageResult<NodeInfo>>('/api/v2/nodes', params)
}

export function createNode(data: Partial<NodeInfo>) {
  return post<NodeInfo>('/api/v2/nodes', data)
}

export function getNode(id: number) {
  return get<NodeInfo>(`/api/v2/nodes/${id}`)
}

export function updateNode(id: number, data: Partial<NodeInfo>) {
  return put<NodeInfo>(`/api/v2/nodes/${id}`, data)
}

export function deleteNode(id: number) {
  return del(`/api/v2/nodes/${id}`)
}