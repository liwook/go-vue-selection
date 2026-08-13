<template>
  <div class="sku">
    <el-card>
      <el-table v-loading="loading" :data="list" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="skuName" label="SKU 名称" />
        <el-table-column label="价格(元)" width="120">
          <template #default="{ row }">
            {{ row.priceCent != null ? (row.priceCent / 100).toFixed(2) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="重量(克)" width="120">
          <template #default="{ row }">
            {{ row.weightMg != null ? (row.weightMg / 1000).toFixed(0) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isSale === 1 ? 'success' : 'info'">
              {{ row.isSale === 1 ? '已上架' : '已下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button
              v-if="row.isSale !== 1"
              size="small"
              type="success"
              @click="onSale(row)"
            >
              上架
            </el-button>
            <el-button v-else size="small" type="warning" @click="cancelSale(row)">
              下架
            </el-button>
            <el-button size="small" :icon="View" @click="openDetail(row)" />
            <el-button size="small" type="danger" :icon="Delete" @click="remove(row)" />
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

    <el-drawer v-model="detailVisible" title="SKU 详情" size="40%">
      <template v-if="detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="SKU 名称">
            {{ detail.skuName }}
          </el-descriptions-item>
          <el-descriptions-item label="价格(元)">
            {{ detail.priceCent != null ? (detail.priceCent / 100).toFixed(2) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="重量(克)">
            {{ detail.weightMg != null ? (detail.weightMg / 1000).toFixed(0) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            {{ detail.isSale === 1 ? '已上架' : '已下架' }}
          </el-descriptions-item>
        </el-descriptions>

        <el-divider>图片</el-divider>
        <el-image
          v-for="img in detail.skuImageList ?? []"
          :key="img.imgId ?? img.imgUrl"
          :src="img.imgUrl"
          style="width: 80px; height: 80px; margin: 4px"
        />

        <el-divider>平台属性</el-divider>
        <el-tag
          v-for="a in detail.skuAttrValueList ?? []"
          :key="a.skuAttrValueId"
          style="margin: 4px"
        >
          {{ a.attrName }}：{{ a.valueName }}
        </el-tag>

        <el-divider>销售属性</el-divider>
        <el-tag
          v-for="s in detail.skuSaleAttrValueList ?? []"
          :key="s.skuSaleAttrValueId"
          style="margin: 4px"
        >
          {{ s.saleAttrName }}：{{ s.saleAttrValueName }}
        </el-tag>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { Delete, View } from '@element-plus/icons-vue'
import { ref } from 'vue'
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import { useCrudTable } from '@/composables/useCrudTable'

type SkuInfo = components['schemas']['types.SkuInfo']
type ResponseSkuInfo = components['schemas']['types.ResponseSkuInfo']

const { list, total, pageNo, pageSize, loading, fetchList, remove } = useCrudTable<SkuInfo>(
  (page, limit) => client.GET('/api/v1/product/sku', { params: { query: { page, limit } } }),
  () => Promise.resolve(),
  async (row) => {
    await client.DELETE('/api/v1/product/sku/{skuId}', {
      params: { path: { skuId: String(row.skuId) } },
    })
  },
)

const detailVisible = ref(false)
const detail = ref<ResponseSkuInfo | null>(null)

async function openDetail(row: SkuInfo) {
  detail.value =
    (
      await client.GET('/api/v1/product/sku/{skuId}', {
        params: { path: { skuId: String(row.skuId) } },
      })
    ).data?.data ?? null
  detailVisible.value = true
}

async function onSale(row: SkuInfo) {
  await client.PUT('/api/v1/product/sku/{skuId}/onsale', {
    params: { path: { skuId: String(row.skuId) } },
  })
  fetchList()
}

async function cancelSale(row: SkuInfo) {
  await client.PUT('/api/v1/product/sku/{skuId}/cancelsale', {
    params: { path: { skuId: String(row.skuId) } },
  })
  fetchList()
}
</script>
