/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  trace_id: string
}

interface UserInfo {
  id: number
  username: string
  email?: string
  avatar?: string
  role: string
  created_at: string
}

interface PageParams {
  page: number
  page_size: number
}

interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
