import { onBeforeUnmount } from "vue"

/**
 * 保证同一时刻只有最新的一个请求在执行：
 * 每次调用 run() 会取消上一次未完成的请求，避免快速翻页/切换筛选时
 * 旧响应覆盖新响应的竞态问题。组件卸载时自动取消未完成请求。
 *
 * 用法：run(signal => userListApi(params, signal))
 * 被取消的请求会 reject（axios CancelledError），由拦截器静默处理。
 */
export function useLatestRequest() {
  let controller: AbortController | null = null

  /** 取消当前未完成的请求 */
  function abort() {
    controller?.abort()
    controller = null
  }

  /** 取消上一个请求并发起新请求 */
  async function run<T>(send: (signal: AbortSignal) => Promise<T>): Promise<T> {
    abort()
    controller = new AbortController()
    return send(controller.signal)
  }

  onBeforeUnmount(abort)

  return { run, abort }
}
