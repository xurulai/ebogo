# 使用官方 Go 镜像作为构建环境
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制 proto 模块（需要先复制，因为 go.mod 依赖它）
COPY proto /proto

# 复制 go.mod 和 go.sum 文件
COPY order-rpc-service/go.mod order-rpc-service/go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY order-rpc-service .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-rpc order-rpc-simple.go

# 使用 alpine 作为运行环境
FROM alpine:latest

# 安装 ca-certificates 和 netcat（用于健康检查）
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/order-rpc .

# 复制配置文件
COPY --from=builder /app/etc ./etc

# 暴露端口
EXPOSE 9002

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD nc -z localhost 9002 || exit 1

# 运行应用
CMD ["./order-rpc"]

