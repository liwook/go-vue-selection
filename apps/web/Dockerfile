# syntax=docker/dockerfile:1.4
# 启用 BuildKit 1.4 前端解析器，支持 RUN --mount=type=cache 等语法（须放第一行）
# 构建上下文必须为仓库根目录（docker-compose 中 context: . + dockerfile: apps/web/Dockerfile）
# 原因：pnpm workspace 需从根目录解析 apps/*，且 COPY 不能超出 build context

# ---- 构建参数（构建时可在外部覆盖版本）----
# 版本对照：根 package.json 的 packageManager=pnpm@11.20.0、engines.node>=20
# node 用 24（LTS，满足 >=20）；nginx 用 1.30.4 稳定线
ARG NODE_VERSION=24.18.0
ARG NGINX_VERSION=1.30.4

# ---- 阶段 1：构建 ----
FROM node:${NODE_VERSION}-alpine AS builder
WORKDIR /build

# 启用并锁定 pnpm（与根 package.json 的 packageManager 一致）
RUN corepack enable && corepack prepare pnpm@11.20.0 --activate

# 依赖下载走国内镜像（境内构建加速；境外机器可删此行）
RUN pnpm config set registry https://registry.npmmirror.com

# 依赖层：只拷依赖清单（变化极少，会被缓存）
# 注意：pnpm-workspace 根目录结构，需同时拷根 package.json / pnpm-lock.yaml / pnpm-workspace.yaml
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY apps/web/package.json apps/web/package.json

# 安装依赖（命中缓存则跳过，避免每次改源码都重装）
# 注意：不能加 --ignore-scripts。vite 依赖 esbuild，其 postinstall 负责放置平台二进制，
# 跳过会导致 vite build 因缺失二进制而失败；根 prepare 的 lefthook install 已容错(|| true)，无 .git 也不致命。
# --mount=type=cache：BuildKit 缓存挂载，依赖落盘到宿主机持久区而非镜像层，跨构建复用加速。
RUN --mount=type=cache,id=vue_admin_pnpm-store,target=/root/.pnpm-store \
    --mount=type=cache,id=vue_admin_pnpm-build,target=/root/.cache/pnpm \
    pnpm install --frozen-lockfile

# 源码层：经常变化，但不影响上面的依赖缓存
# .dockerignore 已排除 apps/server/ 等无关内容，缩小上下文并保住依赖缓存层
COPY . .

# 构建（生产模式）：在 web 子包下执行 build:prod，产物输出到 apps/web/dist
RUN pnpm --filter web build:prod

# ---- 阶段 2：运行 ----
FROM nginx:${NGINX_VERSION}-alpine

# 复制构建产物（pnpm workspace 下 web 产物的实际路径为 apps/web/dist）
COPY --from=builder /build/apps/web/dist /usr/share/nginx/html

# 复制 Nginx 配置（含 /api 反代到 server:9000，server 为 compose service 名非容器名）
# 后端地址在 compose 网络内固定为 server:9000，故直接进镜像，无需部署时挂载注入
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
