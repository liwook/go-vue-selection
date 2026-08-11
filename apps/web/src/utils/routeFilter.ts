import type { RouteRecordRaw } from 'vue-router'

// 注意：被过滤的 route 是「后端返回的 menu.code 集合」（大写首字母，如 User/Role/Permission/Spu/Sku/Trademark/Attr/Category，见 apps/server/init-sql/init.sql）
// 用 toLowerCase() 兜底，前端 asyncRoutes 里的小写 name 才能匹配上。
export function filterAsyncRoute(asyncRoute: RouteRecordRaw[], routes: string[]) {
  // 小写集合只算一次，递归里复用
  const routeNames = routes.map((item: string) => item.toLowerCase())
  return asyncRoute.filter((item: any) => {
    // ① 如果当前 item 没有 children，直接按 name 过滤（整体路由）
    if (!item.children) {
      if (routeNames.includes(item.name as string)) {
        return true
      }
    } else {
      // ② 否则继续递归；拿到「过滤后」的子路由，再判断父路由是否还该留
      item.children = filterAsyncRoute(item.children, routes) // 递归
      if (item.children && item.children.length > 0) {
        return true
      }
    }
    return false
  })
}
