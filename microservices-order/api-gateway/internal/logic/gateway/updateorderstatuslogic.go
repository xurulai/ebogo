package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	orderpb "microservices-order-proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateOrderStatusLogic) UpdateOrderStatus(req *types.UpdateOrderStatusRequest) (resp *types.UpdateOrderStatusResponse, err error) {
	// 调用订单 RPC 服务
	orderResp, err := l.svcCtx.OrderRpc.UpdateOrderStatus(l.ctx, &orderpb.UpdateOrderStatusRequest{
		OrderId: req.OrderId,
		Status:  req.Status,
	})
	if err != nil {
		l.Errorf("UpdateOrderStatus RPC call failed: %v", err)
		return nil, err
	}

	return &types.UpdateOrderStatusResponse{
		Success: orderResp.Success,
		Message: orderResp.Message,
	}, nil
}




