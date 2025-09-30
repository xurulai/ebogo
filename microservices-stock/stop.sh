#!/bin/bash

# 库存微服务停止脚本
# 使用方法：bash stop.sh

set -e

echo "正在停止库存微服务..."

# 停止并删除容器
docker-compose down

echo "✅ 库存微服务已停止！"




