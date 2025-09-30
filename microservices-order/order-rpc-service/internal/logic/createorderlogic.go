package logic

import (
	"context"

	orderpb "microservices-order-proto/order"
	"order-rpc-service/internal/biz"
	"order-rpc-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateOrder 创建订单
func (l *CreateOrderLogic) CreateOrder(in *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 参数校验
	if in.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID无效")
	}
	if in.GetGoodsId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "商品ID无效")
	}
	if in.GetNum() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "商品数量无效")
	}

	// 创建简化的业务逻辑实例
	orderBiz := biz.NewSimpleOrderBiz(l.svcCtx)

	// 调用业务逻辑创建订单
	resp, err := orderBiz.CreateOrder(l.ctx, in)
	if err != nil {
		l.Errorf("CreateOrder failed: %v", err)
		return nil, err
	}

	return resp, nil
}
