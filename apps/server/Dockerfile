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
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 编译成二进制可执行文件 app（swagger 文档已在本地 swag init 生成并提交，无需在容器内生成）
RUN go build -o app .


# 接下来创建一个轻量运行镜像
FROM alpine:3.22

# 安装 CA 证书（用于 HTTPS/TLS 出站请求）和健康检查所需的 wget
RUN apk add --no-cache ca-certificates wget \
    # 创建非 root 用户 appuser（uid/gid 10001），运行时以该用户身份执行
    && addgroup -S appuser \
    && adduser -S -D -H -s /sbin/nologin -G appuser -u 10001 appuser

# 从builder镜像中把静态文件和配置文件拷贝到根目录
COPY --from=builder /build/static /static
COPY --from=builder /build/etc /etc

# 从builder镜像中把 app 拷贝到根目录
COPY --from=builder /build/app /app

# 将所需文件权属交给非 root 用户（app 二进制、配置、静态资源）
RUN chown -R appuser:appuser /app /etc /static

# 声明服务端口
EXPOSE 9000

# 切换到非 root 用户运行
USER appuser

# 需要运行的命令（使用绝对路径）
ENTRYPOINT ["/app", "-f", "/etc/config.yaml"]