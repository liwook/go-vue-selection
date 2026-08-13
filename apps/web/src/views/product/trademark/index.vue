<template>
  <div class="trademark">
    <el-card>
      <el-table v-loading="loading" :data="list" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="tmName" label="品牌名称" />
        <el-table-column label="品牌 Logo">
          <template #default="{ row }">
            <el-image :src="row.logoUrl" style="width: 80px; height: 80px" fit="contain" />
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" />
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
  </div>
</template>

<script setup lang="ts">
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type Trademark = components['schemas']['types.Trademark']

const { list, total, pageNo, pageSize, loading, fetchList } = useCrudTable<Trademark>(
  (page, limit) => client.GET('/api/v1/product/trademark', { params: { query: { page, limit } } }),
  () => Promise.resolve(),
  () => Promise.resolve(),
)

fetchList()
</script>
