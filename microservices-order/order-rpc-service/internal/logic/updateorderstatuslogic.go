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

type UpdateOrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderStatusLogic {
	return &UpdateOrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrderStatus 更新订单状态
func (l *UpdateOrderStatusLogic) UpdateOrderStatus(in *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	// 参数校验
	if in.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "订单ID无效")
	}
	if len(in.GetStatus()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "订单状态不能为空")
	}

	// 创建业务逻辑实例
	orderBiz := biz.NewOrderBiz(l.svcCtx)

	// 调用业务逻辑更新订单状态
	err := orderBiz.UpdateOrderStatus(l.ctx, in.OrderId, in.Status)
	if err != nil {
		l.Errorf("UpdateOrderStatus failed: %v", err)
		return nil, status.Error(codes.Internal, "更新订单状态失败")
	}

	return &orderpb.UpdateOrderStatusResponse{
		Success: true,
		Message: "订单状态更新成功",
	}, nil
}




