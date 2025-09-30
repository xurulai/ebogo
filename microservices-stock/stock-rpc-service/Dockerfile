FROM golang:1.22-alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 复制 proto 生成物（已在仓库内）
COPY proto ./proto

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stock-rpc stock-rpc.go

FROM alpine:latest

# 安装 ca-certificates，用于 HTTPS 请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /build/stock-rpc .
COPY --from=builder /build/etc ./etc

# 暴露端口
EXPOSE 8889

# 运行应用
CMD ["./stock-rpc"]

