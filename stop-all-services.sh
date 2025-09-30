#!/bin/bash

# 停止所有微服务的脚本

echo "🛑 停止所有微服务..."
echo "时间: $(date)"
echo ""

# 停止所有服务
echo "停止商品服务..."
(cd microservices-goods && docker-compose down --remove-orphans) &
GOODS_PID=$!

echo "停止库存服务..."
(cd microservices-stock && docker-compose down --remove-orphans) &
STOCK_PID=$!

echo "停止订单服务..."
(cd microservices-order && docker-compose down --remove-orphans) &
ORDER_PID=$!

# 等待所有停止完成
echo "⏳ 等待所有服务停止..."
wait $GOODS_PID
wait $STOCK_PID
wait $ORDER_PID

echo "✅ 所有服务已停止"

# 恢复原始端口配置
echo "🔧 恢复原始端口配置..."

# 恢复库存服务端口配置
if [ -f "microservices-stock/docker-compose.yml.bak" ]; then
    mv microservices-stock/docker-compose.yml.bak microservices-stock/docker-compose.yml
    echo "✅ 库存服务端口配置已恢复"
fi

# 恢复订单服务端口配置  
if [ -f "microservices-order/docker-compose.yml.bak" ]; then
    mv microservices-order/docker-compose.yml.bak microservices-order/docker-compose.yml
    echo "✅ 订单服务端口配置已恢复"
fi

# 清理日志文件
echo "🧹 清理日志文件..."
rm -f goods.log stock.log order.log

echo ""
echo "🎉 所有微服务已完全停止并清理！"
echo ""
echo "💡 提示："
echo "  - 重新启动: bash start-all-microservices.sh"
echo "  - 检查状态: bash check-all-services.sh"
