#!/bin/bash

# 库存微服务启动脚本
# 使用方法：bash start.sh

set -e

echo "正在启动库存微服务..."

# 检查 Docker 和 docker-compose 是否已安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: docker-compose 未安装。请先安装 docker-compose"
    exit 1
fi

# 生成 proto 代码（如果需要）
if [ ! -f "proto/stock/stock.pb.go" ] || [ ! -f "proto/stock/stock_grpc.pb.go" ]; then
    echo "生成 proto 代码..."
    bash scripts/gen-proto.sh
fi

# 停止并删除现有容器
echo "停止现有服务..."
docker-compose down

# 构建并启动服务
echo "构建并启动服务..."
docker-compose up --build -d

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

echo ""
echo "✅ 库存微服务启动完成！"
echo ""
echo "🚀 API 网关地址: http://localhost:8080"
echo ""
echo "📖 API 文档:"
echo "  - 设置库存: POST http://localhost:8080/api/v1/stock/set"
echo "  - 获取库存: GET http://localhost:8080/api/v1/stock/get?goods_id=1"
echo "  - 扣减库存: POST http://localhost:8080/api/v1/stock/reduce"
echo "  - 回滚库存: POST http://localhost:8080/api/v1/stock/rollback"
echo "  - 批量获取: POST http://localhost:8080/api/v1/stock/batch/get"
echo "  - 批量扣减: POST http://localhost:8080/api/v1/stock/batch/reduce"
echo ""
echo "🔧 查看日志: docker-compose logs -f [service-name]"
echo "🛑 停止服务: docker-compose down"

