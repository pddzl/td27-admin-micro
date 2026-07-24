import { request } from "@/http/axios_n"

interface dictDetailData {
  label: string
  value: string
  sort: number
  dictId: number | null
  parentId?: number
  children?: dictDetailDataModel[]
  description: string
}

export interface dictDetailDataModel extends dictDetailData, Td27Model {}

// List
export type dictDetailListData = ListData<dictDetailDataModel[]>
export type dictDetailFlatData = dictDetailDataModel[]

interface reqDictDetail extends PageInfo {
  dictId: number
}

export function dictDetailListApi(data: reqDictDetail) {
  return request<ApiResponseData<dictDetailListData>>({
    url: "/dict-detail/list",
    method: "post",
    data: { page: data.page, page_size: data.pageSize, dict_id: data.dictId }
  })
}

export function dictDetailFlatApi(data: { dictId: number }) {
  return request<ApiResponseData<dictDetailFlatData>>({
    url: "/dict-detail/flat",
    method: "post",
    data: { dict_id: data.dictId }
  })
}

export function dictDetailCreateApi(data: dictDetailData) {
  const { dictId, parentId, ...rest } = data
  return request<ApiResponseData<dictDetailDataModel>>({
    url: "/dict-detail/create",
    method: "post",
    data: {
      ...rest,
      dict_id: dictId,
      parent_id: parentId
    }
  })
}

export function dictDetailDeleteApi(data: CId) {
  return request<ApiResponseData<dictDetailDataModel>>({
    url: "/dict-detail/delete",
    method: "post",
    data
  })
}

export function dictDetailUpdateApi(data: dictDetailData & CId) {
  const { dictId, parentId, ...rest } = data
  return request<ApiResponseData<dictDetailDataModel>>({
    url: "/dict-detail/update",
    method: "post",
    data: {
      ...rest,
      dict_id: dictId,
      parent_id: parentId
    }
  })
}
