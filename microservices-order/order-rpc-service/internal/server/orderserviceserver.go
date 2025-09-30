package server

import (
	"context"

	orderpb "microservices-order-proto/order"
	"order-rpc-service/internal/logic"
	"order-rpc-service/internal/svc"
)

type OrderServiceServer struct {
	svcCtx *svc.ServiceContext
	orderpb.UnimplementedOrderServiceServer
}

func NewOrderServiceServer(svcCtx *svc.ServiceContext) *OrderServiceServer {
	return &OrderServiceServer{
		svcCtx: svcCtx,
	}
}

// CreateOrder 创建订单
func (s *OrderServiceServer) CreateOrder(ctx context.Context, in *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	l := logic.NewCreateOrderLogic(ctx, s.svcCtx)
	return l.CreateOrder(in)
}

// GetOrderList 获取订单列表
func (s *OrderServiceServer) GetOrderList(ctx context.Context, in *orderpb.GetOrderListRequest) (*orderpb.GetOrderListResponse, error) {
	l := logic.NewGetOrderListLogic(ctx, s.svcCtx)
	return l.GetOrderList(in)
}

// GetOrderDetail 获取订单详情
func (s *OrderServiceServer) GetOrderDetail(ctx context.Context, in *orderpb.GetOrderDetailRequest) (*orderpb.GetOrderDetailResponse, error) {
	l := logic.NewGetOrderDetailLogic(ctx, s.svcCtx)
	return l.GetOrderDetail(in)
}

// UpdateOrderStatus 更新订单状态
func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, in *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	l := logic.NewUpdateOrderStatusLogic(ctx, s.svcCtx)
	return l.UpdateOrderStatus(in)
}




