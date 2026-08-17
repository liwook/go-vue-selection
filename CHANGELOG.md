# Changelog

本项目所有重要变更记录于此文件。格式参考 [Keep a Changelog](https://keepachangelog.com/)，版本号遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)。

## [v1.0.0] - 2026-08-17

首个正式版本。前后端功能完整、可独立运行，monorepo 工程化体系就绪。

### Added
- 前端：`apps/web`（Vue3 + Vite + TypeScript + Element Plus），含管理后台与数据大屏（`/screen`），路由守卫 + 动态权限路由。
- 后端：`apps/server`（Go + Gin + PostgreSQL），handler / service / repository 分层，Swagger 生成 OpenAPI schema。
- 鉴权：JWT + opaque token，配套 RBAC / ACL 权限模型。
- 工程化：pnpm workspace 多包编排、commitlint（Conventional Commits）、lefthook（pre-commit / commit-msg 钩子）。
- 容器化：docker-compose 一键拉起 PostgreSQL + 后端 + 前端 Nginx，clone 即跑。
- 配套文档：架构总览、路由守卫与动态权限路由说明、工程化配置指南。
