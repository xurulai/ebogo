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

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrderList 获取订单列表
func (l *GetOrderListLogic) GetOrderList(in *orderpb.GetOrderListRequest) (*orderpb.GetOrderListResponse, error) {
	// 参数校验
	if in.GetUserId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID无效")
	}
	if in.GetPageNum() <= 0 {
		in.PageNum = 1
	}
	if in.GetPageSize() <= 0 || in.GetPageSize() > 100 {
		in.PageSize = 10
	}

	// 创建业务逻辑实例
	orderBiz := biz.NewOrderBiz(l.svcCtx)

	// 调用业务逻辑获取订单列表
	resp, err := orderBiz.GetOrderList(l.ctx, in.UserId, in.PageNum, in.PageSize)
	if err != nil {
		l.Errorf("GetOrderList failed: %v", err)
		return nil, status.Error(codes.Internal, "获取订单列表失败")
	}

	return resp, nil
}




