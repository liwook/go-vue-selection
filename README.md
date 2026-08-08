# newproject

前后端同仓（monorepo）练习项目。

- 前端：`apps/web`（Vue3 + Vite + TypeScript，工程名 `vue_admin`）
- 后端：`apps/server`（Go + Gin）

## 仓库布局

```
newproject/
├── .gitignore
├── .npmrc                # 强制 pnpm（engine-strict=true）
├── .gitattributes        # 统一换行符为 LF（配合 Biome lineEnding: lf）
├── package.json          # 根编排：pnpm workspace + 跨端脚本（prepare / dev / build / lint 转发到 apps/web）
├── pnpm-workspace.yaml   # pnpm 工作区（apps/*）
├── lefthook.yml          # 统一 Git 钩子（pre-commit / commit-msg）
├── commitlint.config.mjs # 提交信息校验（Conventional Commits）
└── apps/
    ├── server/           # Go + Gin（业务代码自行实现）
    └── web/              # Vue3 + Vite + TS（name: vue_admin）

## 技术栈与工程化

| 层 | 工具 |
| --- | --- |
| 前端框架 | Vue3 + Vite + TypeScript |
| 前端 lint + 格式化 | Biome（`.js/.ts/.vue(<script>)`，主力） |
| 前端模板检查 | ESLint（仅 `.vue` 的 `<template>`） |
| 样式检查 | 暂未启用（Stylelint 待 CSS 规范需求时接入） |
| 类型检查 | vue-tsc（纳入 `build`） |
| 提交信息规范 | commitlint（Conventional Commits） |
| 后端代码检查 | golangci-lint（`.golangci.yml`）+ `gofmt` + `go vet ./...` |
| 统一 Git 钩子 | lefthook（Go 实现，前后端一起管） |
| 包管理 | 强制 pnpm（禁止 npm） |

> 说明：前端格式化与 lint 由 **Biome** 统一承担，**不使用 Prettier**；ESLint 仅补 `.vue` 模板规则，不引入 `typescript-eslint` 解析器（避免与 Biome 在 `.ts` 规则上重叠）。详见 `工程化配置指南.md`。

## 环境要求

- Node.js（建议 LTS）
- [pnpm](https://pnpm.io/)（`npm i -g pnpm` 或 `corepack enable`）
- Go（建议 1.21+）

## 快速开始

```bash
# 1. 安装前端依赖（根目录）
pnpm install

# 2. 启动前端开发服务器（5173）
pnpm dev

# 3. 后端启动方式（二选一）
# 方式 A：本地直接运行（需已装 Go，且存在 ./apps/server/etc/config.yaml）
cd apps/server && go run ./cmd/server

# 方式 B：容器化启动（PostgreSQL + 后端，依赖 docker-compose）
cd apps/server && cp .env.example .env   # 按需修改
docker compose up -d
```

前端开发服务器：http://localhost:5173
后端接口示例：http://localhost:8080/api/health

> 说明：根目录的 `pnpm dev` 仅启动前端；后端需单独运行。后端默认通过 `apps/server/etc/config.yaml` 读取配置，容器化部署时由 `docker-compose.yaml` 中的环境变量覆盖数据库连接等信息，日志统一输出到 stdout。

## 常用脚本

| 命令 | 说明 |
| --- | --- |
| `pnpm dev` | 启动前端开发服务器（仅前端） |
| `pnpm --filter web lint` | 前端 Biome 检查（lint + 格式化） |
| `pnpm --filter web lint:vue` | 前端 ESLint 仅查 `.vue` 模板 |
| `pnpm --filter web format` | 前端 Biome 格式化（`biome check --write`） |
| `pnpm --filter web build` | 前端类型检查（vue-tsc）+ 构建 |
| `cd apps/server && golangci-lint run ./...` | 后端聚合 lint |
| `cd apps/server && go vet ./...` | 后端静态检查 |
| `cd apps/server && go test -tags=integration ./tests/...` | 后端集成测试（需连 `vue_admin_test` 库） |
| `cd apps/server && go test -tags=integration -race ./tests/...` | 集成测试 + 竞态检测（CI 定时/发布前跑，见下） |

### 测试与竞态检测（`-race`）

集成测试通过进程内 `httptest.Server` 串行发起 HTTP 请求，用例之间无并发、包级 `apiClient` 在 `TestMain` 初始化后只读，因此**日常本地回归不必带 `-race`**：

```bash
cd apps/server && go test -tags=integration ./tests/...
```

`-race` 的价值在于"体检"被测的服务器代码（handler / middleware / 全局组件）在并发请求下是否隐藏数据竞争，建议放到 CI 定时任务或发布前检查中跑一遍（成本仅慢几倍）：

```bash
cd apps/server && go test -tags=integration -race ./tests/...
```

前提：必须在 `apps/server` 目录下，且能连上 `vue_admin_test` 库、存在 `./etc/config.yaml`（由 `TestMain` 加载）。若后续编写并发/压测类用例（用了 `t.Parallel()` 或多个 goroutine 复用 `apiClient`），则必须默认带 `-race`。

## Git 提交规范

提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/)，
由 `lefthook` 的 `commit-msg` 钩子通过 `commitlint` 校验。

示例：

```
feat: 新增用户登录接口
chore: 初始化前端工程
fix: 修复 CORS 跨域问题
```

`pre-commit` 钩子会自动并行执行：前端 Biome 检查 + ESLint（仅 `.vue` 模板）+ 后端 `gofmt` 格式化 + `go vet ./...`（golangci-lint 不在钩子里跑，需手动在 `apps/server` 执行）。

## 分工说明

- **配置 / 工程化**：本仓库的脚手架与各 lint / 格式化 / 钩子配置文件。
- **业务代码**：Go 路由、中间件、CORS，以及 Vue 组件、页面逻辑、axios 调用等由使用者自行实现。
