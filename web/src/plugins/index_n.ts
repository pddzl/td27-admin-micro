import type { App } from "vue"
import { vPermission } from "@/directives/permission"
import { installElementPlusIcons } from "./element-plus-icons"
import { installSvgIcon } from "./svg-icon"

export function installPlugins(app: App) {
  installElementPlusIcons(app)
  installSvgIcon(app)
  app.directive("permission", vPermission)
}
