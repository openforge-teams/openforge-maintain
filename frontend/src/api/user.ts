import { get, put, del } from './request'
import type { UserInfo, PageParams, PageResult } from '../env'

export function getProfile() {
  return get<UserInfo>('/api/v2/core/user/profile')
}

export function updateProfile(data: Partial<UserInfo>) {
  return put<UserInfo>('/api/v2/core/user/profile', data)
}

export function listUsers(params: PageParams) {
  return get<PageResult<UserInfo>>('/api/v2/core/users', params)
}

export function getUser(id: number) {
  return get<UserInfo>(`/api/v2/core/users/${id}`)
}

export function updateUser(id: number, data: Partial<UserInfo>) {
  return put<UserInfo>(`/api/v2/core/users/${id}`, data)
}

export function deleteUser(id: number) {
  return del(`/api/v2/core/users/${id}`)
}
