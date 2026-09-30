import type { RouteRecordRaw } from "vue-router"
import type { MenuData, MenuDataModel } from "@/api/sysManagement/menu"
import { describe, expect, it, vi } from "vitest"
import { dynamicImport, formatRouter } from "@/common/utils/router_m"

function makeMenu(overrides: Partial<MenuData> = {}): MenuDataModel {
  return {
    menu_name: "Home",
    path: "/home",
    component: "Home",
    redirect: "",
    parent_id: 0,
    sort: 0,
    hidden: false,
    keepAlive: false,
    title: "Home",
    affix: false,
    alwaysShow: false,
    ...overrides
  } as MenuDataModel
}

describe("dynamicImport", () => {
  it("returns the Layout loader for Layout", () => {
    const loader = dynamicImport("Layout")
    expect(typeof loader).toBe("function")
  })

  it("falls back to the 404 loader for unknown components", () => {
    const warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {})
    const loader = dynamicImport("no/such/component.vue")
    expect(typeof loader).toBe("function")
    expect(warnSpy).toHaveBeenCalledOnce()
    warnSpy.mockRestore()
  })
})

describe("formatRouter", () => {
  it("maps backend menus to route records without children for leaf nodes", () => {
    const output: RouteRecordRaw[] = []
    formatRouter([makeMenu()], output)

    expect(output).toHaveLength(1)
    expect(output[0].name).toBe("Home")
    expect(output[0].path).toBe("/home")
    expect(output[0].meta?.title).toBe("Home")
    expect(output[0].meta?.keepAlive).toBe(false)
    // leaf nodes should not carry an empty children array
    expect(output[0].children).toBeUndefined()
  })

  it("recursively formats nested children", () => {
    const child = makeMenu({ menu_name: "User", path: "/sys/user", component: "sysManagement/user/index.vue", title: "User" })
    const parent = makeMenu({ menu_name: "Sys", path: "/sys", component: "Layout", children: [child] })
    const output: RouteRecordRaw[] = []
    formatRouter([parent], output)

    expect(output[0].path).toBe("/sys")
    expect(Array.isArray(output[0].children)).toBe(true)
    expect(output[0].children?.[0]?.name).toBe("User")
    expect(output[0].children?.[0]?.meta?.title).toBe("User")
  })

  it("preserves hidden/keepAlive meta flags", () => {
    const output: RouteRecordRaw[] = []
    formatRouter([makeMenu({ hidden: true, keepAlive: true })], output)

    expect(output[0].meta?.hidden).toBe(true)
    expect(output[0].meta?.keepAlive).toBe(true)
  })
})
