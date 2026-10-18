import { request } from "@/http/axios_n"

// 服务令牌响应（list / get 返回的原始字段）
export interface ServiceToken {
  id: number
  name: string
  status: boolean
  expires_at?: number | null
  created_at: string
  updated_at: string
}

// 创建服务令牌响应（raw_token 仅创建时返回一次）
export interface CreateServiceTokenResp {
  token_id: number
  raw_token: string
}

// 服务令牌详情（令牌信息 + 已分配的权限 ID，由两个接口合并）
export interface ServiceTokenDetail extends ServiceToken {
  permission_ids: number[]
}

// 创建服务令牌请求
export interface CreateServiceTokenReq {
  name: string
  expires_at?: number // 过期时间戳(秒)，不传表示永不过期
  permission_ids: number[]
}

// 更新服务令牌请求（权限通过 /service-token/assign-permissions 单独分配）
export interface UpdateServiceTokenReq {
  id: number
  name: string
  status: boolean
  expires_at?: number
}

// 列表查询请求
export interface ListServiceTokenReq extends PageInfo {
  name?: string
  status?: boolean
}

// 创建服务令牌
export function createServiceTokenApi(data: CreateServiceTokenReq) {
  return request<ApiResponseData<CreateServiceTokenResp>>({
    url: "/serviceToken/create",
    method: "post",
    data
  })
}

// 更新服务令牌
export function updateServiceTokenApi(data: UpdateServiceTokenReq) {
  return request<ApiResponseData<null>>({
    url: "/serviceToken/update",
    method: "post",
    data
  })
}

// 删除服务令牌
export function deleteServiceTokenApi(id: number) {
  return request<ApiResponseData<null>>({
    url: "/serviceToken/delete",
    method: "post",
    data: { id }
  })
}

// 获取服务令牌详情（合并令牌信息与已分配权限 ID）
export async function getServiceTokenDetailApi(id: number): Promise<ApiResponseData<ServiceTokenDetail>> {
  const [tokenRes, permRes] = await Promise.all([
    request<ApiResponseData<ServiceToken>>({ url: `/service-token/${id}`, method: "get" }),
    request<ApiResponseData<{ permission_ids: number[] }>>({ url: `/service-token/permissions/${id}`, method: "get" })
  ])
  const data: ServiceTokenDetail | undefined = tokenRes.code === 0 && tokenRes.data
    ? { ...tokenRes.data, permission_ids: permRes.code === 0 ? permRes.data.permission_ids || [] : [] }
    : undefined
  return {
    code: tokenRes.code,
    msg: tokenRes.msg,
    // 失败时后端约定 data 为空，这里以 undefined 表示
    data: data as ServiceTokenDetail
  }
}

// 分配服务令牌权限
export function assignTokenPermissionsApi(data: { token_id: number, permission_ids: number[] }) {
  return request<ApiResponseData<null>>({
    url: "/service-token/assign-permissions",
    method: "post",
    data
  })
}

// 切换服务令牌状态
export function toggleTokenStatusApi(data: { id: number, status: boolean }) {
  return request<ApiResponseData<null>>({
    url: "/service-token/toggle-status",
    method: "post",
    data
  })
}

// 获取服务令牌列表
export function listServiceTokenApi(data: ListServiceTokenReq) {
  return request<ApiResponseData<ApiListData<ServiceToken[]>>>({
    url: "/serviceToken/list",
    method: "post",
    data: {
      name: data.name,
      status: data.status,
      page: data.page,
      page_size: data.pageSize
    }
  })
}
