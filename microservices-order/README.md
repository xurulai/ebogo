# 订单微服务 (microservices-order)

基于 go-zero 框架的订单微服务，采用 API Gateway + RPC Service 架构。

## 🚀 快速启动

### 方法一：Docker 环境（推荐）
```bash
# 启动所有服务（MySQL, Redis, RocketMQ, 订单服务）
bash start.sh

# 测试订单创建
bash docker-test.sh

# 停止所有服务
bash stop.sh
```

### 方法二：本地开发
```bash
# 启动基础服务
docker-compose up -d mysql redis rocketmq-nameserver rocketmq-broker

# 启动订单 RPC 服务
cd order-rpc-service
./order-rpc-simple -f etc/order-rpc-simple.yaml

# 启动 API Gateway
cd api-gateway  
./gateway -f etc/gateway-api.yaml
```

## 🧪 测试

### gRPC 直接测试（推荐）
```bash
# 运行完整测试套件
bash test-order-grpc.sh

# 单次快速测试
grpcurl -plaintext -d '{
  "goods_id": 1,
  "num": 1,
  "user_id": 1001,
  "address": "测试地址",
  "name": "测试用户",
  "phone": "13800138000"
}' localhost:9002 order.OrderService/CreateOrder
```

### HTTP API 测试
```bash
# 通过 API Gateway
curl -X POST http://localhost:8080/api/v1/order/create \
  -H "Content-Type: application/json" \
  -d '{
    "goods_id": 1,
    "num": 1,
    "user_id": 1001,
    "address": "测试地址",
    "name": "测试用户",
    "phone": "13800138000"
  }'
```

## 📊 服务端口

- **API Gateway**: 8080
- **订单 RPC**: 9002
- **MySQL**: 3306
- **Redis**: 6379
- **RocketMQ NameServer**: 9876
- **RocketMQ Broker**: 10911

## 🏗️ 架构

```
HTTP Request → API Gateway (8080) → 订单 RPC (9002) → MySQL/Redis
                                  ↓
                              商品服务 RPC (模拟)
                              库存服务 RPC (模拟)
```

## 📁 核心文件

- `start.sh` - 启动脚本
- `stop.sh` - 停止脚本  
- `docker-test.sh` - Docker 环境测试
- `test-order-grpc.sh` - gRPC 测试套件
- `docker-compose.yml` - Docker 编排配置


