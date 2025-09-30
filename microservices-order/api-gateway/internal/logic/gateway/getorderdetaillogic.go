package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	orderpb "microservices-order-proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(req *types.GetOrderDetailRequest) (resp *types.GetOrderDetailResponse, err error) {
	// 调用订单 RPC 服务
	orderResp, err := l.svcCtx.OrderRpc.GetOrderDetail(l.ctx, &orderpb.GetOrderDetailRequest{
		OrderId: req.OrderId,
		UserId:  req.UserId,
	})
	if err != nil {
		l.Errorf("GetOrderDetail RPC call failed: %v", err)
		return nil, err
	}

	// 转换响应格式
	order := orderResp.Order
	var orderDetails []types.OrderDetailInfo
	for _, detail := range order.OrderDetails {
		orderDetails = append(orderDetails, types.OrderDetailInfo{
			OrderId:   detail.OrderId,
			GoodsId:   detail.GoodsId,
			Num:       detail.Num,
			Title:     detail.Title,
			Status:    detail.Status,
			Price:     detail.Price,
			Brief:     detail.Brief,
			PayAmount: detail.PayAmount,
		})
	}

	return &types.GetOrderDetailResponse{
		Order: types.OrderInfo{
			OrderId:        order.OrderId,
			UserId:         order.UserId,
			Status:         order.Status,
			PayChannel:     order.PayChannel,
			PayAmount:      order.PayAmount,
			ReceiveAddress: order.ReceiveAddress,
			ReceiveName:    order.ReceiveName,
			ReceivePhone:   order.ReceivePhone,
			CreatedAt:      order.CreatedAt,
			OrderDetails:   orderDetails,
		},
	}, nil
}




