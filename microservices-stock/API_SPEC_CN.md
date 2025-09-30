# 库存服务 API 规格说明

## 概述

库存服务提供完整的商品库存管理功能，包括库存设置、查询、扣减、回滚等操作。所有接口都通过 HTTP REST API 对外提供服务，内部通过 gRPC 与库存 RPC 服务通信。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **编码**: UTF-8

## 接口列表

### 1. 设置库存

设置指定商品的库存数量。

**接口地址**: `POST /api/v1/stock/set`

**请求参数**:
```json
{
  "goods_id": 1,     // 商品ID，必填，大于0
  "stock": 100       // 库存数量，必填，大于等于0
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "设置库存成功"
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "设置库存失败"
}
```

### 2. 获取库存

查询指定商品的当前库存数量。

**接口地址**: `GET /api/v1/stock/get`

**请求参数**:
- `goods_id`: 商品ID，必填，大于0

**请求示例**:
```
GET /api/v1/stock/get?goods_id=1
```

**响应示例**:
```json
{
  "success": true,
  "message": "获取库存成功",
  "data": {
    "goods_id": 1,
    "stock": 100
  }
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "获取库存失败"
}
```

### 3. 扣减库存

扣减指定商品的库存数量，用于订单下单时的库存预扣。

**接口地址**: `POST /api/v1/stock/reduce`

**请求参数**:
```json
{
  "goods_id": 1,     // 商品ID，必填，大于0
  "num": 10,         // 扣减数量，必填，大于0
  "order_id": 12345  // 订单ID，必填，大于0
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "扣减库存成功"
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "扣减库存失败"
}
```

**业务说明**:
- 扣减库存时会检查可用库存（总库存 - 锁定库存）
- 扣减成功后会增加锁定库存，并创建库存扣减记录
- 使用分布式锁确保并发安全

### 4. 回滚库存

回滚之前扣减的库存，用于订单取消时的库存回滚。

**接口地址**: `POST /api/v1/stock/rollback`

**请求参数**:
```json
{
  "goods_id": 1,       // 商品ID，必填，大于0
  "rollback_num": 10,  // 回滚数量，必填，大于0
  "order_id": 12345    // 订单ID，必填，大于0
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "库存回滚成功"
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "库存回滚失败"
}
```

**业务说明**:
- 回滚时会查找对应的库存扣减记录
- 只有状态为"已扣减"的记录才能回滚
- 回滚成功后会增加总库存，减少锁定库存，并更新记录状态为"已回滚"
- 支持幂等操作，重复回滚不会重复执行

### 5. 批量获取库存

批量查询多个商品的库存信息。

**接口地址**: `POST /api/v1/stock/batch/get`

**请求参数**:
```json
{
  "items": [
    {
      "goods_id": 1,
      "stock": 0    // 查询时此字段可以为0
    },
    {
      "goods_id": 2,
      "stock": 0
    }
  ]
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "批量获取库存成功",
  "data": {
    "items": [
      {
        "goods_id": 1,
        "stock": 100
      },
      {
        "goods_id": 2,
        "stock": 200
      }
    ]
  }
}
```

### 6. 批量扣减库存

批量扣减多个商品的库存。

**接口地址**: `POST /api/v1/stock/batch/reduce`

**请求参数**:
```json
{
  "items": [
    {
      "goods_id": 1,
      "stock": 5    // 要扣减的数量
    },
    {
      "goods_id": 2,
      "stock": 3
    }
  ]
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "批量扣减库存成功"
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "部分商品库存扣减失败"
}
```

**业务说明**:
- 批量操作中，如果部分商品扣减失败，会返回失败状态
- 实际生产环境中建议使用分布式事务确保一致性

## 错误码说明

| HTTP状态码 | 说明 |
|-----------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 500 | 服务器内部错误 |

## 数据类型说明

| 字段类型 | 说明 | 范围 |
|---------|------|------|
| goods_id | 商品ID | int64，大于0 |
| stock | 库存数量 | int64，大于等于0 |
| num | 操作数量 | int64，大于0 |
| order_id | 订单ID | int64，大于0 |
| rollback_num | 回滚数量 | int64，大于0 |

## 业务流程

### 典型的订单流程

1. **下单时**:
   ```
   POST /api/v1/stock/reduce
   {
     "goods_id": 1,
     "num": 2,
     "order_id": 12345
   }
   ```

2. **订单支付成功** - 库存扣减保持不变

3. **订单取消时**:
   ```
   POST /api/v1/stock/rollback
   {
     "goods_id": 1,
     "rollback_num": 2,
     "order_id": 12345
   }
   ```

### 库存管理流程

1. **初始化库存**:
   ```
   POST /api/v1/stock/set
   {
     "goods_id": 1,
     "stock": 1000
   }
   ```

2. **查询库存**:
   ```
   GET /api/v1/stock/get?goods_id=1
   ```

3. **批量查询**:
   ```
   POST /api/v1/stock/batch/get
   {
     "items": [{"goods_id": 1, "stock": 0}, {"goods_id": 2, "stock": 0}]
   }
   ```

## 注意事项

1. **并发安全**: 所有库存变更操作都使用分布式锁确保并发安全
2. **事务一致性**: 库存扣减和记录创建在同一事务中完成
3. **幂等性**: 库存回滚支持幂等操作，重复调用不会产生副作用
4. **性能优化**: 建议在高并发场景下使用批量接口减少网络开销
5. **监控告警**: 建议对库存不足、回滚失败等异常情况设置监控告警




