import { get } from './request'

export interface Overview {
  cpu_usage: number
  memory_usage: number
  memory_total: number
  memory_used: number
  disk_usage: number
  disk_total: number
  disk_used: number
  network_in: number
  network_out: number
  container_count: number
  container_running: number
  container_stopped: number
  uptime: number
  hostname: string
  os: string
}

export interface CPUInfo {
  usage: number
  cores: number
  model: string
  history: { time: string; value: number }[]
}

export interface MemoryInfo {
  total: number
  used: number
  free: number
  cached: number
  buffers: number
  history: { time: string; value: number }[]
}

export interface DiskInfo {
  total: number
  used: number
  free: number
  partitions: {
    device: string
    mount: string
    total: number
    used: number
    free: number
    usage: number
  }[]
}

export interface NetworkInfo {
  interfaces: {
    name: string
    rx_bytes: number
    tx_bytes: number
    rx_speed: number
    tx_speed: number
  }[]
  history: { time: string; rx: number; tx: number }[]
}

export interface ProcessInfo {
  pid: number
  name: string
  user: string
  cpu: number
  memory: number
  status: string
  started_at: string
}

export function getOverview() {
  return get<Overview>('/api/v2/dashboard/overview')
}

export function getCPU() {
  return get<CPUInfo>('/api/v2/dashboard/cpu')
}

export function getMemory() {
  return get<MemoryInfo>('/api/v2/dashboard/memory')
}

export function getDisk() {
  return get<DiskInfo>('/api/v2/dashboard/disk')
}

export function getNetwork() {
  return get<NetworkInfo>('/api/v2/dashboard/network')
}

export function getProcesses() {
  return get<ProcessInfo[]>('/api/v2/dashboard/processes')
}