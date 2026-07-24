import { get, post, del } from './request'
import type { PageParams, PageResult } from '../env'

export interface FirewallRule {
  id: number
  protocol: 'tcp' | 'udp' | 'icmp'
  port: string
  source: string
  action: 'allow' | 'deny'
  enabled: boolean
  comment: string
  created_at: string
}

export interface FirewallStatus {
  enabled: boolean
  default_input: 'allow' | 'deny'
  default_output: 'allow' | 'deny'
}

export function listRules(params?: PageParams) {
  return get<PageResult<FirewallRule>>('/api/v2/firewall/rules', params)
}

export function addRule(data: { protocol: string; port: string; source: string; action: string; comment?: string }) {
  return post<FirewallRule>('/api/v2/firewall/rules', data)
}

export function deleteRule(id: number) {
  return del(`/api/v2/firewall/rules/${id}`)
}

export function enableRule(id: number) {
  return post(`/api/v2/firewall/rules/${id}/enable`)
}

export function disableRule(id: number) {
  return post(`/api/v2/firewall/rules/${id}/disable`)
}

export function getFirewallStatus() {
  return get<FirewallStatus>('/api/v2/firewall/status')
}

export function enableFirewall() {
  return post('/api/v2/firewall/enable')
}

export function disableFirewall() {
  return post('/api/v2/firewall/disable')
}
