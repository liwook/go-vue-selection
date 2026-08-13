<script setup lang="ts">
import { ElMessageBox } from 'element-plus'
import { onMounted } from 'vue'
import { client } from '@/api'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type Role = components['schemas']['types.Role']

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
          <el-button type="primary" @click="addRole">
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
            <el-button size="small" type="warning" @click="openEdit(row)">
              编辑
            </el-button>
            <el-button size="small" type="danger" @click="removeRole(row)">
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
  </div>
</template>
