#!/bin/bash

# 订单微服务启动脚本
# 使用方法：bash start.sh

set -e

echo "🚀 启动订单微服务..."

# 检查 Docker 和 Docker Compose 是否已安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker 未安装。请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: Docker Compose 未安装。请先安装 Docker Compose"
    exit 1
fi

# 停止并清理之前的容器
echo "🧹 清理之前的容器..."
docker-compose down --remove-orphans

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose up -d --build

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

# 等待所有服务健康
echo "🏥 等待服务健康检查..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if docker-compose ps | grep -q "healthy"; then
        echo "✅ 服务启动成功！"
        break
    fi
    
    attempt=$((attempt + 1))
    echo "等待中... ($attempt/$max_attempts)"
    sleep 5
done

if [ $attempt -eq $max_attempts ]; then
    echo "⚠️  警告: 服务可能未完全启动，请检查日志"
    docker-compose logs --tail=50
fi

echo ""
echo "🎉 订单微服务启动完成！"
echo ""
echo "📋 服务信息："
echo "  - API Gateway: http://localhost:8080"
echo "  - MySQL: localhost:3306"
echo "  - Redis: localhost:6379"
echo "  - RocketMQ NameServer: localhost:9876"
echo ""
echo "🔧 常用命令："
echo "  - 查看日志: docker-compose logs -f"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart"
echo ""
echo "🧪 测试 API："
echo "  - 创建订单: curl -X POST http://localhost:8080/api/v1/order/create -H 'Content-Type: application/json' -d '{\"goods_id\":1,\"num\":1,\"user_id\":1,\"address\":\"测试地址\",\"name\":\"测试用户\",\"phone\":\"13800138000\"}'"
echo "  - 获取订单列表: curl 'http://localhost:8080/api/v1/order/list?user_id=1&page_num=1&page_size=10'"




