<template>
  <div class="category">
    <el-card>
      <CategoryCascader v-model:c1="c1Id" v-model:c2="c2Id" v-model:c3="c3Id" />

      <el-button
        v-auth="'btn.Category.add'"
        type="primary"
        :icon="Plus"
        style="margin-top: 12px"
        :disabled="!canAdd"
        @click="openAdd"
      >
        新增{{ addLabel }}
      </el-button>

      <el-table v-loading="loading" :data="tableData" border style="margin-top: 12px">
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="name" label="分类名称" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="`新增${addLabel}`">
      <el-form label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="name" placeholder="请输入分类名称" />
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
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import CategoryCascader from '@/components/CategoryCascader.vue'

type Category2 = components['schemas']['types.Category2']
type Category3 = components['schemas']['types.Category3']

const c1Id = ref('')
const c2Id = ref('')
const c3Id = ref('')

const loading = ref(false)
const tableData = ref<(Category2 | Category3)[]>([])
const name = ref('')
const dialogVisible = ref(false)
const saving = ref(false)

// 当前要新增的层级：选了一级 → 新增二级；选了二级 → 新增三级；选了三级 → 三级为末级不再新增
// 选一级时表格列「该一级下的二级」，选二级时列「该二级下的三级」（即便为空也能直接新增）
const addLevel = computed(() => (c2Id.value ? 3 : c1Id.value ? 2 : 1))
const addLabel = computed(() => (addLevel.value === 2 ? '二级分类' : '三级分类'))
// 一级分类由系统预置，无新增接口；选了一级才能新增二级，选二级才能新增三级
const canAdd = computed(() => addLevel.value >= 2)

async function loadTable() {
  loading.value = true
  try {
    if (c2Id.value) {
      // 选了二级 → 列出该二级下的三级分类（可能为空的初始状态也能直接新增三级）
      tableData.value =
        (
          await client.GET('/api/v1/product/category3/{category2Id}', {
            params: { path: { category2Id: c2Id.value } },
          })
        ).data?.data ?? []
    } else if (c1Id.value) {
      // 选了一级 → 列出该一级下的二级分类
      tableData.value =
        (
          await client.GET('/api/v1/product/category2/{category1Id}', {
            params: { path: { category1Id: c1Id.value } },
          })
        ).data?.data ?? []
    } else {
      tableData.value = []
    }
  } finally {
    loading.value = false
  }
}

watch([c1Id, c2Id, c3Id], loadTable)

function openAdd() {
  name.value = ''
  dialogVisible.value = true
}

async function save() {
  if (!name.value.trim()) {
    ElMessage.warning('请输入名称')
    return
  }
  saving.value = true
  try {
    if (addLevel.value === 2) {
      await client.POST('/api/v1/product/category2', {
        body: { category1Id: c1Id.value, name: name.value },
      })
    } else {
      await client.POST('/api/v1/product/category3', {
        body: { category2Id: c2Id.value, name: name.value },
      })
    }
    ElMessage.success('新增成功')
    dialogVisible.value = false
    await loadTable()
  } finally {
    saving.value = false
  }
}
</script>
