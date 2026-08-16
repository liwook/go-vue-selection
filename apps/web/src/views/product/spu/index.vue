<template>
  <div class="spu">
    <!-- 场景 0：SPU 列表 -->
    <el-card v-if="scene === 0">
      <CategoryCascader v-model:c1="c1Id" v-model:c2="c2Id" v-model:c3="c3Id" />
      <el-button
        v-auth="'btn.Spu.add'"
        type="primary"
        :icon="Plus"
        style="margin-top: 12px"
        :disabled="!c3Id"
        @click="openSpuAdd"
      >
        新增 SPU
      </el-button>
      <el-table v-loading="loading" :data="list" border style="margin-top: 12px">
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="spuName" label="SPU 名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button v-auth="'btn.Spu.addsku'" size="small" type="primary" @click="openAddSku(row)">
              添加 SKU
            </el-button>
            <el-button v-auth="'btn.Spu.update'" size="small" :icon="Edit" @click="openSpuEdit(row)" />
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

    <!-- 场景 1：SPU 表单 -->
    <el-card v-if="scene === 1">
      <el-form label-width="100px">
        <el-form-item label="SPU 名称" required>
          <el-input v-model="spuForm.spuName" placeholder="请输入 SPU 名称" />
        </el-form-item>
        <el-form-item label="品牌" required>
          <el-select v-model="spuForm.tmId" placeholder="请选择品牌">
            <el-option v-for="t in tmList" :key="t.tmId" :label="t.tmName" :value="t.tmId" />
          </el-select>
        </el-form-item>
        <el-form-item label="SPU 描述">
          <el-input v-model="spuForm.description" type="textarea" />
        </el-form-item>
        <el-form-item label="SPU 图片">
          <el-upload
            list-type="picture-card"
            :action="'/api/v1/product/fileUpload'"
            :headers="{ Authorization: `Bearer ${token}` }"
            :on-success="handleSpuImg"
            :before-upload="beforeImg"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="销售属性">
          <el-button :icon="Plus" @click="addSaleAttr">
            添加销售属性
          </el-button>
          <div v-for="(sa, i) in spuForm.spuSaleAttrList" :key="i" style="margin-top: 8px">
            <el-select v-model="sa.baseSaleAttrId" placeholder="选择属性" style="width: 180px">
              <el-option
                v-for="b in baseSaleAttrList"
                :key="b.saleAttrId"
                :label="b.name"
                :value="b.saleAttrId"
              />
            </el-select>
            <el-tag
              v-for="(v, j) in sa.spuSaleAttrValueList"
              :key="j"
              closable
              style="margin-left: 8px"
              @close="sa.spuSaleAttrValueList?.splice(j, 1)"
            >
              {{ v.saleAttrValueName }}
            </el-tag>
            <el-input
              v-model="saleNewValues[i]"
              placeholder="输入值回车"
              style="width: 140px; margin-left: 8px"
              @keyup.enter="addSaleAttrValue(i)"
            />
          </div>
        </el-form-item>
      </el-form>
      <el-button type="primary" :loading="saving" @click="saveSpu">
        保存
      </el-button>
      <el-button @click="scene = 0">
        取消
      </el-button>
    </el-card>

    <!-- 场景 2：SKU 表单 -->
    <el-card v-if="scene === 2">
      <el-form label-width="100px">
        <el-form-item label="SKU 名称" required>
          <el-input v-model="skuForm.skuName" placeholder="请输入 SKU 名称" />
        </el-form-item>
        <el-form-item label="价格(元)" required>
          <el-input v-model="priceYuan" placeholder="单位：元" />
        </el-form-item>
        <el-form-item label="重量(克)" required>
          <el-input v-model="weightGram" placeholder="单位：克" />
        </el-form-item>
        <el-form-item label="SKU 图片">
          <el-image
            v-for="img in skuForm.skuImageList"
            :key="img.imgUrl"
            :src="img.imgUrl"
            style="width: 80px; height: 80px; margin: 4px; cursor: pointer"
            :class="{ active: img.imgUrl === defaultSkuImg }"
            @click="setDefaultImg(String(img.imgUrl))"
          />
        </el-form-item>
        <el-form-item v-if="attrList.length" label="平台属性">
          <div v-for="attr in attrList" :key="attr.attrId">
            <span style="margin-right: 8px">{{ attr.attrName }}：</span>
            <el-select
              v-model="skuAttrMap[String(attr.attrId)]"
              placeholder="请选择"
              style="width: 180px"
              @change="syncSkuAttr"
            >
              <el-option
                v-for="v in attr.attrValueList"
                :key="v.attrValueId"
                :label="v.valueName"
                :value="v.attrValueId"
              />
            </el-select>
          </div>
        </el-form-item>
        <el-form-item v-if="spuForm.spuSaleAttrList?.length" label="销售属性">
          <div v-for="sa in spuForm.spuSaleAttrList" :key="sa.baseSaleAttrId">
            <span style="margin-right: 8px">{{ sa.saleAttrName }}：</span>
            <el-select
              v-model="skuSaleMap[String(sa.baseSaleAttrId)]"
              placeholder="请选择"
              style="width: 180px"
              @change="syncSkuSale"
            >
              <el-option
                v-for="v in sa.spuSaleAttrValueList"
                :key="v.saleAttrValueId"
                :label="v.saleAttrValueName"
                :value="v.saleAttrValueId"
              />
            </el-select>
          </div>
        </el-form-item>
      </el-form>
      <el-button type="primary" :loading="saving" @click="saveSku">
        保存
      </el-button>
      <el-button @click="scene = 0">
        取消
      </el-button>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Edit, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { client } from '@/api/client'
import type { components } from '@/api/schema'
import CategoryCascader from '@/components/CategoryCascader.vue'
import { useCrudTable } from '@/composables/useCrudTable'

type Spu = components['schemas']['types.Spu']
type Trademark = components['schemas']['types.Trademark']
type SaleAttr = components['schemas']['types.SaleAttr']
type SpuImage = components['schemas']['types.SpuImage']
type Attr = components['schemas']['types.Attr']
type SkuInfo = components['schemas']['types.SkuInfo']
type SkuImgDTO = components['schemas']['types.SkuImgDTO']

const token = localStorage.getItem('token') ?? ''

const c1Id = ref('')
const c2Id = ref('')
const c3Id = ref('')
const scene = ref(0)
const saving = ref(false)

const { list, total, pageNo, pageSize, loading, fetchList } = useCrudTable<Spu>(
  async (page, limit) => {
    if (!c3Id.value) return { data: { data: { records: [], total: 0 } } }
    return client.GET('/api/v1/product/spu', {
      params: { query: { page, limit, category3Id: Number(c3Id.value) } },
    })
  },
  () => Promise.resolve(),
  () => Promise.resolve(),
)

// SPU 表单
const spuForm = reactive<Partial<Spu>>({ spuSaleAttrList: [] })
const tmList = ref<Trademark[]>([])
const baseSaleAttrList = ref<SaleAttr[]>([])
const spuImgList = ref<SpuImage[]>([])
const saleNewValues = ref<string[]>([])

// SKU 表单
const skuForm = reactive<Partial<SkuInfo>>({
  skuImageList: [],
  skuAttrValueList: [],
  skuSaleAttrValueList: [],
})
const skuImgList = ref<SpuImage[]>([])
const attrList = ref<Attr[]>([])
const priceYuan = ref('')
const weightGram = ref('')
const defaultSkuImg = ref('')
const skuAttrMap = reactive<Record<string, string>>({})
const skuSaleMap = reactive<Record<string, string>>({})

async function loadTmList() {
  tmList.value = (await client.GET('/api/v1/product/trademark/all')).data?.data ?? []
}

async function openSpuAdd() {
  Object.assign(spuForm, {
    spuName: '',
    description: '',
    tmId: '',
    category3Id: c3Id.value,
    spuImageList: [],
    spuSaleAttrList: [],
  })
  spuImgList.value = []
  saleNewValues.value = []
  await loadTmList()
  baseSaleAttrList.value = (await client.GET('/api/v1/product/baseSaleAttr')).data?.data ?? []
  scene.value = 1
}

async function openSpuEdit(row: Spu) {
  Object.assign(spuForm, { ...row, spuSaleAttrList: [...(row.spuSaleAttrList ?? [])] })
  spuImgList.value =
    (
      await client.GET('/api/v1/product/spu/{spuId}/images', {
        params: { path: { spuId: String(row.spuId) } },
      })
    ).data?.data ?? []
  saleNewValues.value = (row.spuSaleAttrList ?? []).map(() => '')
  await loadTmList()
  baseSaleAttrList.value = (await client.GET('/api/v1/product/baseSaleAttr')).data?.data ?? []
  scene.value = 1
}

function addSaleAttr() {
  spuForm.spuSaleAttrList?.push({ baseSaleAttrId: '', saleAttrName: '', spuSaleAttrValueList: [] })
  saleNewValues.value.push('')
}
function addSaleAttrValue(i: number) {
  const v = saleNewValues.value[i]?.trim()
  if (!v) return
  spuForm.spuSaleAttrList?.[i]?.spuSaleAttrValueList?.push({ saleAttrValueName: v })
  saleNewValues.value[i] = ''
}

async function saveSpu() {
  if (!spuForm.spuName?.trim() || !spuForm.tmId) {
    ElMessage.warning('请填写 SPU 名称与品牌')
    return
  }
  saving.value = true
  try {
    spuForm.spuImageList = spuImgList.value
    if (spuForm.spuId) {
      await client.PUT('/api/v1/product/spu/{spuId}', {
        params: { path: { spuId: String(spuForm.spuId) } },
        body: spuForm as Spu,
      })
    } else {
      await client.POST('/api/v1/product/spu', { body: spuForm as Spu })
    }
    ElMessage.success('保存成功')
    scene.value = 0
    fetchList()
  } finally {
    saving.value = false
  }
}

async function openAddSku(row: Spu) {
  Object.assign(skuForm, {
    spuId: row.spuId,
    tmId: row.tmId,
    category3Id: row.category3Id,
    skuImageList: [],
    skuAttrValueList: [],
    skuSaleAttrValueList: [],
  })
  skuImgList.value =
    (
      await client.GET('/api/v1/product/spu/{spuId}/images', {
        params: { path: { spuId: String(row.spuId) } },
      })
    ).data?.data ?? []
  skuForm.skuImageList = skuImgList.value.map(
    (x) => ({ imgUrl: x.imgUrl, imgName: x.imgName ?? '' }) as SkuImgDTO,
  )
  defaultSkuImg.value = skuImgList.value[0]?.imgUrl ?? ''
  attrList.value =
    (
      await client.GET('/api/v1/product/attr/{category1Id}/{category2Id}/{category3Id}', {
        params: {
          path: { category1Id: c1Id.value, category2Id: c2Id.value, category3Id: c3Id.value },
        },
      })
    ).data?.data ?? []
  for (const k of Object.keys(skuAttrMap)) delete skuAttrMap[k]
  for (const k of Object.keys(skuSaleMap)) delete skuSaleMap[k]
  scene.value = 2
}

function syncSkuAttr() {
  skuForm.skuAttrValueList = Object.entries(skuAttrMap).map(([attrId, valueId]) => {
    const attr = attrList.value.find((a) => a.attrId === attrId)
    const val = attr?.attrValueList?.find((v) => v.attrValueId === valueId)
    return { attrId, attrName: attr?.attrName, valueId, valueName: val?.valueName }
  })
}
function syncSkuSale() {
  skuForm.skuSaleAttrValueList = Object.entries(skuSaleMap).map(([saleAttrId, saleAttrValueId]) => {
    const sa = spuForm.spuSaleAttrList?.find((s) => s.baseSaleAttrId === saleAttrId)
    const val = sa?.spuSaleAttrValueList?.find((v) => v.saleAttrValueId === saleAttrValueId)
    return {
      saleAttrId,
      saleAttrName: sa?.saleAttrName,
      saleAttrValueId,
      saleAttrValueName: val?.saleAttrValueName,
    }
  })
}

async function saveSku() {
  if (!skuForm.skuName?.trim()) {
    ElMessage.warning('请填写 SKU 名称')
    return
  }
  skuForm.priceCent = Math.round(Number(priceYuan.value) * 100)
  skuForm.weightMg = Math.round(Number(weightGram.value) * 1000)
  saving.value = true
  try {
    await client.POST('/api/v1/product/sku', { body: skuForm as SkuInfo })
    ElMessage.success('保存成功')
    scene.value = 0
    fetchList()
  } finally {
    saving.value = false
  }
}

function setDefaultImg(url: string) {
  defaultSkuImg.value = url
}

function handleSpuImg(res: { data?: string }) {
  if (res.data) spuImgList.value.push({ imgUrl: res.data, imgName: '' })
}
function beforeImg(raw: File) {
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
</script>
