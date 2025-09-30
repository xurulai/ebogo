# 订单服务调试指南

## 🐛 VS Code 调试配置

已在 `.vscode/launch.json` 中添加了以下调试配置：

### 单独调试配置

1. **Launch Order RPC Service (9002)** - 调试订单 RPC 服务
   - 使用 `order-rpc-simple.yaml` 配置
   - 端口：9002
   - 推荐用于日常开发调试

2. **Launch Order RPC Service (Docker Config)** - 使用 Docker 配置调试
   - 使用 `order-rpc-docker.yaml` 配置
   - 需要先启动 Docker 基础服务

3. **Launch Order API Gateway (8080)** - 调试 API Gateway
   - 端口：8080
   - 需要先启动 RPC 服务

### 组合调试配置

4. **Launch Order Gateway + RPC** - 同时调试网关和 RPC 服务
5. **Launch Order RPC Only** - 仅调试 RPC 服务

## 🚀 调试步骤

### 方法一：调试 RPC 服务（推荐）

1. **启动基础服务**：
   ```bash
   cd /Users/xurulai/ebogo/microservices-order
   docker-compose up -d mysql redis rocketmq-nameserver rocketmq-broker
   ```

2. **在 VS Code 中**：
   - 按 `F5` 或 `Ctrl+Shift+D` 打开调试面板
   - 选择 **"Launch Order RPC Service (9002)"**
   - 点击绿色播放按钮开始调试

3. **设置断点**：
   - 在 `order-rpc-service/internal/logic/createorderlogic.go` 设置断点
   - 在 `order-rpc-service/internal/biz/order_simple.go` 设置断点

4. **测试调试**：
   ```bash
   grpcurl -plaintext -d '{
     "goods_id": 1,
     "num": 1,
     "user_id": 1001,
     "address": "调试测试地址",
     "name": "调试用户",
     "phone": "13800138000"
   }' localhost:9002 order.OrderService/CreateOrder
   ```

### 方法二：调试完整流程

1. **启动基础服务**（同上）

2. **在 VS Code 中**：
   - 选择 **"Launch Order Gateway + RPC"**
   - 同时调试网关和 RPC 服务

3. **设置断点**：
   - API Gateway: `api-gateway/internal/logic/gateway/createorderlogic.go`
   - RPC Service: `order-rpc-service/internal/logic/createorderlogic.go`

4. **测试调试**：
   ```bash
   curl -X POST http://localhost:8080/api/v1/order/create \
     -H "Content-Type: application/json" \
     -d '{
       "goods_id": 1,
       "num": 1,
       "user_id": 1001,
       "address": "调试测试地址",
       "name": "调试用户",
       "phone": "13800138000"
     }'
   ```

## 🎯 关键调试点

### 1. 订单创建流程
- **入口**: `createorderlogic.go:30` - CreateOrder 方法
- **参数验证**: `createorderlogic.go:32-40`
- **业务逻辑**: `order_simple.go:27` - CreateOrder 方法

### 2. 数据库操作
- **商品查询**: `order_simple.go:32` - GetGoodsDetail
- **库存扣减**: `order_simple.go:40` - ReduceStock
- **订单创建**: `order_simple.go:75` - 数据库事务

### 3. 错误处理
- **参数错误**: `createorderlogic.go:32-40`
- **业务错误**: `order_simple.go:34, 42, 89`

## 🔧 调试技巧

### 1. 查看变量
- 在断点处查看 `req` 请求参数
- 查看 `goodsDetail` 商品信息
- 查看 `orderData` 和 `orderDetail` 数据

### 2. 调试数据库
- 在事务前后设置断点
- 查看 SQL 执行日志
- 验证数据库记录

### 3. 调试 RPC 调用
- 在 `servicecontext.go:99` - `GetGoodsDetail` RPC 调用
- 在 `servicecontext.go:150` - `ReduceStock` RPC 调用
- 查看真实 RPC 响应或降级处理
- 检查服务连接状态和错误处理

## 📊 调试环境要求

### 必需服务
- ✅ MySQL (端口 3306)
- ✅ Redis (端口 6379)
- ⚠️ RocketMQ (可选，用于完整功能)
- 🔗 商品服务 (端口 9001，可选，不可用时使用降级处理)
- 🔗 库存服务 (端口 9003，可选，不可用时使用降级处理)

### 检查服务状态
```bash
# 检查 Docker 服务
docker-compose ps

# 检查端口占用
nc -z localhost 3306 && echo "MySQL OK"
nc -z localhost 6379 && echo "Redis OK"
nc -z localhost 9876 && echo "RocketMQ OK"
nc -z localhost 9001 && echo "Goods Service OK"
nc -z localhost 9003 && echo "Stock Service OK"

# 运行 RPC 集成测试
bash test-rpc-integration.sh
```

## 🚨 常见问题

### 1. 端口冲突
如果 9002 端口被占用：
```bash
lsof -ti:9002 | xargs kill -9
```

### 2. 数据库连接失败
检查 MySQL 服务和配置：
```bash
docker-compose logs mysql
```

### 3. 调试器无法启动
确保 Go 扩展已安装并更新到最新版本。

## 💡 调试建议

1. **优先调试 RPC 服务**：更直接，问题定位更准确
2. **使用 `order-rpc-simple.yaml`**：配置简单，依赖最少
3. **先测试单个功能**：逐步调试，避免复杂问题
4. **查看日志输出**：结合控制台日志分析问题
5. **使用 grpcurl 测试**：比 HTTP 请求更直接
