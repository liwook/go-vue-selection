<template>
  <div class="attr">
    <el-card>
      <CategoryCascader v-model:c1="c1Id" v-model:c2="c2Id" v-model:c3="c3Id" />

      <el-button
        v-auth="'btn.Attr.add'"
        type="primary"
        :icon="Plus"
        style="margin-top: 12px"
        :disabled="!c3Id"
        @click="openAdd"
      >
        新增属性
      </el-button>

      <el-table v-loading="loading" :data="list" border style="margin-top: 12px">
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="attrName" label="属性名称" />
        <el-table-column label="属性值">
          <template #default="{ row }">
            <el-tag v-for="v in row.attrValueList ?? []" :key="v.attrValueId" style="margin-right: 4px">
              {{ v.valueName }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button v-auth="'btn.Attr.update'" size="small" :icon="Edit" @click="openEdit(row)" />
            <el-popconfirm v-auth="'btn.Attr.remove'" title="确认删除?" @confirm="remove(row)">
              <template #reference>
                <el-button size="small" type="danger" :icon="Delete" />
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="title">
      <el-form label-width="100px">
        <el-form-item label="属性名" required>
          <el-input v-model="form.attrName" placeholder="请输入属性名" />
        </el-form-item>
        <el-form-item label="属性值">
          <div>
            <el-button :icon="Plus" @click="addValue">
              添加属性值
            </el-button>
            <div v-for="(v, i) in form.attrValueList" :key="i" style="margin-top: 8px">
              <el-input v-model="v.valueName" placeholder="属性值" style="width: 200px" />
              <el-button type="danger" :icon="Delete" style="margin-left: 8px" @click="removeValue(i)" />
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" :loading="saving" @click="save">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ref, watch } from 'vue'
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import CategoryCascader from '@/components/CategoryCascader.vue'

type Attr = components['schemas']['types.Attr']

const c1Id = ref('')
const c2Id = ref('')
const c3Id = ref('')

const loading = ref(false)
const list = ref<Attr[]>([])
const title = ref('')
const dialogVisible = ref(false)
const saving = ref(false)
const form = ref<Partial<Attr>>({ attrValueList: [] })

async function loadAttr() {
  if (!c3Id.value) {
    list.value = []
    return
  }
  loading.value = true
  try {
    list.value =
      (
        await client.GET('/api/v1/product/attr/{category1Id}/{category2Id}/{category3Id}', {
          params: {
            path: { category1Id: c1Id.value, category2Id: c2Id.value, category3Id: c3Id.value },
          },
        })
      ).data?.data ?? []
  } finally {
    loading.value = false
  }
}

watch([c1Id, c2Id, c3Id], loadAttr)

function openAdd() {
  form.value = { attrName: '', attrValueList: [], categoryId: c3Id.value }
  title.value = '新增属性'
  dialogVisible.value = true
}
function openEdit(row: Attr) {
  form.value = { ...row, attrValueList: [...(row.attrValueList ?? [])] }
  title.value = '修改属性'
  dialogVisible.value = true
}
function addValue() {
  form.value.attrValueList!.push({ valueName: '' })
}
function removeValue(i: number) {
  form.value.attrValueList!.splice(i, 1)
}

async function save() {
  if (!form.value.attrName?.trim()) {
    ElMessage.warning('请输入属性名')
    return
  }
  saving.value = true
  try {
    await client.POST('/api/v1/product/attr', { body: form.value as Attr })
    ElMessage.success('保存成功')
    dialogVisible.value = false
    await loadAttr()
  } finally {
    saving.value = false
  }
}

async function remove(row: Attr) {
  try {
    await client.DELETE('/api/v1/product/attr/{attrId}', {
      params: { path: { attrId: String(row.attrId) } },
    })
    ElMessage.success('删除成功')
    await loadAttr()
  } catch {
    /* 中间件已提示 */
  }
}
</script>
