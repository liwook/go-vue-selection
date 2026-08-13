<template>
  <div class="trademark">
    <el-card>
      <el-button type="primary" :icon="Plus" @click="openAdd">
        新增品牌
      </el-button>
      <el-table v-loading="loading" :data="list" border style="margin-top: 12px">
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="tmName" label="品牌名称" />
        <el-table-column label="品牌 Logo">
          <template #default="{ row }">
            <el-image :src="row.logoUrl" style="width: 80px; height: 80px" fit="contain" />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" :icon="Edit" @click="openEdit(row)" />
            <el-popconfirm title="确认删除?" @confirm="remove(row)">
              <template #reference>
                <el-button size="small" type="danger" :icon="Delete" />
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pageNo"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[5, 10, 20]"
        layout="total, sizes, prev, pager, next"
        @current-change="fetchList"
        @size-change="fetchList"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="title">
      <el-form label-width="100px">
        <el-form-item label="品牌名称" required>
          <el-input v-model="form.tmName" placeholder="请输入品牌名称" />
        </el-form-item>
        <el-form-item label="品牌 Logo" required>
          <el-upload
            :show-file-list="false"
            :action="'/api/v1/product/fileUpload'"
            :headers="{ Authorization: `Bearer ${token}` }"
            :on-success="handleLogoSuccess"
            :before-upload="beforeLogoUpload"
          >
            <el-button type="primary">
              点击上传
            </el-button>
          </el-upload>
          <el-image v-if="form.logoUrl" :src="form.logoUrl" style="width: 80px; height: 80px; margin-left: 12px" />
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

<script setup lang="ts">
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type Trademark = components['schemas']['types.Trademark']

const token = localStorage.getItem('token') ?? ''

const {
  list,
  total,
  pageNo,
  pageSize,
  loading,
  dialogVisible,
  title,
  form,
  fetchList,
  openAdd,
  openEdit,
  save,
  remove,
} = useCrudTable<Trademark>(
  (page, limit) => client.GET('/api/v1/product/trademark', { params: { query: { page, limit } } }),
  (payload) => {
    const tmName = payload.tmName ?? ''
    const logoUrl = payload.logoUrl ?? ''
    return payload.tmId
      ? client.PUT('/api/v1/product/trademark/{trademarkId}', {
          params: { path: { trademarkId: String(payload.tmId) } },
          body: { tmId: String(payload.tmId), tmName, logoUrl },
        })
      : client.POST('/api/v1/product/trademark', { body: { tmName, logoUrl } })
  },
  (row) =>
    client.DELETE('/api/v1/product/trademark/{trademarkId}', {
      params: { path: { trademarkId: String(row.tmId) } },
    }),
)

function handleLogoSuccess(res: { data?: string }) {
  if (res.data) form.logoUrl = res.data
}
function beforeLogoUpload(raw: File) {
  if (!['image/jpeg', 'image/png'].includes(raw.type)) {
    ElMessage.error('只能上传 jpg / png 图片')
    return false
  }
  if (raw.size / 1024 > 800) {
    ElMessage.error('图片大小不能超过 800KB')
    return false
  }
  return true
}

fetchList()
</script>
