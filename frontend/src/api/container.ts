import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface ContainerInfo {
  id: string
  name: string
  image: string
  status: string
  state: string
  ports: { ip: string; private_port: number; public_port: number; type: string }[]
  created_at: string
  labels: Record<string, string>
}

export interface ContainerStats {
  cpu_usage: number
  memory_usage: number
  memory_limit: number
  network_rx: number
  network_tx: number
  block_read: number
  block_write: number
  history: { time: string; cpu: number; memory: number; rx: number; tx: number }[]
}

export interface ImageInfo {
  id: string
  repo_tags: string[]
  size: number
  created_at: string
}

export interface VolumeInfo {
  name: string
  driver: string
  mountpoint: string
  created_at: string
  size: number
}

export interface NetworkInfo {
  id: string
  name: string
  driver: string
  scope: string
  created_at: string
  containers: number
}

export function listContainers(params?: PageParams & { status?: string }) {
  return get<PageResult<ContainerInfo>>('/api/v2/containers', params)
}

export function getContainer(id: string) {
  return get<ContainerInfo>(`/api/v2/containers/${id}`)
}

export function startContainer(id: string) {
  return post(`/api/v2/containers/${id}/start`)
}

export function stopContainer(id: string) {
  return post(`/api/v2/containers/${id}/stop`)
}

export function restartContainer(id: string) {
  return post(`/api/v2/containers/${id}/restart`)
}

export function removeContainer(id: string) {
  return del(`/api/v2/containers/${id}`)
}

export function getContainerLogs(id: string, params?: { tail?: number; since?: number }) {
  return get<string>(`/api/v2/containers/${id}/logs`, params)
}

export function getContainerStats(id: string) {
  return get<ContainerStats>(`/api/v2/containers/${id}/stats`)
}

export function listImages(params?: PageParams) {
  return get<PageResult<ImageInfo>>('/api/v2/images', params)
}

export function pullImage(image: string, tag?: string) {
  return post('/api/v2/images/pull', { image, tag })
}

export function removeImage(id: string) {
  return del(`/api/v2/images/${id}`)
}

export function listVolumes(params?: PageParams) {
  return get<PageResult<VolumeInfo>>('/api/v2/volumes', params)
}

export function removeVolume(name: string) {
  return del(`/api/v2/volumes/${name}`)
}

export function listNetworks(params?: PageParams) {
  return get<PageResult<NetworkInfo>>('/api/v2/networks', params)
}

export function composeUp(name: string, content: string) {
  return post('/api/v2/compose/up', { name, content })
}

export function composeDown(name: string) {
  return post('/api/v2/compose/down', { name })
}