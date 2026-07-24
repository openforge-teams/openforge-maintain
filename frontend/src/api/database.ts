import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface MySQLDB {
  id: number
  name: string
  character_set: string
  collation: string
  size: number
  created_at: string
}

export interface PostgresDB {
  id: number
  name: string
  owner: string
  size: number
  created_at: string
}

export interface RedisInfo {
  version: string
  used_memory: number
  max_memory: number
  connected_clients: number
  uptime: number
  databases: { index: number; keys: number; expires: number }[]
}

export function listMySQL(params?: PageParams) {
  return get<PageResult<MySQLDB>>('/api/v2/databases/mysql', params)
}

export function createMySQLDB(data: { name: string; character_set?: string; user?: string; password?: string }) {
  return post<MySQLDB>('/api/v2/databases/mysql', data)
}

export function deleteMySQLDB(id: number) {
  return del(`/api/v2/databases/mysql/${id}`)
}

export function backupMySQL(id: number) {
  return post(`/api/v2/databases/mysql/${id}/backup`)
}

export function listPostgres(params?: PageParams) {
  return get<PageResult<PostgresDB>>('/api/v2/databases/postgres', params)
}

export function listRedis() {
  return get<RedisInfo>('/api/v2/databases/redis')
}

export function getRedisInfo() {
  return get<RedisInfo>('/api/v2/databases/redis/info')
}