<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { client } from '@/api'
import type { components } from '@/api/schema'

type Menu = components['schemas']['types.Menu']

const menuTree = ref<Menu[]>([])
const dialogVisible = ref(false)
const title = ref('新增')
const form = reactive<Partial<Menu>>({})

async function fetchMenu() {
  try {
    const { data } = await client.GET('/api/v1/acl/permission')
    menuTree.value = data?.data ?? []
  } catch {
    /* 中间件已提示 */
  }
}

onMounted(fetchMenu)

function openAdd(parent?: Menu) {
  title.value = '新增'
  Object.assign(form, {
    menuId: undefined,
    code: '',
    name: '',
    level: parent ? (parent.level ?? 0) + 1 : 1,
    parentId: parent ? (parent.menuId ?? '') : 0,
  } as Partial<Menu>)
  dialogVisible.value = true
}

function openEdit(row: Menu) {
  title.value = '编辑'
  Object.assign(form, row)
  dialogVisible.value = true
}

async function save() {
  try {
    // 新增走 POST /acl/menu（ParamMenuSave：code/level/name/parentId）
    // 编辑走 PUT  /acl/menu/{menuId}（ParamMenuUpdate：含 menuId）
    if (form.menuId) {
      await client.PUT('/api/v1/acl/permission/{permissionId}', {
        params: { path: { permissionId: form.menuId } },
        body: {
          menuId: form.menuId,
          code: form.code ?? '',
          name: form.name ?? '',
          level: form.level ?? 1,
          parentId: String(form.parentId ?? 0),
        },
      })
    } else {
      // 根菜单传 0，子菜单传父级 menuId
      await client.POST('/api/v1/acl/permission', {
        body: {
          code: form.code ?? '',
          name: form.name ?? '',
          level: form.level ?? 1,
          parentId: String(form.parentId ?? 0),
        },
      })
    }
    ElMessage.success('操作成功')
    dialogVisible.value = false
    fetchMenu()
  } catch {
    /* 中间件已提示 */
  }
}

async function removeMenu(row: Menu) {
  try {
    await ElMessageBox.confirm(`确定删除菜单「${row.name ?? ''}」吗？`, '提示', {
      type: 'warning',
    })
    await client.DELETE('/api/v1/acl/permission/{permissionId}', {
      params: { path: { permissionId: row.menuId ?? '' } },
    })
    ElMessage.success('删除成功')
    fetchMenu()
  } catch {
    /* 用户取消或接口失败（中间件已提示），流程终止 */
  }
}
</script>

<template>
  <div class="acl-permission">
    <el-card>
      <el-button v-auth="'btn.Permission.add'" type="primary" @click="openAdd()">
        添加菜单
      </el-button>

      <el-table :data="menuTree" row-key="menuId" border default-expand-all :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="code" label="权限值" />
        <el-table-column prop="level" label="层级" width="100" />
        <el-table-column label="操作" width="320">
          <template #default="{ row }">
            <el-button v-auth="'btn.Permission.update'" size="small" type="warning" @click="openEdit(row)">
              修改
            </el-button>
            <el-button v-auth="'btn.Permission.add'" size="small" type="primary" @click="openAdd(row)">
              添加子菜单
            </el-button>
            <el-button v-auth="'btn.Permission.remove'" size="small" type="danger" @click="removeMenu(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="title" width="40%">
      <el-form label-width="100px">
        <el-form-item label="菜单名">
          <el-input v-model="(form as any).name" />
        </el-form-item>
        <el-form-item label="权限值(code)">
          <el-input v-model="(form as any).code" />
        </el-form-item>
        <el-form-item label="层级">
          <el-input v-model="(form as any).level" type="number" />
        </el-form-item>
        <el-form-item label="父级ID">
          <el-input v-model="(form as any).parentId" type="number" />
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
