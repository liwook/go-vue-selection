# apps/server Makefile
# 管理 Go 后端常用开发动作，与前端（web）解耦，不依赖 Node/pnpm。

GO ?= go
PKG := ./...

.PHONY: test test-unit test-it vet build fmt tidy swagger

## swagger: 根据代码注解生成 OpenAPI 文档到 api/ 目录
swagger:
	$(GO) run github.com/swaggo/swag/cmd/swag init -g main.go -o api --parseDependency --parseInternal

## test: 运行全部单元测试（默认）
test: test-unit

## test-unit: 运行单元测试（不含集成测试 build tag）
test-unit:
	$(GO) test $(PKG)

## test-it: 运行集成测试（需要 -tags=integration，且 tests/ 目录存在）
test-it:
	$(GO) test -tags=integration -v ./tests/...

## vet: 静态检查
vet:
	$(GO) vet $(PKG)

## build: 编译服务
build:
	$(GO) build -o bin/server .

## fmt: 格式化
fmt:
	$(GO) fmt $(PKG)

## tidy: 整理依赖
tidy:
	$(GO) mod tidy
