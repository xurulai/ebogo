#!/bin/bash

# 微服务架构统一启动脚本

echo "=== 微服务商品系统启动 ==="
echo "当前目录: $(pwd)"
echo "时间: $(date)"
echo ""

# 检查必要工具
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装"
    exit 1
fi

# 构建服务
echo "🔨 构建微服务..."

# 构建商品服务
echo "  - 构建商品服务..."
cd goods-service
if GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o goods-service .; then
    echo "    ✅ 商品服务构建成功"
else
    echo "    ❌ 商品服务构建失败"
    exit 1
fi
cd ..

# 构建网关服务
echo "  - 构建网关服务..."
cd api-gateway
if GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway-api .; then
    echo "    ✅ 网关服务构建成功"
else
    echo "    ❌ 网关服务构建失败"
    exit 1
fi
cd ..

# 停止现有服务
echo "🧹 停止现有服务..."
docker-compose down --volumes --remove-orphans

# 启动微服务架构
echo "🚀 启动微服务架构..."
docker-compose up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
echo "  - MySQL 启动中..."
sleep 20
echo "  - Redis 启动中..."
sleep 10
echo "  - 商品服务启动中..."
sleep 15
echo "  - 网关服务启动中..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

echo ""
echo "📝 查看服务日志..."
docker-compose logs --tail=5 api-gateway

echo ""
echo "🧪 测试微服务API..."
echo ""

# 测试网关API
echo "测试网关API（对外唯一入口）:"
if curl -s "http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1" > /dev/null; then
    echo "✅ 网关API测试成功"
    echo "响应数据："
    curl -s "http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1" | head -3
else
    echo "❌ 网关API测试失败"
    echo "查看详细日志："
    docker-compose logs api-gateway
fi

echo ""
echo "=== 微服务架构启动完成 ==="
echo ""
echo "🌐 服务架构："
echo "  外部访问："
echo "    - API网关:    http://localhost:8080 （统一入口）"
echo "  内部服务："
echo "    - 商品服务:    goods-service:8888 （内部网络）"
echo "    - MySQL:      mysql:3306"
echo "    - Redis:      redis:6379"
echo ""
echo "🧪 测试命令："
echo "  # 通过网关访问（推荐）"
echo "  curl \"http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1\""
echo "  curl \"http://localhost:8080/api/v1/goods/detail?goods_id=1&user_id=1\""
echo ""
echo "🔧 管理命令："
echo "  - 查看状态: docker-compose ps"
echo "  - 查看日志: docker-compose logs -f [service]"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart [service]"
echo ""
echo "📋 架构特点："
echo "  ✅ 统一网关入口 (8080)"
echo "  ✅ 内部服务隔离"
echo "  ✅ 服务发现和负载均衡"
echo "  ✅ 统一认证和监控"
echo "  ✅ 容器化部署"
