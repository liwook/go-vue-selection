import type { RouteRecordRaw } from 'vue-router'

// 扩展 RouteMeta：给路由记录附加菜单所需的 title / icon / hidden
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    hidden?: boolean
  }
}

// 常量路由：不需要权限，任何人都可访问（登录页 / 首页 / 大屏 / 404 页）
export const constantRoutes: RouteRecordRaw[] = [
  // 登录页（不挂在 layout 上，菜单隐藏）
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true },
  },
  // 404 页面本体
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404/index.vue'),
    meta: { hidden: true },
  },
  // layout 容器 + 首页
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'home',
        component: () => import('@/views/home/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled' },
      },
      // 个人信息：从右上角头像下拉进入，不在左侧菜单显示
      {
        path: 'profile',
        name: 'profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人信息', hidden: true },
      },
    ],
  },
  // 数据大屏（全屏独立页）
  {
    path: '/screen',
    name: 'screen',
    component: () => import('@/views/screen/index.vue'),
    meta: { title: '数据大屏', icon: 'Histogram' },
  },
]

// 异步路由：需要权限，登录后按用户 routes 动态挂载
// 父路由【不需要 name】（过滤逻辑是「父路由靠子项存活」，不参与名字匹配）
// 子路由的 name【用小写】，如 'user' / 'role' / 'permission' / 'attr' / ...
// ⚠️ 命名约定：后端 GetUserInfo 返回的 routes 是 menu.code，为【大写首字母】形式
//    （User / Role / Permission / Trademark / Attr / Category / Spu / Sku，见 apps/server/init-sql/init.sql）。
//    前端 name 用小写即可——filterAsyncRoute 会统一 toLowerCase() 比较，大小写不一致已被兜住，
//    不需要把前端 name 改成大写，也不要在 mock 里写前端不存在的 name（如 'Home'，前端无此路由）。
export const asyncRoutes: RouteRecordRaw[] = [
  {
    path: '/acl',
    component: () => import('@/layout/index.vue'),
    meta: { title: '权限管理', icon: 'Lock' },
    children: [
      {
        path: 'user',
        name: 'user',
        component: () => import('@/views/acl/user/index.vue'),
        meta: { title: '用户管理', icon: 'UserFilled' },
      },
      {
        path: 'role',
        name: 'role',
        component: () => import('@/views/acl/role/index.vue'),
        meta: { title: '角色管理', icon: 'UserFilled' },
      },
      {
        path: 'permission',
        name: 'permission',
        component: () => import('@/views/acl/permission/index.vue'),
        meta: { title: '菜单管理', icon: 'Menu' },
      },
    ],
  },
  {
    path: '/product',
    component: () => import('@/layout/index.vue'),
    meta: { title: '商品管理', icon: 'Goods' },
    children: [
      {
        path: 'attr',
        name: 'attr',
        component: () => import('@/views/product/attr/index.vue'),
        meta: { title: '属性管理', icon: 'Files' },
      },
      {
        path: 'category',
        name: 'category',
        component: () => import('@/views/product/category/index.vue'),
        meta: { title: '分类管理', icon: 'Menu' },
      },
      {
        path: 'spu',
        name: 'spu',
        component: () => import('@/views/product/spu/index.vue'),
        meta: { title: 'Spu管理', icon: 'Goods' },
      },
      {
        path: 'sku',
        name: 'sku',
        component: () => import('@/views/product/sku/index.vue'),
        meta: { title: 'Sku管理', icon: 'Goods' },
      },
      {
        path: 'trademark',
        name: 'trademark',
        component: () => import('@/views/product/trademark/index.vue'),
        meta: { title: '品牌管理', icon: 'Stamp' },
      },
    ],
  },
]

// 任意路由：必须最后 addRoute，做 404 兜底
export const anyRoute: RouteRecordRaw = {
  path: '/:pathMatch(.*)*',
  redirect: '/404',
  meta: { hidden: true },
}
