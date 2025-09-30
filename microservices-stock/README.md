# 库存服务（microservices-stock）

基于 go-zero 的库存微服务示例，包含 API 网关与库存 RPC 服务，并已统一使用共享 Proto 模块。

## 架构概览

```
客户端 → API Gateway(HTTP) ── gRPC ─→ Stock RPC Service → 业务/存储层
                         │
                         └─ 统一鉴权、限流、日志、监控
```

- **API Gateway**: 对外统一 HTTP 入口，负责参数校验、协议转换（HTTP→gRPC）、结果映射。
- **Stock RPC Service**: 实现库存领域能力（设置库存、获取库存、扣减库存、回滚库存等）。
- **Shared Proto Module**: `microservices-stock/proto`，两端共享的 gRPC/Protobuf 定义与生成代码。

## 目录结构

```
microservices-stock/
├── api-gateway/                  # API 网关
│   ├── etc/                      # 网关配置（yaml）
│   ├── gateway.go                # 网关入口
│   ├── internal/
│   │   ├── config/               # 配置结构体
│   │   ├── handler/              # 路由与 handler
│   │   ├── logic/                # 网关业务逻辑（参数映射、RPC 调用、出参映射）
│   │   └── svc/                  # ServiceContext: 注入 RPC 客户端
│   └── go.mod                    # 引用 shared proto（通过 replace）
├── stock-rpc-service/            # 库存 RPC 服务
│   ├── etc/                      # 服务端配置（yaml）
│   ├── stock-rpc.go              # gRPC 服务入口
│   ├── internal/
│   │   ├── biz/                  # 领域业务逻辑
│   │   ├── config/               # 配置结构体
│   │   ├── logic/                # 参数校验、调用 biz、错误码
│   │   ├── server/               # gRPC Server 实现（委派到 logic）
│   │   └── svc/                  # ServiceContext: 注入配置/依赖
│   └── go.mod                    # 引用 shared proto（通过 replace）
├── proto/                        # 共享 Proto 模块（单一事实源）
│   ├── go.mod                    # module microservices-stock-proto
│   └── stock/
│       ├── stock.proto           # gRPC/Protobuf 定义
│       ├── stock.pb.go           # 生成的消息类型
│       └── stock_grpc.pb.go      # 生成的 gRPC 接口
└── scripts/
    └── gen-proto.sh             # 仅对共享模块生成代码
```

## 快速启动

### 方式一：使用 Docker（推荐）

```bash
# 启动所有服务（包括 MySQL、Redis、API Gateway、Stock RPC Service）
bash start.sh

# 测试 API
bash test.sh

# 停止服务
bash stop.sh
```

### 方式二：本地开发

1. **安装依赖**：
   ```bash
   # 安装 protoc
   # macOS: brew install protobuf
   # Ubuntu: sudo apt-get install protobuf-compiler
   
   # 安装 Go 插件
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
   ```

2. **生成 proto 代码**：
   ```bash
   bash scripts/gen-proto.sh
   ```

3. **启动 MySQL 和 Redis**：
   ```bash
   docker-compose up mysql redis -d
   ```

4. **启动 RPC 服务**：
   ```bash
   cd stock-rpc-service
   go run stock-rpc.go
   ```

5. **启动 API 网关**：
   ```bash
   cd api-gateway
   go run gateway.go
   ```

## API 接口

### 1. 设置库存
```bash
curl -X POST "http://localhost:8080/api/v1/stock/set" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "stock": 100}'
```

### 2. 获取库存
```bash
curl "http://localhost:8080/api/v1/stock/get?goods_id=1"
```

### 3. 扣减库存
```bash
curl -X POST "http://localhost:8080/api/v1/stock/reduce" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "num": 10, "order_id": 12345}'
```

### 4. 回滚库存
```bash
curl -X POST "http://localhost:8080/api/v1/stock/rollback" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "rollback_num": 10, "order_id": 12345}'
```

### 5. 批量获取库存
```bash
curl -X POST "http://localhost:8080/api/v1/stock/batch/get" \
  -H "Content-Type: application/json" \
  -d '{"items": [{"goods_id": 1, "stock": 0}, {"goods_id": 2, "stock": 0}]}'
```

### 6. 批量扣减库存
```bash
curl -X POST "http://localhost:8080/api/v1/stock/batch/reduce" \
  -H "Content-Type: application/json" \
  -d '{"items": [{"goods_id": 1, "stock": 5}, {"goods_id": 2, "stock": 3}]}'


  
grpcurl -plaintext -d '{"goodsId": 1}' localhost:9003 stock.StockService/GetStock
```

## 核心特性

### 1. 分布式锁
- 使用 Redis 分布式锁确保库存操作的原子性
- 防止并发扣减导致的库存超卖问题

### 2. 事务支持
- 使用 GORM 事务确保数据一致性
- 库存扣减和记录创建在同一事务中完成

### 3. 库存记录
- 记录每次库存变更的详细信息
- 支持幂等性回滚操作

### 4. 高可用设计
- 服务间通过 gRPC 通信，性能更高
- 支持水平扩展和负载均衡
- 完整的健康检查机制

## 数据库设计

### 库存表 (xx_stock)
```sql
CREATE TABLE `xx_stock` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `goods_id` bigint NOT NULL COMMENT '商品ID',
  `stocknum` bigint NOT NULL DEFAULT '0' COMMENT '库存数量',
  `lock` bigint NOT NULL DEFAULT '0' COMMENT '锁定库存',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_id` (`goods_id`)
);
```

### 库存记录表 (xx_stock_record)
```sql
CREATE TABLE `xx_stock_record` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `goods_id` bigint NOT NULL COMMENT '商品ID',
  `num` bigint NOT NULL COMMENT '数量',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态：1-已扣减，3-已回滚',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_goods` (`order_id`, `goods_id`)
);
```

## 配置说明

### API 网关配置
- `api-gateway/etc/gateway-api.yaml`: 本地开发配置
- `api-gateway/etc/gateway-api-docker.yaml`: Docker 环境配置

### RPC 服务配置
- `stock-rpc-service/etc/stock-rpc.yaml`: 本地开发配置
- `stock-rpc-service/etc/stock-rpc-docker.yaml`: Docker 环境配置

## 监控与日志

- 所有服务都集成了 go-zero 的日志和监控功能
- 支持链路追踪和性能监控
- Docker 环境下可通过 `docker-compose logs` 查看日志

## 扩展说明

1. **注册中心**: 可集成 Etcd/Consul 替代直连配置
2. **消息队列**: 可集成 RocketMQ/Kafka 处理异步库存回滚
3. **缓存优化**: 可增加 Redis 缓存提升查询性能
4. **监控告警**: 可集成 Prometheus + Grafana 监控
5. **限流熔断**: 可配置 go-zero 的限流和熔断机制

