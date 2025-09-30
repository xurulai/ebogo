#!/bin/bash

# 测试修复后的微服务功能
# 使用方法：bash test_fixed_services.sh

set -e

echo "🧪 开始测试修复后的微服务..."
echo ""

# 测试商品服务
echo "===== 测试商品服务 ====="
GOODS_BASE_URL="http://localhost:8080"

echo "1. 测试获取直播间商品列表..."
curl -X GET "${GOODS_BASE_URL}/api/v1/goods/room?room_id=1&user_id=1" \
  -w "\n状态码: %{http_code}\n\n"

echo "2. 测试获取商品详情..."
curl -X GET "${GOODS_BASE_URL}/api/v1/goods/detail?goods_id=1&user_id=1" \
  -w "\n状态码: %{http_code}\n\n"

echo "3. 测试更新商品价格..."
curl -X POST "${GOODS_BASE_URL}/api/v1/goods/update" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "price": 899900}' \
  -w "\n状态码: %{http_code}\n\n"

echo "4. 再次获取商品详情验证更新..."
curl -X GET "${GOODS_BASE_URL}/api/v1/goods/detail?goods_id=1&user_id=1" \
  -w "\n状态码: %{http_code}\n\n"

echo ""
echo "===== 测试库存服务 ====="
STOCK_BASE_URL="http://localhost:8080"

echo "1. 测试设置库存..."
curl -X POST "${STOCK_BASE_URL}/api/v1/stock/set" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "stock": 100}' \
  -w "\n状态码: %{http_code}\n\n"

echo "2. 测试获取库存..."
curl -X GET "${STOCK_BASE_URL}/api/v1/stock/get?goods_id=1" \
  -w "\n状态码: %{http_code}\n\n"

echo "3. 测试扣减库存..."
curl -X POST "${STOCK_BASE_URL}/api/v1/stock/reduce" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "num": 10, "order_id": 12345}' \
  -w "\n状态码: %{http_code}\n\n"

echo "4. 再次获取库存查看变化..."
curl -X GET "${STOCK_BASE_URL}/api/v1/stock/get?goods_id=1" \
  -w "\n状态码: %{http_code}\n\n"

echo "5. 测试回滚库存..."
curl -X POST "${STOCK_BASE_URL}/api/v1/stock/rollback" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "rollback_num": 10, "order_id": 12345}' \
  -w "\n状态码: %{http_code}\n\n"

echo "✅ 测试完成！"
echo ""
echo "📊 测试说明："
echo "  - 商品服务现在使用真实的MySQL数据库和Redis缓存"
echo "  - 库存服务保持原有的分布式锁和事务逻辑"
echo "  - 所有服务都支持本地缓存和Redis缓存"
echo "  - 数据持久化到MySQL数据库"




