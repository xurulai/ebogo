package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	orderpb "microservices-order-proto/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListRequest) (resp *types.GetOrderListResponse, err error) {
	// 调用订单 RPC 服务
	orderResp, err := l.svcCtx.OrderRpc.GetOrderList(l.ctx, &orderpb.GetOrderListRequest{
		UserId:   req.UserId,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("GetOrderList RPC call failed: %v", err)
		return nil, err
	}

	// 转换响应格式
	var orders []types.OrderInfo
	for _, order := range orderResp.Orders {
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

		orders = append(orders, types.OrderInfo{
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
		})
	}

	return &types.GetOrderListResponse{
		Total:  orderResp.Total,
		Orders: orders,
	}, nil
}




