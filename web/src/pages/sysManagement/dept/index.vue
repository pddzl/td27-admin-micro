<script setup lang="ts">
import type { FormInstance, FormRules, TableInstance } from "element-plus"
import type { CreateDeptReq, Dept, UpdateDeptReq } from "@/api/sysManagement/dept"
import { onMounted, reactive, ref } from "vue"
import {
  createDeptApi,
  deleteDeptApi,
  deptListApi,
  getElTreeDeptsApi,
  updateDeptApi
} from "@/api/sysManagement/dept"

// Search form
const searchForm = reactive({
  dept_name: "",
  status: undefined as boolean | undefined
})

// Table data
const loading = ref(false)
const deptTree = ref<Dept[]>([])
const allDepts = ref<Dept[]>([])
const expandedRowKeys = ref<number[]>([])
const tableRef = ref<TableInstance>()

// Dialog
const dialogVisible = ref(false)
const dialogTitle = ref("")
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const deptOptions = ref<Dept[]>([])

const formData = reactive<{
  id?: number
  dept_name: string
  parent_id: number | undefined
  sort: number
  status: boolean
}>({
  dept_name: "",
  parent_id: undefined,
  sort: 0,
  status: true
})

const formRules: FormRules = {
  dept_name: [{ required: true, message: "请输入部门名称", trigger: "blur" }]
}

// Get department list
async function getDeptList() {
  loading.value = true
  try {
    const res = await deptListApi({
      dept_name: searchForm.dept_name,
      status: searchForm.status
    })
    if (res.code === 0) {
      const data = res.data.list ?? res.data
      deptTree.value = data
      allDepts.value = flattenTree(data)
    }
  } finally {
    loading.value = false
  }
}

// Flatten tree for options
function flattenTree(tree: Dept[]): Dept[] {
  const result: Dept[] = []
  const traverse = (nodes: Dept[]) => {
    nodes.forEach((node) => {
      result.push(node)
      if (node.children) {
        traverse(node.children)
      }
    })
  }
  traverse(tree)
  return result
}

// Get department options for tree-select
async function getDeptOptions() {
  try {
    const res = await getElTreeDeptsApi()
    if (res.code === 200) {
      deptOptions.value = res.data.tree
    }
  } catch (error) {
    console.error(error)
  }
}

// Search
function handleSearch() {
  getDeptList()
}

// Reset
function handleReset() {
  searchForm.dept_name = ""
  searchForm.status = undefined
  getDeptList()
}

// Collect all row IDs from a tree
// Expand/Collapse all
function handleExpandAll() {
  if (expandedRowKeys.value.length === 0) {
    const ids: number[] = []
    function collect(tree: Dept[]) {
      for (const n of tree) {
        if (n.id) ids.push(n.id)
        if (n.children?.length) collect(n.children)
      }
    }
    collect(deptTree.value)
    expandedRowKeys.value = ids
  } else {
    expandedRowKeys.value = []
  }
}

// Create
function handleCreate() {
  isEdit.value = false
  dialogTitle.value = "新增部门"
  resetForm()
  getDeptOptions()
  dialogVisible.value = true
}

// Update
function handleUpdate(row: Dept) {
  isEdit.value = true
  dialogTitle.value = "修改部门"
  Object.assign(formData, {
    id: row.id,
    dept_name: row.dept_name,
    parent_id: row.parent_id || undefined,
    sort: row.sort,
    status: row.status
  })
  getDeptOptions()
  dialogVisible.value = true
}

// Add child
function handleAddChild(row: Dept) {
  isEdit.value = false
  dialogTitle.value = "新增子部门"
  resetForm()
  formData.parent_id = row.id
  getDeptOptions()
  dialogVisible.value = true
}

// Delete
function handleDelete(row: Dept) {
  ElMessageBox.confirm(`确定要删除部门 "${row.dept_name}" 吗？`, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  })
    .then(async () => {
      try {
        const res = await deleteDeptApi({ id: row.id! })
        if (res.code === 0) {
          ElMessage.success("删除成功")
          getDeptList()
        }
      } catch (error) {
        console.error(error)
      }
    })
    .catch(() => {})
}

function closeDialog() {
  dialogVisible.value = false
  resetForm()
}

// Submit form
async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          const res = await updateDeptApi(formData as UpdateDeptReq)
          if (res.code === 200) {
            ElMessage.success("更新成功")
            dialogVisible.value = false
            getDeptList()
          }
        } else {
          const res = await createDeptApi(formData as CreateDeptReq)
          if (res.code === 200) {
            ElMessage.success("创建成功")
            dialogVisible.value = false
            getDeptList()
          }
        }
      } catch (error) {
        console.error(error)
      }
      closeDialog()
    }
  })
}

// Reset form
function resetForm() {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(formData, {
    dept_name: "",
    parent_id: undefined,
    sort: 0,
    status: true
  })
}

onMounted(() => {
  getDeptList()
})
</script>

<template>
  <div class="app-container">
    <!-- Search Form -->
    <el-card v-loading="loading" shadow="never" class="search-wrapper">
      <el-form :model="searchForm" :inline="true">
        <el-form-item prop="dept_name">
          <template #label>部门名称</template>
          <el-input v-model="searchForm.dept_name" placeholder="请输入部门名称" clearable />
        </el-form-item>
        <el-form-item prop="status">
          <template #label>状态</template>
          <el-select v-model="searchForm.status" placeholder="部门状态" clearable style="width: 120px">
            <el-option label="正常" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">
            搜索
          </el-button>
          <el-button icon="Refresh" @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Table -->
    <el-card v-loading="loading" shadow="never">
      <div class="toolbar-wrapper">
        <div>
          <el-button type="primary" icon="Plus" @click="handleCreate">
            新增
          </el-button>
          <el-button icon="Expand" @click="handleExpandAll">
            展开/折叠
          </el-button>
        </div>
        <div>
          <el-tooltip content="刷新" effect="light">
            <el-button type="primary" icon="RefreshRight" circle plain />
          </el-tooltip>
        </div>
      </div>
      <div class="table-wrapper">
        <el-table
          ref="tableRef"
          v-loading="loading"
          :data="deptTree"
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          border
          :expand-row-keys="expandedRowKeys"
          highlight-current-row
        >
          <el-table-column prop="dept_name" label="部门名称" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="dept-name">{{ row.dept_name }}</span>
            </template>
          </el-table-column>

          <el-table-column prop="status" label="状态" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status ? 'success' : 'danger'">
                {{ row.status ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="操作" align="center" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link icon="Edit" @click="handleUpdate(row)">
                修改
              </el-button>
              <el-button type="primary" link icon="Plus" @click="handleAddChild(row)">
                新增
              </el-button>
              <el-button
                v-if="!row.children || row.children.length === 0"
                type="danger"
                link
                icon="Delete"
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="30%"
      @closed="resetForm"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="80px">
        <el-form-item prop="parentId">
          <template #label>上级部门</template>
          <el-tree-select
            v-model="formData.parent_id"
            :data="deptOptions"
            :props="{ label: 'dept_name', value: 'id', children: 'children' }"
            placeholder="请选择上级部门"
            clearable
            check-strictly
            :disabled="isEdit && formData.id === 1"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="部门名称" prop="dept_name">
          <el-input v-model="formData.dept_name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item prop="sort">
          <template #label>排序</template>
          <el-input-number v-model="formData.sort" :min="0" :max="999" style="width: 100%" />
        </el-form-item>
        <el-form-item prop="status">
          <template #label>状态</template>
          <el-radio-group v-model="formData.status">
            <el-radio :value="true">
              正常
            </el-radio>
            <el-radio :value="false">
              禁用
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
.search-wrapper {
  margin-bottom: 5px;
  :deep(.el-card__body) {
    padding-bottom: 2px;
  }
}
</style>
