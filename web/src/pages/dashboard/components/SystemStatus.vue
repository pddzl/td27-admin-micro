<script setup lang="ts">
import type { SystemInfo } from "@/api/sysMonitor/dashboard"
import { CollectionTag, Connection, Cpu, Monitor, Platform, Timer } from "@element-plus/icons-vue"
import { computed } from "vue"

const props = defineProps<{
  info: SystemInfo
}>()

// 系统负载：基于协程数估算（后端暂无负载指标）
const systemLoad = computed(() => {
  const baseLoad = Math.min(props.info.goroutine_count / 10, 100)
  return Math.min(Math.round(baseLoad), 100)
})

// 运行时长：由 uptime_seconds 格式化
const uptimeText = computed(() => {
  const s = props.info.uptime_seconds
  if (!s || s <= 0) return "-"
  const days = Math.floor(s / 86400)
  const hours = Math.floor((s % 86400) / 3600)
  const minutes = Math.floor((s % 3600) / 60)
  if (days > 0) return `${days}天 ${hours}小时`
  if (hours > 0) return `${hours}小时 ${minutes}分钟`
  return `${minutes}分钟`
})

const loadColor = computed(() => {
  if (systemLoad.value < 50) return "#67c23a"
  if (systemLoad.value < 80) return "#e6a23c"
  return "#f56c6c"
})
</script>

<template>
  <div class="system-status">
    <div class="status-header">
      <h3>系统状态</h3>
      <el-tag type="success" size="small">
        运行中
      </el-tag>
    </div>

    <div class="status-content">
      <div class="status-item">
        <div class="status-label">
          <el-icon><Monitor /></el-icon>
          <span>应用名称</span>
        </div>
        <div class="status-value">
          TD27 Admin
        </div>
      </div>

      <div class="status-item">
        <div class="status-label">
          <el-icon><CollectionTag /></el-icon>
          <span>版本</span>
        </div>
        <div class="status-value">
          {{ info.go_version || "-" }}
        </div>
      </div>

      <div class="status-item">
        <div class="status-label">
          <el-icon><Platform /></el-icon>
          <span>运行环境</span>
        </div>
        <div class="status-value">
          {{ info.os }} / {{ info.arch }}
        </div>
      </div>

      <div class="status-item">
        <div class="status-label">
          <el-icon><Cpu /></el-icon>
          <span>CPU 核心</span>
        </div>
        <div class="status-value">
          {{ info.cpu_cores }} 核
        </div>
      </div>

      <div class="status-item">
        <div class="status-label">
          <el-icon><Connection /></el-icon>
          <span>Go 协程</span>
        </div>
        <div class="status-value">
          {{ info.goroutine_count }}
        </div>
      </div>

      <div class="status-item">
        <div class="status-label">
          <el-icon><Timer /></el-icon>
          <span>运行时长</span>
        </div>
        <div class="status-value">
          {{ uptimeText }}
        </div>
      </div>
    </div>

    <div class="status-footer">
      <el-progress
        :percentage="systemLoad"
        :color="loadColor"
        :stroke-width="8"
        :show-text="true"
      />
      <div class="load-text">
        系统负载
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.system-status {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.status-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;

  .el-icon {
    color: #909399;
  }
}

.status-value {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
  font-family: monospace;
}

.status-footer {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.load-text {
  text-align: center;
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}
</style>
