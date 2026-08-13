import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import 'element-plus/es/components/message/style/css'

/**
 * 通用分页 CRUD 表格
 * @param listApi   (page, limit, keyword) => Promise<{ data?: {...} }>，内部直接调 client.GET，业务列表在 data.data.records / data.data.total
 * @param saveApi   (payload) => Promise，新增或编辑提交（由 isEdit 决定）
 * @param removeApi (row) => Promise，单条删除
 * @param options.key 用于搜索的关键字字段名（默认 'username'）
 */
export function useCrudTable<T extends object>(
  listApi: (
    page: number,
    limit: number,
    keyword: string,
  ) => Promise<{ data?: { data?: { records?: T[]; total?: number } } }>,
  saveApi: (payload: Partial<T>) => Promise<unknown>,
  removeApi: (row: T) => Promise<unknown>,
  options: { key?: string } = {},
) {
  // ⚠️ 落地提醒：本项目真实实体主键命名各异且为 string（ResponseUser.userId /
  //    ResponseRole.roleId / Menu.menuId），并非统一的 `id`。本封装不假设主键名——
  //    removeApi 直接收到整行 row，由调用方在内部取对应主键（如 row.userId）；
  //    saveApi 也由调用方根据是否存在真实主键判断新增/编辑。
  const key = options.key ?? 'username'

  const list = ref<T[]>([])
  const total = ref(0)
  const pageNo = ref(1)
  const pageSize = ref(5)
  const keyword = ref('')
  const loading = ref(false)

  const dialogVisible = ref(false)
  const title = ref('新增')
  const form = reactive<Partial<T>>({})

  async function fetchList() {
    loading.value = true
    try {
      // listApi 内部直接调 client.GET，成功后 data 一定有值；失败已 throw 到 catch
      const { data } = await listApi(pageNo.value, pageSize.value, keyword.value)
      list.value = data?.data?.records ?? []
      total.value = data?.data?.total ?? 0
    } finally {
      loading.value = false
    }
  }

  function search() {
    pageNo.value = 1
    fetchList()
  }
  function reset() {
    keyword.value = ''
    pageNo.value = 1
    fetchList()
  }
  function openAdd(blank: Partial<T>) {
    title.value = '新增'
    Object.assign(form, blank)
    dialogVisible.value = true
  }
  function openEdit(row: T) {
    title.value = '编辑'
    Object.assign(form, row)
    dialogVisible.value = true
  }
  async function save() {
    try {
      // form 为 reactive<Partial<T>>，运行时即普通对象；此处断言以消除 vue-tsc 对
      // Reactive 深层 readonly 与 Partial<T> 互推的误报
      await saveApi(form as Partial<T>)
      ElMessage.success('操作成功')
      dialogVisible.value = false
      fetchList()
    } catch {
      /* 中间件已 ElMessage.error，仅阻止后续流程 */
    }
  }
  async function remove(row: T) {
    try {
      await removeApi(row)
      ElMessage.success('删除成功')
      if (list.value.length === 1 && pageNo.value > 1) pageNo.value--
      fetchList()
    } catch {
      /* 中间件已提示错误 */
    }
  }

  return {
    key,
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
  }
}
