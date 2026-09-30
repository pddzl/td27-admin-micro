import { request } from "@/http/axios_n"

export interface PermissionData {
  id: number
  name: string
  domain: string
  resource: string
  action: string
  effect: string
  domain_id: number
}

/** 拉取某角色当前拥有的所有权限行（含 menu/api/button/data 全部域） */
export function getRolePermissionsApi(roleId: number) {
  return request<ApiResponseData<PermissionData[]>>({
    url: "/permission/by-role",
    method: "post",
    data: { id: roleId }
  })
}
