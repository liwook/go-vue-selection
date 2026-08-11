FROM golang:1.26-alpine AS builder
LABEL authors="li"

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum 文件并下载依赖信息
# 使用 BuildKit cache mount 让 mod cache 跨构建复用，避免每次重新下载
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 将代码复制到容器中
COPY . .

# api/docs.go 被 .gitignore 忽略、不进 git，需在容器内用 swag 现生成
# --parseDependency 解析外部包（如 types/）的结构体字段
# --parseInternal 扫描 internal/ 目录
RUN --mount=type=cache,target=/go/pkg/mod \
    go run github.com/swaggo/swag/cmd/swag init -g main.go -o api --parseDependency --parseInternal

# 编译成二进制可执行文件 app
# 复用 mod cache 与 go build cache，大幅加速二次构建
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o app .


# 接下来创建一个轻量运行镜像
FROM alpine:3.22

# 将 Alpine 软件源替换为国内镜像（清华源），避免官方源在境内访问超时
RUN sed -i 's|https://dl-cdn.alpinelinux.org|https://mirrors.tuna.tsinghua.edu.cn|g' /etc/apk/repositories

# 安装 CA 证书（用于 HTTPS/TLS 出站请求）和健康检查所需的 wget
RUN apk add --no-cache ca-certificates wget \
    # 创建非 root 用户 appuser（uid/gid 10001），运行时以该用户身份执行
    && addgroup -S appuser \
    && adduser -S -D -H -s /sbin/nologin -G appuser -u 10001 appuser

# 从builder镜像中把静态文件和配置拷贝到根目录
# 只复制应用自身配置文件，避免覆盖/搬动系统 /etc 下的只读文件（hosts/resolv.conf）
COPY --from=builder /build/static /static
COPY --from=builder /build/etc/config.yaml /etc/config.yaml

# 从builder镜像中把 app 拷贝到根目录
COPY --from=builder /build/app /app

# 将所需文件权属交给非 root 用户（app 二进制、静态资源、应用配置）
# 注意：不要 chown 整个 /etc，其中 hosts/resolv.conf 是运行时注入的只读文件
RUN chown -R appuser:appuser /app /static \
    && chown appuser:appuser /etc/config.yaml

# 声明服务端口
EXPOSE 9000

# 切换到非 root 用户运行
USER appuser

# 需要运行的命令（使用绝对路径）
ENTRYPOINT ["/app", "-f", "/etc/config.yaml"]