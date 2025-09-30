package biz

import (
	"context"
	"time"

	orderpb "microservices-order-proto/order"
	"order-rpc-service/internal/svc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// SimpleOrderBiz 简化的订单业务逻辑（不使用 RocketMQ）
type SimpleOrderBiz struct {
	svcCtx *svc.ServiceContext
}

// NewSimpleOrderBiz 创建简化的订单业务实例
func NewSimpleOrderBiz(svcCtx *svc.ServiceContext) *SimpleOrderBiz {
	return &SimpleOrderBiz{
		svcCtx: svcCtx,
	}
}

// CreateOrder 创建订单 - 简化版本，不使用 RocketMQ
func (b *SimpleOrderBiz) CreateOrder(ctx context.Context, param *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 1. 生成订单号
	orderId := b.svcCtx.SnowflakeNode.Generate().Int64()

	// 2. 查询商品详情
	goodsDetail, err := b.svcCtx.GoodsRpc.GetGoodsDetail(ctx, &svc.GetGoodsDetailReq{
		GoodsId: param.GoodsId,
		UserId:  param.UserId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "获取商品详情失败")
	}

	// 3. 扣减库存
	_, err = b.svcCtx.StockRpc.ReduceStock(ctx, &svc.ReduceStockInfo{
		GoodsId: param.GoodsId,
		Num:     int64(param.Num),
		OrderId: orderId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "库存扣减失败")
	}

	// 4. 创建订单记录
	payAmount := int64(8999) // 模拟价格：89.99 元 = 8999 分

	orderData := Order{
		OrderId:        orderId,
		UserId:         param.UserId,
		PayAmount:      payAmount,
		ReceiveAddress: param.Address,
		ReceiveName:    param.Name,
		ReceivePhone:   param.Phone,
		Status:         "pending", // 待支付
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	orderDetail := OrderDetail{
		OrderId:   orderId,
		UserId:    param.UserId,
		GoodsId:   param.GoodsId,
		Num:       int64(param.Num),
		Title:     goodsDetail.Title,
		Status:    "pending",
		Price:     payAmount,
		Brief:     goodsDetail.Brief,
		PayAmount: payAmount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 使用事务创建订单和订单详情记录
	err = b.svcCtx.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&orderData).Error; err != nil {
			return err
		}
		if err := tx.Create(&orderDetail).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "创建订单失败")
	}

	return &orderpb.CreateOrderResponse{
		Success: true,
		Message: "订单创建成功",
		OrderId: orderId,
		Price:   "89.99",
	}, nil
}
