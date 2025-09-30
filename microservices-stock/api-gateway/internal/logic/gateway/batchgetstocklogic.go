package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	stockpb "api-gateway/proto/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetStockLogic {
	return &BatchGetStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchGetStockLogic) BatchGetStock(req *types.BatchStockRequest) (resp *types.CommonResponse, err error) {
	// 组装 gRPC 请求
	var stockItems []*stockpb.GoodsStockInfo
	for _, item := range req.Items {
		stockItems = append(stockItems, &stockpb.GoodsStockInfo{
			GoodsId: item.GoodsId,
			Stock:   item.Stock,
		})
	}

	stockReq := &stockpb.StockInfoList{
		Data: stockItems,
	}

	// 调用 RPC 服务
	stockResp, err := l.svcCtx.StockRpc.BatchGetStock(l.ctx, stockReq)
	if err != nil {
		l.Errorf("BatchGetStock RPC call failed: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "批量获取库存失败",
		}, nil
	}

	// 转换响应数据
	var items []types.StockResponse
	for _, item := range stockResp.Data {
		items = append(items, types.StockResponse{
			GoodsId: item.GoodsId,
			Stock:   item.Stock,
		})
	}

	data := types.BatchStockResponse{
		Items: items,
	}

	// 返回响应
	return &types.CommonResponse{
		Success: true,
		Message: "批量获取库存成功",
		Data:    data,
	}, nil
}
