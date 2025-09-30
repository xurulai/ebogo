# 商品服务（microservices-goods）

基于 go-zero 的商品微服务示例，包含 API 网关与商品 RPC 服务，并已统一使用共享 Proto 模块。

> 接口字段详解请见：`API_SPEC_CN.md`

## 架构概览

```
客户端 → API Gateway(HTTP) ── gRPC ─→ Goods RPC Service → 业务/存储层
                         │
                         └─ 统一鉴权、限流、日志、监控
```

- **API Gateway**: 对外统一 HTTP 入口，负责参数校验、协议转换（HTTP→gRPC）、结果映射。
- **Goods RPC Service**: 实现商品领域能力（直播间商品列表、商品详情、更新价格等）。
- **Shared Proto Module**: `microservices-goods/proto`，两端共享的 gRPC/Protobuf 定义与生成代码。

## 目录结构

```
microservices-goods/
├── api-gateway/                  # API 网关
│   ├── etc/                      # 网关配置（yaml）
│   ├── gateway.go                # 网关入口
│   ├── internal/
│   │   ├── config/               # 配置结构体（已中文注释）
│   │   ├── handler/              # 路由与 handler
│   │   ├── logic/                # 网关业务逻辑（参数映射、RPC 调用、出参映射）
│   │   └── svc/                  # ServiceContext: 注入 RPC 客户端
│   └── go.mod                    # 引用 shared proto（通过 replace）
├── goods-rpc-service/            # 商品 RPC 服务
│   ├── etc/                      # 服务端配置（yaml）
│   ├── goods-rpc.go              # gRPC 服务入口
│   ├── internal/
│   │   ├── biz/                  # 领域业务（当前为演示用静态数据）
│   │   ├── config/               # 配置结构体（已中文注释）
│   │   ├── logic/                # 参数校验、调用 biz、错误码
│   │   ├── server/               # gRPC Server 实现（委派到 logic）
│   │   └── svc/                  # ServiceContext: 注入配置/依赖
│   └── go.mod                    # 引用 shared proto（通过 replace）
├── proto/                        # 共享 Proto 模块（单一事实源）
│   ├── go.mod                    # module microservices-goods-proto
│   └── goods/
│       ├── goods.proto           # gRPC/Protobuf 定义
│       ├── goods.pb.go           # 生成的消息类型
│       └── goods_grpc.pb.go      # 生成的 gRPC 接口
└── scripts/
    └── gen-proto.sh             # 仅对共享模块生成代码
```

## 启动与构建

- 安装依赖：需要 `protoc`、`protoc-gen-go`、`protoc-gen-go-grpc`
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5`
  - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1`
- 生成 proto（如有改动）：
  - `bash scripts/gen-proto.sh`
- 构建：
  - API 网关：`cd api-gateway && go build ./...`
  - RPC 服务：`cd goods-rpc-service && go build ./...`

配置文件样例：
- 网关：`api-gateway/etc/gateway-api.yaml`
- 服务：`goods-rpc-service/etc/goods-rpc.yaml`

## 配置说明（节选）

- 网关 `internal/config/config.go`
  - `RestConf`: go-zero HTTP 服务基础配置（端口、超时、日志等）
  - `GoodsRpc`: gRPC 客户端配置（直连地址或注册中心）
- 服务 `internal/config/config.go`
  - `RpcServerConf`: gRPC 服务端配置（监听地址、模式、限流、链路追踪等）

## Proto 与接口

- 单一事实源：`proto/goods/goods.proto`
- `go_package`: `microservices-goods-proto/goods`
- 接口（节选）：
  - `GetGoodsByRoom(GetGoodsByRoomRequest) returns (GetGoodsByRoomResponse)`
  - `GetGoodsDetail(GetGoodsDetailRequest) returns (GetGoodsDetailResponse)`
  - `UpdateGoodsDetail(UpdateGoodsDetailRequest) returns (UpdateGoodsDetailResponse)`

## 业务流程（核心用例）

- 获取直播间商品列表
  1. 客户端调用网关 GET `/api/v1/goods/room?room_id=&user_id=`
  2. 网关 `logic/gateway` 组装 `GetGoodsByRoomRequest` 并调用 RPC `GoodsService.GetGoodsByRoom`
  3. RPC `logic` 校验参数 → 调用 `biz.GoodsBiz.GetGoodsByRoom`（当前静态数据）
  4. RPC 返回 `GetGoodsByRoomResponse`，网关转换为对外响应结构并返回

- 获取商品详情
  1. 客户端调用网关 GET `/api/v1/goods/detail?goods_id=&user_id=`
  2. 网关调用 RPC `GoodsService.GetGoodsDetail`
  3. RPC `logic` 校验参数 → `biz.GetGoodsDetail`
  4. 返回 `GoodsDetailResponse`，网关响应

- 更新商品详情（价格）
  1. 客户端调用网关 POST `/api/v1/goods/update`，JSON: `{goods_id, price}`（分）
  2. 网关调用 RPC `GoodsService.UpdateGoodsDetail`
  3. RPC `logic` 校验参数 → `biz.UpdateGoodsDetail`
  4. 返回成功/失败结果

## 设计说明

- 分层职责
  - 网关：HTTP 参数/结果与 RPC 的映射，聚合与编排，横切能力（鉴权、限流、日志）
  - RPC：参数校验、错误码统一、领域业务调度
  - Biz：领域逻辑（当前示例使用静态数据，可替换为DB/Redis/MQ等）
- 共享 Proto 模块
  - 杜绝“双份 proto 定义”带来的包路径与版本漂移
  - 两端通过 `replace microservices-goods-proto => ../proto` 引入

## API 示例

```bash
# 获取直播间商品列表
curl "http://localhost:8080/api/v1/goods/room?room_id=1&user_id=1"

# 获取商品详情
curl "http://localhost:8080/api/v1/goods/detail?goods_id=1&user_id=1"

# 更新商品价格（单位：分）
curl -X POST "http://localhost:8080/api/v1/goods/update" \
  -H "Content-Type: application/json" \
  -d '{"goods_id": 1, "price": 799900}'
```

## 后续可扩展

- 将 `biz` 接入真实 MySQL/Redis，并在 `svc` 中注入连接
- 引入注册中心（如 Etcd/Consul），替代直连配置
- 加入鉴权中间件与统一错误码映射
- 增加熔断/限流与链路追踪配置


# 获取商品详情
grpcurl -plaintext -d '{"goods_id":1,"user_id":123}' 127.0.0.1:9001 goods.GoodsService/GetGoodsDetail

# 获取直播间商品
grpcurl -plaintext -d '{"room_id":1,"user_id":123}' 127.0.0.1:9001 goods.GoodsService/GetGoodsByRoom

# 更新商品价格（单位：分）
grpcurl -plaintext -d '{"goods_id":1,"price":8999}' 127.0.0.1:9001 goods.GoodsService/UpdateGoodsDetail
