# go-vue-selection

前后端同仓（monorepo）的**中后台脚手架**，内置 Vue3 管理后台与数据大屏（`/screen`）。

- 前端：`apps/web`（Vue3 + Vite + TypeScript，工程名 `vue_admin`）
- 后端：`apps/server`（Go + Gin）

## 文档索引

- **工程化**：[`工程化配置指南.md`](./工程化配置指南.md)——pnpm workspace / Biome / lefthook / commitlint 配置详解
- **前端**：
  - [阶段02·前端工程搭建与选型总结](./apps/web/docs/02.前端工程搭建与选型总结.md)——脚手架、Element Plus、openapi-fetch 请求封装
  - [阶段15·数据大屏实现指南](./apps/web/docs/12.阶段15·数据大屏实现指南.md)——ECharts + 中国地图
  - [阶段16·打包与上线部署](./apps/web/docs/13.阶段16·打包与上线部署.md)——Vite 构建与 Nginx / Docker 部署
- **后端**：[`apps/server/docs/`](./apps/server/docs/)——后端设计、测试规范（含集成测试 `-race` 说明）

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
    ├── server/           # Go + Gin（业务代码自行实现）
    └── web/              # Vue3 + Vite + TS（name: vue_admin）
```

## 技术栈

- 前端：Vue3 + Vite + TypeScript，UI 用 Element Plus，请求封装用 `openapi-fetch`（类型来自 Swagger 生成的 `src/api/schema.d.ts`，clone 即跑）
- 后端：Go + Gin，接口经 Swagger 生成 OpenAPI schema，配置走 `apps/server/etc/config.yaml`（可用 `VUE_ADMIN_*` 环境变量覆盖）
- 数据：PostgreSQL（容器化自带 18.3），前端构建为纯静态产物，由 Nginx / Docker 托管

> 工程化（pnpm workspace / Biome / lefthook / commitlint 等）配置详解见 `工程化配置指南.md`。

## 环境要求

- Node.js ≥ 20（建议 LTS）
- [pnpm](https://pnpm.io/)（`corepack enable` 或 `npm i -g pnpm`；仓库锁定 `pnpm@11.20.0`）
- Go 1.26+
- Docker + Docker Compose（仅容器化启动后端时需要；自带 PostgreSQL 18.3）

## 快速开始

### 方式一：容器化一键全栈（推荐，一条命令同时启动前端 + 后端 + 数据库）

在**仓库根目录**执行：

```bash
# 1. 准备环境变量（填入 POSTGRES_PASSWORD 与 AUTH_JWT_SECRET）
cp apps/server/.env.example apps/server/.env

# 2. 一条命令拉起 PostgreSQL + 后端(Go) + 前端(Nginx)
docker compose up -d --build
```

- 访问站点：<http://localhost>
- 对外仅暴露 80（前端 Nginx），`/api` 由 Nginx 反代到后端 `server:9000`，数据库端口不对外暴露
- 容器内使用 `etc/config.yaml.example` 作默认配置，postgres 连接 / 日志输出 / jwt_secret 等由 `.env` 注入，**无需本地私有 `config.yaml`**

### 方式二：本地开发（前后端分离，热更新）

```bash
# 1. 安装依赖（根目录）
pnpm install

# 2. 启动前端开发服务器（5173）
pnpm dev

# 3. 后端本地运行（需已装 Go，且存在 apps/server/etc/config.yaml）
#   - config.yaml 中 postgres.host 默认指向远端，请改成你本地/可达的 PG 地址
#   - 启动后监听 :9000，支持 -f 指定配置文件路径
cd apps/server && go run .
```

- 前端开发服务器：<http://localhost:5173>
- 后端健康检查：<http://localhost:9000/health>
- 后端 Swagger：执行 `cd apps/server && make swagger` 生成到 `api/`，访问 <http://localhost:9000/swagger/index.html>

> 根目录 `pnpm dev` 仅启动前端；后端需单独运行。容器化部署时由 `docker-compose.yaml` 的环境变量（前缀 `VUE_ADMIN_`，如 `VUE_ADMIN_POSTGRES_HOST=pg`）覆盖数据库与日志，日志统一输出 stdout，首次启动自动执行 `apps/server/init-sql/` 下建表脚本。

后端设计与测试（含集成测试 `-race` 竞态检测）、`make` 命令与 `config.yaml` 字段说明见 `apps/server/docs/`。

## 生产部署

前端构建为纯静态产物，打包与 Nginx / Docker 部署见 [`apps/web/docs/13.阶段16·打包与上线部署.md`](./apps/web/docs/13.阶段16·打包与上线部署.md)。
