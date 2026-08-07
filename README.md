# newproject

前后端同仓（monorepo）练习项目。

- 前端：`apps/frontend`（Vue3 + Vite + TypeScript，工程名 `vue_admin`）
- 后端：`apps/backend`（Go + Gin）

## 仓库布局

```
newproject/
├── .gitignore
├── .editorconfig
├── .npmrc                # 强制 pnpm（engine-strict=true）
├── package.json          # 根启动编排（concurrently 一键启动）
├── pnpm-workspace.yaml   # pnpm 工作区
├── lefthook.yml          # 统一 Git 钩子（pre-commit / commit-msg）
└── apps/
    ├── backend/          # Go + Gin（业务代码自行实现）
    └── frontend/         # Vue3 + Vite + TS（name: vue_admin）
```

## 技术栈与工程化

| 层 | 工具 |
| --- | --- |
| 前端框架 | Vue3 + Vite + TypeScript |
| 前端代码检查 | ESLint（typescript-eslint + eslint-plugin-vue） |
| 样式检查 | Stylelint |
| 代码格式化 | Prettier |
| 类型检查 | vue-tsc（纳入 `build`） |
| 提交信息规范 | commitlint（Conventional Commits） |
| 后端代码检查 | golangci-lint + `go vet ./...` |
| 统一 Git 钩子 | lefthook（Go 实现，前后端一起管） |
| 包管理 | 强制 pnpm（禁止 npm） |

## 环境要求

- Node.js（建议 LTS）
- [pnpm](https://pnpm.io/)（`npm i -g pnpm` 或 `corepack enable`）
- Go（建议 1.21+）

## 快速开始

```bash
# 1. 安装依赖（根目录，会安装前后端所需）
pnpm install

# 2. 一键启动前后端（前端 5173，后端 8080）
pnpm dev
```

前端开发服务器：http://localhost:5173
后端接口示例：http://localhost:8080/api/health

## 常用脚本

| 命令 | 说明 |
| --- | --- |
| `pnpm dev` | 同时启动前端与后端 |
| `pnpm -C apps/frontend lint` | 前端 ESLint 检查 |
| `pnpm -C apps/frontend stylelint` | 前端 Stylelint 检查 |
| `pnpm -C apps/frontend format` | 前端 Prettier 格式化 |
| `pnpm -C apps/frontend build` | 前端类型检查 + 构建 |
| `pnpm -C apps/backend exec golangci-lint run` | 后端 lint（或进入目录执行） |
| `cd apps/backend && go vet ./...` | 后端静态检查 |

## Git 提交规范

提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/)，
由 `lefthook` 的 `commit-msg` 钩子通过 `commitlint` 校验。

示例：

```
feat: 新增用户登录接口
chore: 初始化前端工程
fix: 修复 CORS 跨域问题
```

`pre-commit` 钩子会自动并行执行：前端 lint + 后端 golangci-lint + `go vet ./...`。

## 分工说明

- **配置 / 工程化**：本仓库的脚手架与各 lint / 格式化 / 钩子配置文件。
- **业务代码**：Go 路由、中间件、CORS，以及 Vue 组件、页面逻辑、axios 调用等由使用者自行实现。
