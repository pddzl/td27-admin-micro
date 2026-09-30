import type { Directive, DirectiveBinding } from "vue"
import { ensurePermissionsLoaded, permissionCache } from "@/composables/usePermission"

/**
 * v-permission="['user:delete']" —— 无权限时隐藏元素
 * v-permission.disable="['user:delete']" —— 无权限时禁用元素
 *
 * 权限码来自登录后预加载的 user-buttons 缓存（见 usePermission.ts），
 * 缓存就绪后判断是同步的；仅在缓存尚未就绪时等待一次。
 */
async function updateEl(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
  const codes = Array.isArray(binding.value) ? binding.value : [binding.value]
  const modifier = binding.modifiers

  // 缓存未就绪时等待首次加载（幂等，后续调用立即返回）
  await ensurePermissionsLoaded()

  const hasPermission = codes.some(code => permissionCache.value[code] === true)

  if (!hasPermission) {
    if (modifier.disable) {
      // 同时设置原生属性与 pointer-events，保证 element-plus 组件也被禁用
      el.setAttribute("disabled", "true")
      ;(el as any).classList?.add("is-disabled")
      ;(el as HTMLElement).style.pointerEvents = "none"
    } else {
      el.style.display = "none"
    }
  } else {
    el.removeAttribute("disabled")
    ;(el as any).classList?.remove?.("is-disabled")
    ;(el as HTMLElement).style.pointerEvents = ""
    el.style.display = ""
  }
}

export const vPermission: Directive = {
  mounted(el, binding) {
    updateEl(el, binding)
  },
  updated(el, binding) {
    updateEl(el, binding)
  }
}
