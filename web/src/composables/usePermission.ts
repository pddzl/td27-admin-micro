import { ref } from "vue"
import { getUserButtonsApi } from "@/api/sysManagement/button"

/**
 * 按钮权限的全局缓存（模块级单例）。
 * key 为按钮权限码（如 "user:delete"），value 恒为 true ——
 * 未出现在缓存中的码即为无权限。
 */
export const permissionCache = ref<Record<string, boolean>>({})

let loadingPromise: Promise<void> | null = null

/**
 * 拉取当前用户的所有按钮权限码（幂等，并发去重）。
 * 在路由守卫（登录后/刷新后）调用一次即可，指令与组件直接读缓存。
 */
export function ensurePermissionsLoaded(): Promise<void> {
  if (!loadingPromise) {
    loadingPromise = (async () => {
      try {
        const res = await getUserButtonsApi()
        if (res.code === 0 && Array.isArray(res.data)) {
          const perms: Record<string, boolean> = {}
          res.data.forEach((code) => {
            perms[code] = true
          })
          permissionCache.value = perms
        }
      } catch (error) {
        console.error("加载按钮权限失败", error)
      } finally {
        loadingPromise = null
      }
    })()
  }
  return loadingPromise
}

/** 登出时清空按钮权限缓存，防止账号切换后残留 */
export function resetPermissions() {
  permissionCache.value = {}
  loadingPromise = null
}

/** 同步读取：缓存中存在该码即有权限 */
export function hasPermission(buttonCode: string) {
  return permissionCache.value[buttonCode] === true
}

/** 同步读取：任一码有权限即通过 */
export function hasAnyPermission(buttonCodes: string[]) {
  return buttonCodes.some(code => permissionCache.value[code] === true)
}

/** 同步读取：所有码都有权限才通过 */
export function hasAllPermissions(buttonCodes: string[]) {
  return buttonCodes.every(code => permissionCache.value[code] === true)
}

/**
 * 组件内使用的组合式 API（与指令共享同一份缓存）。
 */
export function usePermission() {
  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    permissionCache
  }
}
