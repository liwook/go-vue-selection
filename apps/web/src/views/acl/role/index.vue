<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { client } from '@/api'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type Role = components['schemas']['types.Role']
type Menu = components['schemas']['types.Menu']

const {
  list,
  total,
  pageNo,
  pageSize,
  keyword,
  loading,
  dialogVisible,
  title,
  form,
  fetchList,
  search,
  reset,
  openAdd,
  openEdit,
  save,
  remove,
} = useCrudTable<Role>(
  (page, limit, roleName) =>
    client.GET('/api/v1/acl/role', {
      params: { query: { page, limit, roleName: roleName || undefined } },
    }),
  async (payload) => {
    // 新增走 POST /acl/role（ParamRoleSave：roleName/remark）
    // 编辑走 PUT  /acl/role（ParamRoleUpdate：roleId/roleName/remark）
    if ((payload as Partial<Role>).roleId) {
      await client.PUT('/api/v1/acl/role/{roleId}', {
        params: { path: { roleId: payload.roleId ?? '' } },
        body: {
          roleId: payload.roleId,
          roleName: payload.roleName,
          remark: payload.remark,
        },
      })
    } else {
      await client.POST('/api/v1/acl/role', {
        body: { roleName: payload.roleName ?? '', remark: payload.remark },
      })
    }
  },
  (row) =>
    client.DELETE('/api/v1/acl/role/{roleId}', { params: { path: { roleId: row.roleId ?? '' } } }),
  { key: 'roleName' },
)

onMounted(fetchList)

function addRole() {
  openAdd({ roleName: '', remark: '' } as Partial<Role>)
}

async function removeRole(row: Role) {
  try {
    await ElMessageBox.confirm(`确定删除角色「${row.roleName ?? ''}」吗？`, '提示', {
      type: 'warning',
    })
    await remove(row)
  } catch {
    /* 用户取消或接口失败（中间件已提示），流程终止 */
  }
}

// —— 分配权限（菜单树）——
const assignVisible = ref(false)
const menuTree = ref<Menu[]>([])
const checkedKeys = ref<string[]>([])
const currentRoleId = ref('')
const treeRef = ref()

// 递归收集所有 select===true 的 menuId（含父节点，后端去重）
function collectSelected(nodes: Menu[] | undefined, acc: string[]): string[] {
  for (const n of nodes ?? []) {
    if (n.select) acc.push(n.menuId ?? '')
    collectSelected(n.children, acc)
  }
  return acc
}

async function openAssign(row: Role) {
  currentRoleId.value = row.roleId ?? ''
  try {
    const { data } = await client.GET('/api/v1/acl/permission/role/{roleId}', {
      params: { path: { roleId: currentRoleId.value } },
    })
    const tree = data?.data ?? []
    menuTree.value = tree
    checkedKeys.value = collectSelected(tree, [])
    assignVisible.value = true
  } catch {
    /* 中间件已提示 */
  }
}

async function confirmAssign() {
  // 取树中当前勾选（含半选父节点）
  const half = treeRef.value?.getHalfCheckedKeys?.() ?? []
  const checked = treeRef.value?.getCheckedKeys?.() ?? []
  const ids = [...half, ...checked].map(String)
  try {
    await client.POST('/api/v1/acl/permission/role/{roleId}', {
      params: { path: { roleId: currentRoleId.value }, query: { permissionId: ids.join(',') } },
    })
    ElMessage.success('分配成功')
    assignVisible.value = false
  } catch {
    /* 中间件已提示 */
  }
}
</script>

<template>
  <div class="acl-role">
    <el-card>
      <el-form inline>
        <el-form-item label="角色名">
          <el-input v-model="keyword" placeholder="请输入角色名" clearable @clear="reset" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">
            查询
          </el-button>
          <el-button @click="reset">
            重置
          </el-button>
          <el-button type="primary" v-auth="'btn.Role.add'" @click="addRole">
            新增角色
          </el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="roleName" label="角色名称" />
        <el-table-column prop="remark" label="角色备注" />
        <el-table-column prop="createTime" label="创建时间" />
        <el-table-column label="操作" width="320">
          <template #default="{ row }">
            <el-button size="small" type="warning" v-auth="'btn.Role.update'" @click="openEdit(row)">
              编辑
            </el-button>
            <el-button size="small" type="primary" v-auth="'btn.Role.update'" @click="openAssign(row)">
              分配权限
            </el-button>
            <el-button size="small" type="danger" v-auth="'btn.Role.remove'" @click="removeRole(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pageNo"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[5, 10, 20]"
        layout="prev, pager, next, jumper, ->, sizes, total"
        @current-change="fetchList"
        @size-change="fetchList"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="title" width="40%">
      <el-form label-width="100px">
        <el-form-item label="角色名称">
          <el-input v-model="(form as any).roleName" />
        </el-form-item>
        <el-form-item label="角色备注">
          <el-input v-model="(form as any).remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="save">
          确定
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="assignVisible" title="分配权限" width="40%">
      <el-tree
        ref="treeRef"
        :data="menuTree"
        :props="{ label: 'name', children: 'children' }"
        node-key="menuId"
        show-checkbox
        default-expand-all
        :default-checked-keys="checkedKeys"
      />
      <template #footer>
        <el-button @click="assignVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="confirmAssign">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>
