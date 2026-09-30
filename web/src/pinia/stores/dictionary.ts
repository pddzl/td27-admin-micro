import type { DictModel } from "@/api/sysManagement/dict"
import type { dictDetailDataModel } from "@/api/sysManagement/dictDetail"
import { dictListApi } from "@/api/sysManagement/dict"
import { dictDetailFlatApi } from "@/api/sysManagement/dictDetail"

export const useDictionaryStore = defineStore("dictionary", () => {
  const dictionaries = ref<DictModel[]>([])
  // cache details by dictId
  const detailsMap = ref<Record<number, dictDetailDataModel[]>>({})

  // 并发去重：同一时刻相同请求只发一次
  let dictionariesPromise: Promise<void> | null = null
  const detailsPromises = new Map<number, Promise<void>>()

  const fetchDictionaries = async () => {
    if (dictionaries.value.length > 0) return
    if (!dictionariesPromise) {
      dictionariesPromise = (async () => {
        const res = await dictListApi({})
        if (res.code === 0) {
          dictionaries.value = res.data.list
        }
      })().finally(() => {
        dictionariesPromise = null
      })
    }
    await dictionariesPromise
  }

  const fetchDictionaryDetail = async (dictId: number) => {
    if (detailsMap.value[dictId]) return // ✅ cached
    if (!detailsPromises.has(dictId)) {
      detailsPromises.set(
        dictId,
        (async () => {
          const res = await dictDetailFlatApi({ dictId })
          if (res.code === 0) {
            detailsMap.value[dictId] = res.data
          }
        })().finally(() => {
          detailsPromises.delete(dictId)
        })
      )
    }
    await detailsPromises.get(dictId)
  }

  // ✅ Helper: get options by en_name
  const getOptions = async (en_name: string) => {
    // find dictId
    if (dictionaries.value.length === 0) {
      await fetchDictionaries()
    }

    const dict = dictionaries.value.find(d => d.en_name === en_name)
    if (!dict) return []

    // fetch details if needed
    await fetchDictionaryDetail(dict.id)

    return (detailsMap.value[dict.id] || []).map((item: dictDetailDataModel) => ({
      label: item.label,
      value: item.value
    }))
  }

  /** 字典/字典项管理页在增删改后调用，使缓存失效，下次读取时重新拉取 */
  const invalidate = (dictId?: number | null) => {
    if (dictId == null) {
      dictionaries.value = []
      detailsMap.value = {}
    } else {
      delete detailsMap.value[dictId]
    }
  }

  return {
    dictionaries,
    detailsMap,
    fetchDictionaries,
    fetchDictionaryDetail,
    getOptions,
    invalidate // ✅ expose helper
  }
})
