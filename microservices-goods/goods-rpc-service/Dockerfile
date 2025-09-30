# 多阶段构建
FROM golang:1.22-alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 复制 proto 生成物（已在仓库内）
COPY proto ./proto

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -a -installsuffix cgo -o goods-rpc .

# 运行阶段
FROM alpine:latest

# 安装ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/goods-rpc .
COPY --from=builder /app/etc ./etc

# 暴露端口
EXPOSE 9001

# 运行应用
CMD ["./goods-rpc", "-f", "etc/goods-rpc.yaml"]



