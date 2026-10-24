import { request } from "@/http/axios_n"

interface reqFiles {
  name?: string
  orderKey?: string
  desc?: boolean
  page?: number
  page_size?: number
}

// 字段与后端 /file/list 的 snake_case 返回一一对应
interface fileData {
  id: number
  file_name: string
  full_path: string
  mime: string
  created_at: string
  updated_at: string
}

export interface fileDataModel extends fileData {}

// List
export type fileListData = ListData<fileDataModel[]>

// 分页获取文件信息
export function fileListApi(data: reqFiles) {
  return request<ApiResponseData<fileListData>>({
    url: "/file/list",
    method: "post",
    data
  })
}

// 下载文件（后端按 id 路由：GET /file/download/:id，返回原始字节流）
export function fileDownloadApi(params: { id: number; name: string }) {
  return request({
    url: `/file/download/${params.id}`,
    method: "get",
    responseType: "blob"
  })
}

// 删除文件（后端为 POST /file/delete，body 为 { id }）
export function fileDeleteApi(params: { id: number }) {
  return request<ApiResponseData<null>>({
    url: "/file/delete",
    method: "post",
    data: params
  })
}
