package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	orderpb "microservices-order-proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (resp *types.CreateOrderResponse, err error) {
	// 调用订单 RPC 服务
	orderResp, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &orderpb.CreateOrderRequest{
		GoodsId: req.GoodsId,
		Num:     req.Num,
		UserId:  req.UserId,
		Address: req.Address,
		Name:    req.Name,
		Phone:   req.Phone,
	})
	if err != nil {
		l.Errorf("CreateOrder RPC call failed: %v", err)
		return nil, err
	}

	return &types.CreateOrderResponse{
		Success: orderResp.Success,
		Message: orderResp.Message,
		OrderId: orderResp.OrderId,
		Price:   orderResp.Price,
	}, nil
}
