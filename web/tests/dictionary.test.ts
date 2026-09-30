import type { DictModel } from "@/api/sysManagement/dict"
import { createPinia, setActivePinia } from "pinia"
import { beforeEach, describe, expect, it, vi } from "vitest"
import { dictListApi } from "@/api/sysManagement/dict"
import { dictDetailFlatApi } from "@/api/sysManagement/dictDetail"
import { useDictionaryStore } from "@/pinia/stores/dictionary"

vi.mock("@/api/sysManagement/dict", () => ({ dictListApi: vi.fn() }))
vi.mock("@/api/sysManagement/dictDetail", () => ({ dictDetailFlatApi: vi.fn() }))

type DictListRes = Awaited<ReturnType<typeof dictListApi>>
type DetailRes = Awaited<ReturnType<typeof dictDetailFlatApi>>

const dicts: Partial<DictModel>[] = [
  { id: 1, cn_name: "状态", en_name: "status" },
  { id: 2, cn_name: "性别", en_name: "gender" }
]

function dictListOk() {
  return { code: 0, msg: "ok", data: { list: dicts, total: dicts.length } } as unknown as DictListRes
}

function detailOk(list: Array<{ label: string, value: string }>) {
  return { code: 0, msg: "ok", data: list } as unknown as DetailRes
}

beforeEach(() => {
  setActivePinia(createPinia())
  vi.mocked(dictListApi).mockReset()
  vi.mocked(dictDetailFlatApi).mockReset()
})

describe("dictionary store", () => {
  it("fetches dictionaries once and caches options", async () => {
    vi.mocked(dictListApi).mockResolvedValue(dictListOk())
    vi.mocked(dictDetailFlatApi).mockResolvedValue(
      detailOk([
        { label: "启用", value: "1" },
        { label: "禁用", value: "0" }
      ])
    )

    const store = useDictionaryStore()
    const first = await store.getOptions("status")
    const second = await store.getOptions("status")

    expect(dictListApi).toHaveBeenCalledOnce()
    expect(dictDetailFlatApi).toHaveBeenCalledOnce()
    expect(first).toEqual([
      { label: "启用", value: "1" },
      { label: "禁用", value: "0" }
    ])
    expect(second).toEqual(first)
  })

  it("deduplicates concurrent getOptions calls", async () => {
    let resolveDictList!: (value: DictListRes) => void
    vi.mocked(dictListApi).mockImplementation(
      () => new Promise<DictListRes>((resolve) => {
        resolveDictList = resolve
      })
    )
    vi.mocked(dictDetailFlatApi).mockResolvedValue(detailOk([{ label: "男", value: "1" }]))

    const store = useDictionaryStore()
    const p1 = store.getOptions("status")
    const p2 = store.getOptions("gender")

    expect(dictListApi).toHaveBeenCalledOnce()

    resolveDictList(dictListOk())
    await Promise.all([p1, p2])

    expect(dictListApi).toHaveBeenCalledOnce()
    expect(dictDetailFlatApi).toHaveBeenCalledTimes(2) // one per dictId
  })

  it("returns [] for unknown dict en_name", async () => {
    vi.mocked(dictListApi).mockResolvedValue(dictListOk())

    const store = useDictionaryStore()
    const options = await store.getOptions("not_exist")

    expect(options).toEqual([])
    expect(dictDetailFlatApi).not.toHaveBeenCalled()
  })

  it("invalidate() clears cache and refetches", async () => {
    vi.mocked(dictListApi).mockResolvedValue(dictListOk())
    vi.mocked(dictDetailFlatApi).mockResolvedValue(detailOk([{ label: "启用", value: "1" }]))

    const store = useDictionaryStore()
    await store.getOptions("status")
    expect(dictListApi).toHaveBeenCalledOnce()

    store.invalidate()

    await store.getOptions("status")
    expect(dictListApi).toHaveBeenCalledTimes(2)
    expect(dictDetailFlatApi).toHaveBeenCalledTimes(2)
  })

  it("invalidate(dictId) only clears that dict's details", async () => {
    vi.mocked(dictListApi).mockResolvedValue(dictListOk())
    vi.mocked(dictDetailFlatApi).mockResolvedValue(detailOk([{ label: "启用", value: "1" }]))

    const store = useDictionaryStore()
    await store.getOptions("status")
    await store.getOptions("gender")
    expect(dictDetailFlatApi).toHaveBeenCalledTimes(2)

    store.invalidate(1)

    // dictionaries list stays cached, only dict 1 details refetch
    await store.getOptions("status")
    await store.getOptions("gender")
    expect(dictListApi).toHaveBeenCalledOnce()
    expect(dictDetailFlatApi).toHaveBeenCalledTimes(3)
  })
})
