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

type GetOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrderDetail 获取订单详情
func (l *GetOrderDetailLogic) GetOrderDetail(in *orderpb.GetOrderDetailRequest) (*orderpb.GetOrderDetailResponse, error) {
	// 参数校验
	if in.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "订单ID无效")
	}
	if in.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID无效")
	}

	// 创建业务逻辑实例
	orderBiz := biz.NewOrderBiz(l.svcCtx)

	// 调用业务逻辑获取订单详情
	resp, err := orderBiz.GetOrderDetail(l.ctx, in.OrderId, in.UserId)
	if err != nil {
		l.Errorf("GetOrderDetail failed: %v", err)
		return nil, status.Error(codes.Internal, "获取订单详情失败")
	}

	return resp, nil
}




