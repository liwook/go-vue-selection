// 扩展 RouteMeta：给路由记录附加菜单所需的 title / icon / hidden
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    hidden?: boolean
  }
}

// 对外暴露常量路由（本阶段为静态路由，后续权限阶段会在此基础上筛选）
export const constantRoute = [
  // 登录页（不挂 Layout，隐藏于菜单）
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true },
  },

  // Layout 容器 + 首页
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        component: () => import('@/views/home/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled' },
      },
    ],
  },

  // 数据大屏（全屏独立页，不挂 Layout，单层即可）
  {
    path: '/screen',
    component: () => import('@/views/screen/index.vue'),
    meta: { title: '数据大屏', icon: 'Histogram' },
  },

  // 权限管理 acl
  {
    path: '/acl',
    component: () => import('@/layout/index.vue'),
    redirect: '/acl/user',
    meta: { title: '权限管理', icon: 'Lock' },
    children: [
      {
        path: 'user',
        component: () => import('@/views/acl/user/index.vue'),
        meta: { title: '用户管理', icon: 'UserFilled' },
      },
      {
        path: 'role',
        component: () => import('@/views/acl/role/index.vue'),
        meta: { title: '角色管理', icon: 'UserFilled' },
      },
      {
        path: 'permission',
        component: () => import('@/views/acl/permission/index.vue'),
        meta: { title: '菜单管理', icon: 'Lock' },
      },
    ],
  },

  // 商品管理 product
  {
    path: '/product',
    component: () => import('@/layout/index.vue'),
    redirect: '/product/trademark',
    meta: { title: '商品管理', icon: 'Goods' },
    children: [
      {
        path: 'trademark',
        component: () => import('@/views/product/trademark/index.vue'),
        meta: { title: '品牌管理', icon: 'Stamp' },
      },
      {
        path: 'attr',
        component: () => import('@/views/product/attr/index.vue'),
        meta: { title: '属性管理', icon: 'Files' },
      },
      {
        path: 'category',
        component: () => import('@/views/product/category/index.vue'),
        meta: { title: '分类管理', icon: 'Menu' },
      },
      {
        path: 'spu',
        component: () => import('@/views/product/spu/index.vue'),
        meta: { title: 'Spu管理', icon: 'Goods' },
      },
      {
        path: 'sku',
        component: () => import('@/views/product/sku/index.vue'),
        meta: { title: 'Sku管理', icon: 'Goods' },
      },
    ],
  },

  // 404
  {
    path: '/404',
    component: () => import('@/views/404/index.vue'),
    meta: { title: '404', hidden: true },
  },

  // 任意未匹配路径 -> 重定向到 404（hidden 避免菜单渲染出空白项）
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    meta: { hidden: true },
  },
]
