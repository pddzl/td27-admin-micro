import { request } from "@/http/axios_n"

// 统计数据（字段与后端 /dashboard/statistics 的 snake_case 返回一一对应）
export interface DashboardStats {
  user_count: number
  role_count: number
  api_count: number
  dept_count: number
  menu_count: number
  dict_count: number
  log_count: number
}

// 最近操作记录（/dashboard/recent-operations 返回 { list: [...] }）
export interface RecentOperation {
  id: number
  ip: string
  method: string
  path: string
  status: number
  user_agent: string
  resp_time: number
  created_at: string
}

// 系统信息（/dashboard/system-info）
export interface SystemInfo {
  go_version: string
  os: string
  arch: string
  cpu_cores: number
  uptime_seconds: number
  goroutine_count: number
  memory_allocated_mb: number
  total_allocated_mb: number
  gc_count: number
}

// 获取统计数据
export function getDashboardStatisticsApi() {
  return request<ApiResponseData<DashboardStats>>({
    url: "/dashboard/statistics",
    method: "get"
  })
}

// 获取最近操作记录
export function getRecentOperationsApi() {
  return request<ApiResponseData<{ list: RecentOperation[] }>>({
    url: "/dashboard/recent-operations",
    method: "get"
  })
}

// 获取系统信息
export function getSystemInfoApi() {
  return request<ApiResponseData<SystemInfo>>({
    url: "/dashboard/system-info",
    method: "get"
  })
}
