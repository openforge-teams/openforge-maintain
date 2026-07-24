import { post } from './request'

export interface LoginParams {
  username: string
  password: string
  totp_code?: string
}

export interface RegisterParams {
  username: string
  password: string
  email?: string
}

export interface ChangePasswordParams {
  old_password: string
  new_password: string
}

export function login(data: LoginParams) {
  return post<{ access_token: string; refresh_token: string }>('/api/v2/core/auth/login', data)
}

export function logout() {
  return post('/api/v2/core/auth/logout')
}

export function refresh(refresh_token: string) {
  return post<{ token: string }>('/api/v2/core/auth/refresh', { refresh_token })
}

export function register(data: RegisterParams) {
  return post('/api/v2/core/auth/register', data)
}

export function changePassword(data: ChangePasswordParams) {
  return post('/api/v2/core/auth/change-password', data)
}
