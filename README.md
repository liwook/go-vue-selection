# go-vue-selection

前后端同仓（monorepo）的**中后台脚手架**，内置 Vue3 管理后台与数据大屏（`/screen`）。

- 前端：`apps/web`（Vue3 + Vite + TypeScript，工程名 `vue_admin`）
- 后端：`apps/server`（Go + Gin），数据库：PostgreSQL

## 仓库布局

```text
go-vue-selection/
├── .npmrc                # 强制 pnpm（engine-strict=true）
├── .gitattributes        # 统一换行符为 LF（配合 Biome lineEnding: lf）
├── package.json          # 根编排：pnpm workspace + 跨端脚本（prepare / dev / build / lint 转发到 apps/web）
├── pnpm-workspace.yaml   # pnpm 工作区（apps/*）
├── lefthook.yml          # 统一 Git 钩子（pre-commit / commit-msg）
├── commitlint.config.mjs # 提交信息校验（Conventional Commits）
└── apps/
    ├── server/           # Go + Gin（handler / service / repository 分层）
    └── web/              # Vue3 + Vite + TS（name: vue_admin）
```

## 技术栈

- 前端：Vue3 + Vite + TypeScript，UI 用 Element Plus，请求封装用 `openapi-fetch`（类型来自 Swagger 生成的 `src/api/schema.d.ts`，clone 即跑）
- 后端：Go + Gin，接口经 Swagger 生成 OpenAPI schema，配置走 `apps/server/etc/config.yaml`（可用 `VUE_ADMIN_*` 环境变量覆盖）
- 数据：PostgreSQL（容器化自带 18.3），前端构建为纯静态产物，由 Nginx / Docker 托管

## 环境要求

- Node.js ≥ 20（建议 LTS）
- [pnpm](https://pnpm.io/)（`corepack enable` 或 `npm i -g pnpm`；仓库锁定 `pnpm@11.20.0`）
- Go 1.26+
- Docker + Docker Compose（仅容器化启动时需要）
- PostgreSQL：容器化方式由 docker-compose 中的 pg 服务（postgres:18.3-alpine）提供
- 两种启动方式均依赖 PostgreSQL：方式一由 docker-compose 的 pg 服务提供；方式二可用自备实例或 docker 启动的 PG

## 快速开始

### 方式一：容器化一键全栈（推荐，一条命令同时启动前端 + 后端 + 数据库）

在**仓库根目录**执行：

```bash
# 无需 .env：pg 密码与 jwt_secret 均已在 compose / config.yaml.example 内置默认值
# 1. 一条命令拉起 PostgreSQL + 后端(Go) + 前端(Nginx)
docker compose up -d --build
```

- 访问站点：<http://localhost>
- 对外仅暴露 80（前端 Nginx），`/api` 由 Nginx 反代到后端 `server:9000`，数据库端口不对外暴露
- 容器内使用 `etc/config.yaml.example` 作默认配置，postgres 密码与 jwt_secret 均内置默认（`postgres` / 演示密钥），**无需 `.env`、无需本地私有 `config.yaml`**；如需自定义，可用 `VUE_ADMIN_*` 环境变量覆盖（如 `VUE_ADMIN_POSTGRES_PASSWORD`、`VUE_ADMIN_AUTH_JWT_SECRET`）

### 方式二：本地开发（前后端分离，热更新）

```bash
# 1. 安装依赖（根目录）
pnpm install

# 2. 准备前端环境变量
cp apps/web/.env.example apps/web/.env.development
#   然后将 .env.development 中的 API_PROXY_TARGET 改成你本地/可达的后端地址（默认 http://localhost:9000）

# 3. 准备后端本地配置
cp apps/server/etc/config.yaml.example apps/server/etc/config.yaml
#   - 若 PG 是远程的，请把 config.yaml 的 postgres.host 改成对应的远程 IP（密码也一并对应修改）

# 4. 启动 PostgreSQL 容器（后端依赖 PG，需先起；无需 .env，密码有默认值 postgres）
cd apps/server && docker compose up -d pg
#   - 容器内 5432 映射到宿主机 127.0.0.1:5432，正好对上后端默认配置

# 5. 后端本地运行（需已装 Go，沿用上一步已 cd 的 apps/server 目录）
#    - 先生成 Swagger 文档文件（swag 命令可全局安装：go install github.com/swaggo/swag/cmd/swag@latest）
swag init -g main.go -o api --parseInternal
go run .

# 6. 另开一个新终端启动前端开发服务器（5173），其 /api 代理到后端 :9000，故建议后端就绪后再起
pnpm dev
```

- 前端开发服务器：<http://localhost:5173>
- 后端健康检查：<http://localhost:9000/health>
- 后端 Swagger：执行 `cd apps/server && make swagger` 生成到 `apps/server/api/`，访问 <http://localhost:9000/swagger/index.html>

> 根目录 `pnpm dev` 仅启动前端；后端需单独运行。容器化部署时由 `docker-compose.yaml` 的环境变量（前缀 `VUE_ADMIN_`，如 `VUE_ADMIN_POSTGRES_HOST=pg`）覆盖数据库与日志，日志统一输出 stdout，首次启动自动执行 `apps/server/init-sql/` 下建表脚本。

## 生产部署

前端构建为纯静态产物，打包与 Nginx / Docker 部署见 [`apps/web/docs/13.阶段16·打包与上线部署.md`](./apps/web/docs/13.阶段16·打包与上线部署.md)。

## 文档索引

分阶段的实现指南、后端分层规范等详细文档见仓库 [`apps/web/docs/`](./apps/web/docs/) 与 [`apps/server/docs/`](./apps/server/docs/) 目录；工程化配置（pnpm workspace / Biome / lefthook / commitlint）见 [`工程化配置指南.md`](./工程化配置指南.md)。
