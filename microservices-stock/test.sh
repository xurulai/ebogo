#!/bin/bash

# 库存微服务测试脚本
# 使用方法：bash test.sh

set -e

BASE_URL="http://localhost:8080"

echo "🧪 开始测试库存微服务 API..."
echo ""

# 测试设置库存
echo "1. 测试设置库存..."
curl -X POST "${BASE_URL}/api/v1/stock/set" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "stock": 100}' \
  -w "\n状态码: %{http_code}\n\n"

# 测试获取库存
echo "2. 测试获取库存..."
curl -X GET "${BASE_URL}/api/v1/stock/get?goods_id=1" \
  -w "\n状态码: %{http_code}\n\n"

# 测试扣减库存
echo "3. 测试扣减库存..."
curl -X POST "${BASE_URL}/api/v1/stock/reduce" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "num": 10, "order_id": 12345}' \
  -w "\n状态码: %{http_code}\n\n"

# 再次获取库存查看变化
echo "4. 再次获取库存查看变化..."
curl -X GET "${BASE_URL}/api/v1/stock/get?goods_id=1" \
  -w "\n状态码: %{http_code}\n\n"

# 测试回滚库存
echo "5. 测试回滚库存..."
curl -X POST "${BASE_URL}/api/v1/stock/rollback" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "rollback_num": 10, "order_id": 12345}' \
  -w "\n状态码: %{http_code}\n\n"

# 测试批量获取库存
echo "6. 测试批量获取库存..."
curl -X POST "${BASE_URL}/api/v1/stock/batch/get" \
  -H "Content-Type: application/json" \
  -d '{"items": [{"goods_id": 1, "stock": 0}, {"goods_id": 2, "stock": 0}]}' \
  -w "\n状态码: %{http_code}\n\n"

echo "✅ 测试完成！"




