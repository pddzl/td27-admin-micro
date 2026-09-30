<script lang="ts" setup>
import type { ElTree as ElTree1 } from "element-plus"
import type { ButtonData } from "@/api/sysManagement/button"
import { ref } from "vue"
import { listButtonApi } from "@/api/sysManagement/button"
import { getRolePermissionsApi } from "@/api/sysManagement/permission"
import { rebuildRolePermissionApi } from "@/api/sysManagement/role_permission"

const props = defineProps({
  id: {
    type: Number,
    default: 0
  }
})

const treeRef = ref<InstanceType<typeof ElTree1>>()

const buttonDefaultProps = {
  label(data: any) {
    return `${data.buttonName}（${data.buttonCode}）`
  }
}

const buttonIds = ref<number[]>([]) // 勾选的按钮 domain_id
const buttonsTreeData = ref<ButtonData[]>([])

function getTreeData() {
  Promise.all([
    listButtonApi({ page: 1, pageSize: 500 }),
    getRolePermissionsApi(props.id)
  ])
    .then(([btnRes, permRes]) => {
      buttonsTreeData.value = btnRes.data.list
      buttonIds.value = (permRes.data || [])
        .filter(p => p.domain === "button")
        .map(p => p.domain_id)
    })
    .catch(() => {})
}
getTreeData()

function UpdateHandle() {
  const domainIds = treeRef.value?.getCheckedKeys() as number[]
  rebuildRolePermissionApi({ role_id: props.id, domain_ids: domainIds, domain: "button" })
    .then((res) => {
      if (res.code === 0 || res.code === 200) {
        ElMessage({ type: "success", message: "更新成功" })
      }
    })
    .catch(() => {})
}
</script>

<template>
  <div>
    <div class="clearfix">
      <el-button type="primary" class="button" @click="UpdateHandle">
        更新
      </el-button>
    </div>
    <div class="tree-content">
      <ElTree
        ref="treeRef"
        :data="buttonsTreeData"
        :default-checked-keys="buttonIds"
        node-key="id"
        highlight-current
        :props="buttonDefaultProps"
        show-checkbox
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.button {
  float: right;
  margin-right: 5%;
}
.tree-content {
  overflow: auto;
  height: calc(100vh - 160px);
  margin-top: 10px;
}
.clearfix::after {
  content: "";
  display: block;
  height: 0;
  clear: both;
  visibility: hidden;
}
</style>
