<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { client } from '@/api'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type ResponseUser = components['schemas']['types.ResponseUser']
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
} = useCrudTable<ResponseUser>(
  (page, limit, username) =>
    client.GET('/api/v1/acl/user', {
      params: { query: { page, limit, username: username || undefined } },
    }),
  async (payload) => {
    // 新增走 POST /acl/user（ParamUserSignUp：username/password/name）
    // 编辑走 PUT  /acl/user（ParamUserUpdate：仅 avatar/name，userId 在行内但 body 不含）
    if ((payload as Partial<ResponseUser>).userId) {
      await client.PUT('/api/v1/acl/user', {
        body: { name: payload.name },
      })
    } else {
      await client.POST('/api/v1/acl/user', {
        body: {
          username: (payload as { username?: string }).username ?? '',
          password: (payload as { password?: string }).password ?? '',
          name: payload.name ?? '',
        },
      })
    }
  },
  (row) =>
    client.DELETE('/api/v1/acl/user/{userId}', { params: { path: { userId: row.userId ?? '' } } }),
  { key: 'username' },
)

onMounted(fetchList)

function addUser() {
  openAdd({ name: '', username: '', password: '' } as Partial<ResponseUser>)
}

async function removeUser(row: ResponseUser) {
  try {
    await ElMessageBox.confirm(`确定删除用户「${row.username ?? ''}」吗？`, '提示', {
      type: 'warning',
    })
    await remove(row)
  } catch {
    /* 用户取消或接口失败（中间件已提示），流程终止 */
  }
}

// —— 分配角色 ——
const assignVisible = ref(false)
const allRoles = ref<Role[]>([])
const checkedRoles = ref<string[]>([])
const currentUserId = ref('')

async function openAssign(row: ResponseUser) {
  currentUserId.value = row.userId ?? ''
  try {
    const { data } = await client.GET('/api/v1/acl/user/{userId}/role', {
      params: { path: { userId: currentUserId.value } },
    })
    allRoles.value = data?.data?.allRolesList ?? []
    checkedRoles.value = (data?.data?.assignRoles ?? []).map((r) => r.roleId ?? '')
    assignVisible.value = true
  } catch {
    /* 中间件已提示 */
  }
}

async function confirmAssign() {
  try {
    await client.POST('/api/v1/acl/user/{userId}/role', {
      params: { path: { userId: currentUserId.value } },
      body: { userId: currentUserId.value, roleIdList: checkedRoles.value },
    })
    ElMessage.success('分配成功')
    assignVisible.value = false
  } catch {
    /* 中间件已提示 */
  }
}
</script>

<template>
  <div class="acl-user">
    <el-card>
      <el-form inline>
        <el-form-item label="用户名">
          <el-input v-model="keyword" placeholder="请输入用户名" clearable @clear="reset" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">
            查询
          </el-button>
          <el-button @click="reset">
            重置
          </el-button>
          <el-button type="primary" @click="addUser">
            新增用户
          </el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="name" label="用户昵称" />
        <el-table-column prop="roleName" label="角色" />
        <el-table-column prop="createTime" label="创建时间" />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" type="warning" @click="openEdit(row)">
              编辑
            </el-button>
            <el-button size="small" type="primary" @click="openAssign(row)">
              分配角色
            </el-button>
            <el-button size="small" type="danger" @click="removeUser(row)">
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
        <el-form-item label="用户名">
          <el-input v-model="(form as any).username" :disabled="title === '编辑'" />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model="(form as any).name" />
        </el-form-item>
        <el-form-item v-if="title === '新增'" label="密码">
          <el-input v-model="(form as any).password" type="password" />
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

    <el-dialog v-model="assignVisible" title="分配角色" width="40%">
      <el-checkbox-group v-model="checkedRoles">
        <el-checkbox v-for="r in allRoles" :key="r.roleId" :value="r.roleId">
          {{ r.roleName }}
        </el-checkbox>
      </el-checkbox-group>
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
