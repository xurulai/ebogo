#!/bin/bash

# 生成 proto 代码脚本
# 使用方法：bash scripts/gen-proto.sh

set -e

echo "开始生成 proto 代码..."

# 检查必要工具是否已安装
if ! command -v protoc &> /dev/null; then
    echo "错误: protoc 未安装。请先安装 Protocol Buffers 编译器"
    exit 1
fi

if ! command -v protoc-gen-go &> /dev/null; then
    echo "错误: protoc-gen-go 未安装。请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5"
    exit 1
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "错误: protoc-gen-go-grpc 未安装。请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1"
    exit 1
fi

# 进入 proto 目录
cd proto

# 生成 Go 代码
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       stock/stock.proto

echo "proto 代码生成完成！"
echo "生成的文件："
echo "  - stock/stock.pb.go (消息定义)"
echo "  - stock/stock_grpc.pb.go (gRPC 服务定义)"




