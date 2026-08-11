import type { RouteRecordRaw } from 'vue-router'

// 注意：被过滤的 route 是「后端返回的 menu.code 集合」（大写首字母，如 User/Role/Permission/Spu/Sku/Trademark/Attr/Category，见 apps/server/init-sql/init.sql）
// 用 toLowerCase() 兜底，前端 asyncRoutes 里的小写 name 才能匹配上。
export function filterAsyncRoute(asyncRoute: RouteRecordRaw[], routes: string[]): RouteRecordRaw[] {
  // 小写集合只算一次，递归里复用
  const routeNames = routes.map((item: string) => item.toLowerCase())
  return asyncRoute.reduce((acc: RouteRecordRaw[], item: any) => {
    // ① 如果当前 item 没有 children，直接按 name 过滤（整体路由）
    if (!item.children) {
      if (routeNames.includes(item.name as string)) {
        acc.push(item)
      }
    } else {
      // ② 否则继续递归；拿到「过滤后」的子路由，再判断父路由是否还该留
      const filteredChildren = filterAsyncRoute(item.children, routes)
      if (filteredChildren.length > 0) {
        // 不修改原路由对象，生成新的 children 数组
        acc.push({
          ...item,
          children: filteredChildren,
        } as RouteRecordRaw)
      }
    }
    return acc
  }, [])
}
