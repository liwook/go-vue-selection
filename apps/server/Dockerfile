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


# 接下来创建一个小镜像
FROM scratch

# 从builder镜像中把静态文件和配置文件拷贝到根目录
COPY --from=builder /build/static /static
COPY --from=builder /build/etc /etc

# 从builder镜像中把 app 拷贝到根目录
COPY --from=builder /build/app /app

# 声明服务端口
EXPOSE 9000

# 需要运行的命令（使用绝对路径，避免 scratch 镜像中工作目录不确定导致的找不到配置文件问题）
ENTRYPOINT ["/app", "-f", "/etc/config.yaml"]