// 项目级配置（标题、logo 等），供 Layout 各组件读取
const settings = {
  // 浏览器标签 / 菜单顶部标题
  title: '甄选后台管理系统',
  // 走 public/ 静态资源，vite 构建时会原样拷贝到 dist 根，配合 BASE_URL 兼容子路径部署
  logo: `${import.meta.env.BASE_URL}chopper.png`,
  // 是否在菜单顶部展示 logo + 标题
  logoShow: true,
}

export default settings
