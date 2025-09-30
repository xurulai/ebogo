#!/bin/bash

# 统一启动三个微服务的脚本
# 商品服务、库存服务、订单服务

set -e

echo "🚀 启动所有微服务系统..."
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

# 创建共享 external 网络
echo "🔧 确保共享网络存在..."
docker network inspect microservice-network >/dev/null 2>&1 || docker network create microservice-network >/dev/null
echo "✅ 共享网络就绪"

# 停止所有现有服务
echo "🧹 停止所有现有服务..."
(cd microservices-order && docker-compose down --remove-orphans) &
(cd microservices-goods && docker-compose down --remove-orphans) &
(cd microservices-stock && docker-compose down --remove-orphans) &
wait

echo "✅ 所有服务已停止"

# 启动所有服务
echo "🚀 启动所有微服务..."

echo "  - 启动商品服务..."
(cd microservices-goods && docker-compose up -d --build > ../goods.log 2>&1) &
echo "  - 启动库存服务..."
(cd microservices-stock && docker-compose up -d --build > ../stock.log 2>&1) &
echo "  - 启动订单服务..."
(cd microservices-order && docker-compose up -d --build > ../order.log 2>&1) &
wait

echo "✅ 所有服务启动命令已执行"

# 等待服务健康检查
echo "🏥 等待服务健康检查..."
sleep 30

# 检查服务状态
echo "📊 检查服务状态..."
echo ""
echo "=== 商品服务状态 ==="
(cd microservices-goods && docker-compose ps)
echo ""
echo "=== 库存服务状态 ==="
(cd microservices-stock && docker-compose ps)
echo ""
echo "=== 订单服务状态 ==="
(cd microservices-order && docker-compose ps)

echo ""
echo "🧪 测试所有服务API..."

# 测试商品服务
echo "测试商品服务 API:"
if curl -s "http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1" > /dev/null; then
    echo "✅ 商品服务 API 测试成功"
else
    echo "❌ 商品服务 API 测试失败"
fi

sleep 5

# 通过订单网关发起下单联调
echo "测试创建订单 API:"
curl -s -X POST http://localhost:8082/api/v1/order/create \
  -H 'Content-Type: application/json' \
  -d '{"goods_id":1,"num":1,"user_id":1001,"address":"测试地址","name":"测试用户","phone":"13800138000"}' | cat

echo ""
echo "🎉 所有微服务启动完成！"
echo ""
echo "📋 服务访问地址："
echo "  🛍️  商品服务:  http://localhost:8080"
echo "  📦 库存服务:  http://localhost:8081" 
echo "  📋 订单服务:  http://localhost:8082"
echo ""
echo "💾 数据库连接："
echo "  商品 MySQL:  localhost:3306"
echo "  库存 MySQL:  localhost:3307"
echo "  订单 MySQL:  localhost:3308"
echo ""
echo "🔄 缓存连接："
echo "  商品 Redis:  localhost:6379"
echo "  库存 Redis:  localhost:6380"
echo "  订单 Redis:  localhost:6381"
echo ""
echo "📨 消息队列："
echo "  RocketMQ NameServer: localhost:9877"
echo "  RocketMQ Broker: localhost:10910, 10912"
echo ""
echo "🧪 API 测试示例："
echo "  # 商品服务"
echo "  curl \"http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1\""
echo "  curl \"http://localhost:8080/api/v1/goods/detail?goods_id=1&user_id=1\""
echo ""
echo "  # 库存服务"
echo "  curl \"http://localhost:8081/api/v1/stock/get?goods_id=1\""
echo ""
echo "  # 订单服务"
echo "  curl \"http://localhost:8082/api/v1/order/list?user_id=1&page_num=1&page_size=10\""
echo ""
echo "🔧 管理命令："
echo "  - 查看所有服务状态: bash check-all-services.sh"
echo "  - 停止所有服务: bash stop-all-services.sh"
echo "  - 查看服务日志: tail -f goods.log stock.log order.log"
echo ""
echo "📝 注意事项："
echo "  ✅ 所有服务使用不同端口，避免冲突"
echo "  ✅ 每个服务有独立的数据库和缓存"
echo "  ✅ 支持并行启动，提高启动速度"
echo "  ✅ 包含健康检查和API测试"
