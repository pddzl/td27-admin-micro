<script lang="ts" setup>
import type { ElTree } from "element-plus"
import type { ApiChild, ApiTreeData } from "@/api/sysManagement/api"
import type { ServiceToken } from "@/api/sysTool/service_token"
import { usePagination } from "@@/composables/usePagination_n"
import { formatDateTime } from "@@/utils/datetime"
import { apiGetElTreeApi } from "@/api/sysManagement/api"
import {
  assignTokenPermissionsApi,
  createServiceTokenApi,
  deleteServiceTokenApi,
  getServiceTokenDetailApi,
  listServiceTokenApi,
  toggleTokenStatusApi,
  updateServiceTokenApi
} from "@/api/sysTool/service_token"

const { paginationData, changeCurrentPage, changePageSize } = usePagination()

// Table data
const tableData = ref<ServiceToken[]>([])
const total = ref(0)
const loading = ref(false)
const searchForm = ref({
  name: "",
  status: undefined as boolean | undefined
})

// Dialog
const dialogVisible = ref(false)
const dialogTitle = ref("")
const isEdit = ref(false)
const saving = ref(false)

// Form
const form = ref<{
  id?: number
  name: string
  status?: boolean
  expiresAt?: number
  apiIds: number[] // 已分配的 API 权限 ID（domain_id）
}>({
  name: "",
  expiresAt: undefined,
  apiIds: [],
  status: true
})
const formRef = ref()

// API Tree
const apiTreeRef = ref<InstanceType<typeof ElTree>>()
const apiTreeData = ref<ApiTreeData[]>([])
const apiTreeLoading = ref(false)

// Generated token (only shown once after creation)
const generatedToken = ref("")
const showTokenDialog = ref(false)

// Rules
const rules = {
  name: [
    { required: true, message: "请输入令牌名称", trigger: "blur" },
    { max: 100, message: "长度不能超过100个字符", trigger: "blur" }
  ]
}

// Load table data
async function loadData() {
  loading.value = true
  try {
    const res = await listServiceTokenApi({
      page: paginationData.currentPage,
      pageSize: paginationData.pageSize,
      name: searchForm.value.name || undefined,
      status: searchForm.value.status
    })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

// Load API tree（id=0 表示不按角色预勾选，令牌的已选权限由 /service-token/permissions/:id 提供）
async function loadApiTree() {
  apiTreeLoading.value = true
  try {
    const res = await apiGetElTreeApi({ id: 0, from_source: "token" })
    if (res.code === 0) {
      apiTreeData.value = res.data.list || []
    } else {
      ElMessage.error(res.msg || "加载API列表失败")
    }
  } catch (err) {
    console.error("Load API tree error:", err)
    ElMessage.error("加载API列表失败")
  } finally {
    apiTreeLoading.value = false
  }
}

// Search
function handleSearch() {
  paginationData.currentPage = 1
  paginationData.pageSize = 10
  loadData()
}

// 分页
function handleSizeChange(value: number) {
  changePageSize(value)
  loadData()
}

function handleCurrentChange(value: number) {
  changeCurrentPage(value)
  loadData()
}

// Reset
function handleReset() {
  searchForm.value = { name: "", status: undefined }
  handleSearch()
}

// Add new
function handleAdd() {
  isEdit.value = false
  dialogTitle.value = "创建服务令牌"
  form.value = {
    name: "",
    expiresAt: undefined,
    apiIds: [],
    status: true
  }
  dialogVisible.value = true
  loadApiTree()
}

// Edit
async function handleEdit(row: ServiceToken) {
  isEdit.value = true
  dialogTitle.value = "编辑服务令牌"
  dialogVisible.value = true
  loadApiTree()

  const res = await getServiceTokenDetailApi(row.id)
  if (res.code === 0) {
    form.value.id = res.data.id
    form.value.name = res.data.name
    form.value.status = res.data.status
    form.value.expiresAt = res.data.expires_at || undefined
    form.value.apiIds = res.data.permission_ids || []
  }
}

// 已分配权限 ID -> API 树节点 key（path:method）映射，用于 default-checked-keys
const checkedTreeKeys = computed<string[]>(() => {
  const keyById = new Map<number, string>()
  for (const group of apiTreeData.value) {
    for (const api of group.children || []) {
      if (api.id) {
        keyById.set(api.id, api.key)
      }
    }
  }
  return (form.value.apiIds || [])
    .map(id => keyById.get(id))
    .filter((key): key is string => Boolean(key))
})

// Delete
function handleDelete(row: ServiceToken) {
  ElMessageBox.confirm(`确定要删除令牌 "${row.name}" 吗？删除后将无法恢复。`, "提示", {
    type: "warning"
  }).then(async () => {
    const res = await deleteServiceTokenApi(row.id)
    if (res.code === 0) {
      ElMessage.success("删除成功")
      loadData()
    }
  })
}

// Get checked API IDs from tree
function getCheckedApiIds(): number[] {
  const checkedNodes = apiTreeRef.value?.getCheckedNodes(false, true) as ApiChild[]
  if (!checkedNodes) return []

  return checkedNodes
    .filter(node => node.id && node.path && node.method)
    .map(node => node.id)
}

// 后端语义不一致：创建时 expires_at 为「自now起的秒数(TTL)」，更新时为绝对时间戳
function relativeTtl(expiresAt?: number): number | undefined {
  if (!expiresAt) return undefined
  const delta = Math.floor(expiresAt - Date.now() / 1000)
  return delta > 0 ? delta : undefined
}

// Submit
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (form.value.expiresAt && form.value.expiresAt <= Date.now() / 1000) {
    ElMessage.warning("过期时间必须晚于当前时间")
    return
  }

  saving.value = true
  try {
    const permissionIds = getCheckedApiIds()

    if (isEdit.value && form.value.id) {
      const res = await updateServiceTokenApi({
        id: form.value.id,
        name: form.value.name,
        status: form.value.status!,
        expires_at: form.value.expiresAt ? Number(form.value.expiresAt) : undefined
      })
      if (res.code === 0) {
        // 权限单独分配
        const pRes = await assignTokenPermissionsApi({
          token_id: form.value.id,
          permission_ids: permissionIds
        })
        if (pRes.code !== 0) {
          ElMessage.error(pRes.msg || "权限分配失败")
          return
        }
        ElMessage.success("更新成功")
        dialogVisible.value = false
        loadData()
      }
    } else {
      const res = await createServiceTokenApi({
        name: form.value.name,
        expires_at: relativeTtl(form.value.expiresAt),
        permission_ids: permissionIds
      })
      if (res.code === 0) {
        generatedToken.value = res.data.raw_token
        showTokenDialog.value = true
        dialogVisible.value = false
        loadData()
      }
    }
  } finally {
    saving.value = false
  }
}

// Copy token
function copyToken() {
  navigator.clipboard.writeText(generatedToken.value)
  ElMessage.success("已复制到剪贴板")
}

// Status change
async function handleStatusChange(row: ServiceToken) {
  const res = await toggleTokenStatusApi({ id: row.id, status: row.status })
  if (res.code === 0) {
    ElMessage.success("状态更新成功")
  } else {
    row.status = !row.status
  }
}

// Format expiration
function formatExpires(expiresAt?: number | null) {
  if (!expiresAt) return "永不过期"
  return formatDateTime(expiresAt)
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="app-container">
    <el-card class="mb-[5px]" shadow="never">
      <el-form :model="searchForm" inline>
        <el-form-item label="令牌名称">
          <el-input v-model="searchForm.name" placeholder="请输入" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            搜索
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="flex justify-between items-center">
          <el-button v-permission="['token:create']" type="primary" icon="Plus" @click="handleAdd">
            创建令牌
          </el-button>
          <el-tooltip content="刷新" effect="light">
            <el-button type="primary" icon="RefreshRight" circle plain @click="loadData" />
          </el-tooltip>
        </div>
      </template>

      <el-table v-loading="loading" :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="令牌名称" min-width="120px" />
        <el-table-column label="状态" width="100px">
          <template #default="{ row }">
            <el-switch v-model="row.status" @change="handleStatusChange(row)" />
          </template>
        </el-table-column>
        <el-table-column label="过期时间" width="180">
          <template #default="{ row }">
            {{ formatExpires(row.expires_at) }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="['token:update']" type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button v-permission="['token:delete']" type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="mt-5 flex justify-end">
        <el-pagination
          background
          :layout="paginationData.layout"
          :page-sizes="paginationData.pageSizes"
          :total="paginationData.total"
          v-model:page-size="paginationData.pageSize"
          v-model:current-page="paginationData.currentPage"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="令牌名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：第三方系统集成" />
        </el-form-item>

        <el-form-item label="过期时间">
          <el-date-picker
            v-model="form.expiresAt"
            type="datetime"
            placeholder="选择过期时间（留空永不过期）"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="X"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="API权限">
          <div class="api-tree-container">
            <el-tree
              v-if="!apiTreeLoading && apiTreeData.length > 0"
              ref="apiTreeRef"
              :data="apiTreeData"
              node-key="key"
              :default-checked-keys="checkedTreeKeys"
              default-expand-all
              show-checkbox
              :props="{
                label: (data: any) => data.description || data.key,
                children: 'children',
              }"
            />
            <el-empty v-else-if="!apiTreeLoading && apiTreeData.length === 0" description="暂无API数据，请先创建API" />
            <el-skeleton v-else :rows="6" animated />
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" :loading="saving" @click="handleSubmit">
          保存
        </el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="showTokenDialog"
      title="服务令牌创建成功"
      width="600px"
      :close-on-click-modal="false"
      :show-close="false"
    >
      <el-alert
        title="请立即复制保存"
        description="此令牌仅显示一次，关闭后将无法再次查看。请妥善保管。"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <div class="token-display">
        <code>{{ generatedToken }}</code>
        <el-button type="primary" @click="copyToken">
          复制
        </el-button>
      </div>

      <div class="token-usage">
        <h4>使用说明</h4>
        <p>在请求头中添加以下字段：</p>
        <pre>X-Service-Token: {{ generatedToken }}</pre>
      </div>

      <template #footer>
        <el-button type="primary" @click="showTokenDialog = false">
          我已保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
.api-tree-container {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
}

.token-display {
  display: flex;
  gap: 10px;
  align-items: center;
  background: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;

  code {
    flex: 1;
    font-size: 14px;
    word-break: break-all;
    font-family: "Courier New", monospace;
  }
}

.token-usage {
  background: #f5f7fa;
  padding: 15px;
  border-radius: 4px;

  h4 {
    margin-top: 0;
    margin-bottom: 10px;
  }

  pre {
    background: #2d2d2d;
    color: #f8f8f2;
    padding: 10px;
    border-radius: 4px;
    overflow-x: auto;
  }
}
</style>
