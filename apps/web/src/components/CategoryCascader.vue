<template>
  <div class="category-cascader">
    <el-select
      v-model="c1"
      placeholder="一级分类"
      clearable
      :disabled="disabled"
      @change="onC1Change"
    >
      <el-option v-for="c in c1List" :key="c.category1Id" :label="c.name" :value="c.category1Id" />
    </el-select>
    <el-select
      v-model="c2"
      placeholder="二级分类"
      clearable
      :disabled="disabled || !c1"
      @change="onC2Change"
    >
      <el-option v-for="c in c2List" :key="c.category2Id" :label="c.name" :value="c.category2Id" />
    </el-select>
    <el-select
      v-model="c3"
      placeholder="三级分类"
      clearable
      :disabled="disabled || !c2"
      @change="onC3Change"
    >
      <el-option v-for="c in c3List" :key="c.category3Id" :label="c.name" :value="c.category3Id" />
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { client } from '@/api/client'
import type { components } from '@/api/schema'

type Category1 = components['schemas']['types.Category1']
type Category2 = components['schemas']['types.Category2']
type Category3 = components['schemas']['types.Category3']

const props = defineProps<{ disabled?: boolean }>()
const c1 = defineModel<string>('c1')
const c2 = defineModel<string>('c2')
const c3 = defineModel<string>('c3')

const c1List = ref<Category1[]>([])
const c2List = ref<Category2[]>([])
const c3List = ref<Category3[]>([])

async function loadC1() {
  c1List.value = (await client.GET('/api/v1/product/category1')).data?.data ?? []
}
async function loadC2(category1Id: string) {
  c2List.value = (await client.GET('/api/v1/product/category2/{category1Id}', {
    params: { path: { category1Id } },
  })).data?.data ?? []
}
async function loadC3(category2Id: string) {
  c3List.value = (await client.GET('/api/v1/product/category3/{category2Id}', {
    params: { path: { category2Id } },
  })).data?.data ?? []
}

function onC1Change(val?: string) {
  c2.value = ''
  c3.value = ''
  c2List.value = []
  c3List.value = []
  if (val) loadC2(val)
}
function onC2Change(val?: string) {
  c3.value = ''
  c3List.value = []
  if (val) loadC3(val)
}
function onC3Change() {
  /* c3 变化由 v-model 同步给父组件 */
}

onMounted(loadC1)

// 外部回填：c1/c2/c3 已有值时逐级加载
watch(c1, async (val) => {
  if (val && !c2List.value.length) await loadC2(val)
})
watch(c2, async (val) => {
  if (val && !c3List.value.length) await loadC3(val)
})

// 暴露给父组件在新增完成后重置
defineExpose({ reset: () => { c1.value = ''; c2.value = ''; c3.value = ''; c2List.value = []; c3List.value = [] } })
void props
</script>

<style scoped>
.category-cascader {
  display: flex;
  gap: 8px;
}
.category-cascader .el-select {
  min-width: 140px;
}
</style>
